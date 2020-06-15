package xmp

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"
)

const (
	rdfXmlNamespace = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

var (
	// ErrPropertyNotFound represents an error for a get operation that produced
	// no results.
	ErrPropertyNotFound = errors.New("property not found")
)

var (
	rdfTag = xml.Name{
		Space: rdfXmlNamespace,
		Local: "RDF",
	}

	rdfDescriptionTag = xml.Name{
		Space: rdfXmlNamespace,
		Local: "Description",
	}
)

var (
	// standardXpacketId is the expected value of the xpacket ID for XMP data.
	standardXpacketId = "W5M0MpCehiHzreSzNTczkc9d"
)

type rawAttributeAssignment string

// Parse will return a name and value if the string looks like
// 'name="value"'. Else, nil.
func (raa rawAttributeAssignment) Parse() (string, string) {
	len_ := len(raa)

	// At least one character in the name, the equals, and the quote characters
	// on the value.
	if len_ < 4 {
		return "", ""
	}

	// Split. Make sure there *was* an equals. If there was more than one, the
	// others must have been in the value.

	parts := strings.SplitN(string(raa), "=", 2)
	if len(parts) != 2 {
		return "", ""
	}

	name := parts[0]
	valueRaw := parts[1]

	// Validate the start-quote of the value.

	valueQuoteChar := valueRaw[0]

	if valueQuoteChar != '"' && valueQuoteChar != '\'' {
		return "", ""
	}

	// Make sure the end-quote of the value matches the start quote character.

	if valueRaw[len(valueRaw)-1] != valueQuoteChar {
		return "", ""
	}

	value := valueRaw[1 : len(valueRaw)-1]

	// Unescape any quotes.
	value = strings.ReplaceAll(value, "\\\"", "\"")
	value = strings.ReplaceAll(value, "\\\\", "\\")

	return name, value
}

// Parser parses an XMP document.
type Parser struct {
	xd *xml.Decoder

	// TODO(dustin): !! Add an accessor for this. Investigate whether we even have this value in our test-data.
	// encodingSequence will be nil if xpacket has empty "start" value.
	encodingSequence bom.Sequence

	packetIsOpen         bool
	rdfIsOpen            bool
	rdfDescriptionIsOpen bool

	stack []XmlName

	lastCharData *string

	lastToken xml.Token
}

// NewParser returns a new Parser struct.
func NewParser(r io.Reader) *Parser {
	xd := xml.NewDecoder(r)

	stack := make([]XmlName, 0)

	return &Parser{
		xd:    xd,
		stack: stack,
	}
}

type XmlName xml.Name

func (xn XmlName) String() string {
	prefix := LookupPreferredNamespacePrefix(xn.Space)
	if prefix == "" {
		// They should notify us of the unknown namespace so that we
		// can register it and they can handle it properly.
		prefix = "?"
	}

	return fmt.Sprintf("[%s]%s", prefix, xn.Local)
}

type XmpPropertyName []XmlName

func (xpn XmpPropertyName) Parts() (parts []string) {
	parts = make([]string, len(xpn))
	for i, tag := range xpn {
		prefix := LookupPreferredNamespacePrefix(tag.Space)
		if prefix == "" {
			// They should notify us of the unknown namespace so that we
			// can register it and they can handle it properly.
			prefix = "?"
		}

		parts[i] = fmt.Sprintf("[%s]%s", prefix, tag.Local)
	}

	return parts
}

func (xpn XmpPropertyName) String() string {
	parts := xpn.Parts()
	return strings.Join(parts, ".")
}

type XmpPropertyIndex struct {
	subindices map[string]*XmpPropertyIndex
	leaves     map[string][]interface{}
}

func newXmpPropertyIndex() *XmpPropertyIndex {
	subindices := make(map[string]*XmpPropertyIndex)
	leaves := make(map[string][]interface{})

	xpi := &XmpPropertyIndex{
		subindices: subindices,
		leaves:     leaves,
	}

	return xpi
}

func (xpi *XmpPropertyIndex) add(name XmpPropertyName, value interface{}) {
	currentNodeName := name[0]
	currentNodeNamePhrase := currentNodeName.String()

	if len(name) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			subindex = newXmpPropertyIndex()
		}

		subindex.add(name[1:], value)

		if found == false {
			xpi.subindices[currentNodeNamePhrase] = subindex
		}
	} else {

		// TODO(dustin): !! Remember, repetition is allowed. Each value should be a slice so get() can return a slice.

		// if _, found := xpi.leaves[currentNodeNamePhrase]; found == true {
		// 	log.Panicf("property value set more than once: [%s]", currentNodeNamePhrase)
		// }

		if currentLeaves, found := xpi.leaves[currentNodeNamePhrase]; found == true {
			xpi.leaves[currentNodeNamePhrase] = append(currentLeaves, value)
		} else {
			xpi.leaves[currentNodeNamePhrase] = []interface{}{value}
		}
	}
}

