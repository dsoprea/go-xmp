package xmp

import (
	"io"
	"strings"

	"encoding/binary"
	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"
)

const (
	rdfXmlNamespace = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
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
func (raa rawAttributeAssignment) parse() (string, string) {
	// At least one character in the name, the equals, and the quote characters
	// on the value.
	if len(raa) < 4 {
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
	bomEncoding  bom.Encoding
	bomByteOrder binary.ByteOrder

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

func (xp *Parser) parseStartElementToken(xpi *XmpPropertyIndex, t xml.StartElement) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	xp.lastCharData = nil

	if t.Name == rdfTag {
		if xp.rdfIsOpen == true {
			log.Panicf("RDF is already open")
		}

		xp.rdfIsOpen = true

		return nil
	} else if t.Name == rdfDescriptionTag {
		if xp.rdfDescriptionIsOpen == true {
			log.Panicf("RDF description is already open")
		}

		xp.rdfDescriptionIsOpen = true

		return nil
	}

	xp.stack = append(xp.stack, XmlName(t.Name))

	return nil
}

func (xp *Parser) parseEndElementToken(xpi *XmpPropertyIndex, t xml.EndElement) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if t.Name == rdfTag {
		if xp.rdfIsOpen == false {
			log.Panicf("RDF is not open")
		}

		xp.rdfIsOpen = false

		return nil
	} else if t.Name == rdfDescriptionTag {
		if xp.rdfDescriptionIsOpen == false {
			log.Panicf("RDF description is not open")
		}

		xp.rdfDescriptionIsOpen = false

		return nil
	}

	if xp.rdfDescriptionIsOpen == false {
		return nil
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

	return nil
}

func (xp *Parser) parseCharDataToken(xpi *XmpPropertyIndex, t xml.CharData, lastToken xml.Token) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// Scope any intermediate nodes that we don't care about (like "
	// xmpmeta"-space tags).
	if xp.rdfDescriptionIsOpen == false {
		return nil
	}

	// Ignore any character data between adjacent closing-tags. We only
	// want character-data between an open and a close tag (leaf/scalar
	// nodes).
	if _, ok := lastToken.(xml.EndElement); ok == true {
		return nil
	}

	// We only care about char-data for leaf nodes. So, we'll
	// temporarily stash it, we'll always clear stashed char-data when
	// we encounter a start-tag, and then process any stashed data when
	// we encounter a stop-tag.
	value := string(t)
	xp.lastCharData = &value

	return nil
}

func (xp *Parser) parseProcInstToken(xpi *XmpPropertyIndex, t xml.ProcInst) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if t.Target != "xpacket" {
		return nil
	}

	fragment := string(t.Inst)
	parts := strings.Split(fragment, " ")

	if xp.packetIsOpen == false {
		xp.packetIsOpen = true

		foundBegin := false
		foundId := false

		for _, part := range parts {
			raa := rawAttributeAssignment(part)
			name, value := raa.parse()

			if name == "begin" {
				foundBegin = true

				if len(value) == 0 {
					// NOTE(dustin): Currently clearing this, though we might later recommend that the user use a new struct to process each subsequence xpacket, but keeping this just in case a) they don't listen, or b) we decide to support that method.
					// NOTE(dustin): <-- No current way for the user to know where we are in the stream upon return.
					xp.bomEncoding = 0
					xp.bomByteOrder = nil
				} else {
					encoding, byteOrder, err := bom.GetEncoding([]byte(value))
					log.PanicIf(err)

					xp.bomEncoding = encoding
					xp.bomByteOrder = byteOrder
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

	return nil
}

func (xp *Parser) parseToken(xpi *XmpPropertyIndex, token xml.Token) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// We do our last-token management here since many of the conditionals
	// below will skip and continue.
	lastToken := xp.lastToken
	xp.lastToken = token

	switch t := token.(type) {
	case xml.StartElement:
		err := xp.parseStartElementToken(xpi, t)
		log.PanicIf(err)

		return nil

	case xml.EndElement:
		err := xp.parseEndElementToken(xpi, t)
		log.PanicIf(err)

		return nil

	case xml.CharData:

		err := xp.parseCharDataToken(xpi, t, lastToken)
		log.PanicIf(err)

		return nil

	case xml.ProcInst:
		err := xp.parseProcInstToken(xpi, t)
		log.PanicIf(err)

		return nil
	}

	return nil
}

// Parse parses the XMP document.
func (xp *Parser) Parse() (xpi *XmpPropertyIndex, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	xpi = newXmpPropertyIndex()

	for {
		token, err := xp.xd.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Panic(err)
		}

		err = xp.parseToken(xpi, token)
		log.PanicIf(err)
	}

	return xpi, nil
}