func (xpi *XmpPropertyIndex) get(namePhraseSlice []string) (results []interface{}, err error) {
	currentNodeNamePhrase := namePhraseSlice[0]

	if len(namePhraseSlice) > 1 {
		if subindex, found := xpi.subindices[currentNodeNamePhrase]; found == false {
			return nil, ErrPropertyNotFound
		} else {
			values, err := subindex.get(namePhraseSlice[1:])
			if err != nil {
				if err == ErrPropertyNotFound {
					return nil, err
				}

				log.Panic(err)
			}

			return values, nil
		}
	}

	// If we get here, we are expecting to find a leaf-node.

	if values, found := xpi.leaves[currentNodeNamePhrase]; found == true {
		return values, nil
	}

	return nil, ErrPropertyNotFound
}

func (xpi *XmpPropertyIndex) dump(prefix []string) {
	for name, subindex := range xpi.subindices {
		subindex.dump(append(prefix, name))
	}

	for name, values := range xpi.leaves {
		fqName := append(prefix, name)
		fqNamePhrase := strings.Join(fqName, ".")

		for _, value := range values {
			fmt.Printf("%s = [%s]\n", fqNamePhrase, value)
		}
	}
}

func (xpi *XmpPropertyIndex) Dump() {
	xpi.dump([]string{})
}

// Parse parses the XMP document.
func (xp *Parser) Parse() (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	xpi := newXmpPropertyIndex()

	for {
		token, err := xp.xd.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Panic(err)
		}

		// We do our last-token management here since many of the conditionals
		// below will skip and continue.
		lastToken := xp.lastToken
		xp.lastToken = token

		switch t := token.(type) {
		case xml.StartElement:
			xp.lastCharData = nil

			if t.Name == rdfTag {
				if xp.rdfIsOpen == true {
					log.Panicf("RDF is already open")
				}

				xp.rdfIsOpen = true

				continue
			} else if t.Name == rdfDescriptionTag {
				if xp.rdfDescriptionIsOpen == true {
					log.Panicf("RDF description is already open")
				}

				xp.rdfDescriptionIsOpen = true

				continue
			}

			xp.stack = append(xp.stack, XmlName(t.Name))
		case xml.EndElement:
			if t.Name == rdfTag {
				if xp.rdfIsOpen == false {
					log.Panicf("RDF is not open")
				}

				xp.rdfIsOpen = false

				continue
			} else if t.Name == rdfDescriptionTag {
				if xp.rdfDescriptionIsOpen == false {
					log.Panicf("RDF description is not open")
				}

				xp.rdfDescriptionIsOpen = false

				continue
			}

			// Process any stash char-data. Since this is cleared whenever we
			// encounter a start-tag, this tells us that we were a leaf/scalar
			// node.
			if xp.lastCharData != nil {
				xpn := XmpPropertyName(xp.stack)

				// TODO(dustin): !! We still need to parse the values to proper types.
				valuePhrase := strings.Trim(string(*xp.lastCharData), "\r\n\t ")

				xpi.add(xpn, valuePhrase)

				xp.lastCharData = nil
			}

			// Go already validates that the tags are balanced.
			xp.stack = xp.stack[:len(xp.stack)-1]

		case xml.CharData:

			// Scope any intermediate nodes that we don't care about (like "
			// xmpmeta"-space tags).
			if xp.rdfDescriptionIsOpen == false {
				continue
			}

			// Ignore any character data between adjacent closing-tags. We only
			// want character-data between an open and a close tag (leaf/scalar
			// nodes).
			if _, ok := lastToken.(xml.EndElement); ok == true {
				continue
			}

			// We only care about char-data for leaf nodes. So, we'll
			// temporarily stash it, we'll always clear stashed char-data when
			// we encounter a start-tag, and then process any stashed data when
			// we encounter a stop-tag.
			value := string(t)
			xp.lastCharData = &value
		case xml.ProcInst:
			if t.Target != "xpacket" {
				continue
			}

			fragment := string(t.Inst)
			parts := strings.Split(fragment, " ")

			if xp.packetIsOpen == false {
				xp.packetIsOpen = true

				foundBegin := false
				foundId := false

				for _, part := range parts {
					raa := rawAttributeAssignment(part)
					name, value := raa.Parse()

					if name == "begin" {
						foundBegin = true

						if len(value) == 0 {
							// NOTE(dustin): Currently clearing this, though we might later recommend that the user use a new struct to process each subsequence xpacket, but keeping this just in case a) they don't listen, or b) we decide to support that method.
							// NOTE(dustin): <-- No current way for the user to know where we are in the stream upon return.
							xp.encodingSequence = nil
						} else {
							bs := bom.Sequence(value)

							_, _, err := bs.GetEncoding()
							log.PanicIf(err)

							xp.encodingSequence = bs
						}
					} else if name == "id" {
						foundId = true

						if value != standardXpacketId {
							log.Panicf("xpacket ID not expected: [%s]", value)
						}
					}
				}

				if foundBegin == false || foundId == false {
					log.Panicf("'begin' or 'id' attributes of xpacket tag missing")
				}
			} else {
				// We don't currently do anything on the closing tag.

				xp.packetIsOpen = false
			}
		}

	}

	return nil
}
