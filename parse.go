package xmp

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"encoding/binary"
	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/type"
)

var (
	parseLogger = log.NewLogger("xmp.parse")
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

	// nameStack is a stack comprised of xml.Name structs.
	nameStack []XmlName

	// stringStack is a stack comprised of simple names, used for matching.
	stringStack []string

	lastCharData *string
	lastToken    xml.Token

	unknownNamespaces map[string]struct{}
}

// NewParser returns a new Parser struct.
func NewParser(r io.Reader) *Parser {
	xd := xml.NewDecoder(r)

	nameStack := make([]XmlName, 0)
	stringStack := make([]string, 0)

	unknownNamespaces := make(map[string]struct{})

	return &Parser{
		xd:                xd,
		nameStack:         nameStack,
		stringStack:       stringStack,
		unknownNamespaces: unknownNamespaces,
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

	nodeName := XmlName(t.Name)
	xp.nameStack = append(xp.nameStack, nodeName)
	xp.stringStack = append(xp.stringStack, nodeName.String())

	for _, attribute := range t.Attr {
		namespaceUri := attribute.Name.Space
		localName := attribute.Name.Local
		rawValue := attribute.Value

		namespace, err := xmpnamespace.Get(namespaceUri)
		if err != nil {
			if err == xmpnamespace.ErrNamespaceNotFound {
				if _, found := xp.unknownNamespaces[namespaceUri]; found == false {
					parseLogger.Warningf(
						nil,
						"Namespace [%s] for attribute [%s] is not known. Skipping.",
						namespaceUri, localName)
				}

				continue
			}

			log.Panic(err)
		}

		parsedValue, err := parseValue(namespace, localName, rawValue)
		if err != nil {
			if err == ErrChildFieldNotFound || err == xmptype.ErrValueNotValid {
				parseLogger.Warningf(
					nil,
					"Could not parse attribute [%s] [%s] value: [%s]",
					namespaceUri, localName, rawValue)

				continue
			}

			log.Panic(err)
		}

		// TODO(dustin): !! Still need to store this value somewhere.
		// fmt.Printf("Parsed ATTRIBUTE [%s] [%s] [%s] -> [%s] [%v]\n", namespaceUri, localName, rawValue, reflect.TypeOf(parsedValue), parsedValue)
		parsedValue = parsedValue
	}

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

	if xp.lastCharData != nil {
		err := xp.parseCharData(t.Name, *xp.lastCharData)
		log.PanicIf(err)

		xp.lastCharData = nil
	}

	// Go already validates that the tags are balanced.
	xp.nameStack = xp.nameStack[:len(xp.nameStack)-1]
	xp.stringStack = xp.stringStack[:len(xp.stringStack)-1]

	return nil
}

// matchNodeName returns whether the given string-slice matches against the end
// of the current string stack (the stringicized names of all of the nodes that
// comprise where we're currently at in the tree).
func (xp *Parser) matchNodeName(queryName []string) bool {
	queryLen := len(queryName)
	if len(xp.stringStack) < queryLen {
		return false
	}

	candidateSuffix := xp.stringStack[len(xp.stringStack)-queryLen:]

	for i, part := range queryName {
		if candidateSuffix[i] != part {
			return false
		}
	}

	return true
}

// parseCharData parses the char-data that exists in leaf-nodes (not in nodes
// that have child-nodes).
func (xp *Parser) parseCharData(nodeName xml.Name, rawValue string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// Handle an array item.

	if nodeName.Local == "li" {
		// TODO(dustin): !! Note that, if/when we're building a list of array items, we should also capture the attributes of the 'li' node as well. This would be necessary with language-alternative nodes, for instance (the data we need in order for that to be of value partially exists as attributes).
		// TODO(dustin): !! Keep in mind that may the encapsulating nodes likely specify which collection type to expect, so we can proactively process them from the top rather than reacting to them from the bottom (like we're doing, here).

		xpn := XmpPropertyName(xp.nameStack[:len(xp.nameStack)-2])

		if xp.matchNodeName([]string{"[rdf]Alt", "[rdf]li"}) == true {
			fmt.Printf("Found alt under [%s]\n", xpn.String())

			// TODO(dustin): !! Finish

		} else if xp.matchNodeName([]string{"[rdf]Seq", "[rdf]li"}) == true {
			fmt.Printf("Found seq under [%s]\n", xpn.String())

			// TODO(dustin): !! Finish

		} else if xp.matchNodeName([]string{"[rdf]Bag", "[rdf]li"}) == true {
			fmt.Printf("Found bag under [%s]\n", xpn.String())

			// TODO(dustin): !! Finish

		}

		return nil
	}

	// Parse a normal node.

	namespaceUri := nodeName.Space
	localName := nodeName.Local

	namespace, err := xmpnamespace.Get(namespaceUri)
	if err != nil {
		if err == xmpnamespace.ErrNamespaceNotFound {
			if _, found := xp.unknownNamespaces[namespaceUri]; found == false {
				parseLogger.Warningf(
					nil,
					"Namespace [%s] for node [%s] with char-data is not known. Skipping.",
					namespaceUri, localName)

				fmt.Printf(
					"Namespace [%s] for node [%s] with char-data is not known. Skipping.\n",
					namespaceUri, localName)
			}

			return nil
		}

		log.Panic(err)
	}

	parsedValue, err := parseValue(namespace, localName, rawValue)
	if err != nil {
		if err == ErrChildFieldNotFound || err == xmptype.ErrValueNotValid {
			parseLogger.Warningf(
				nil,
				"Could not parse char-data under node [%s] [%s] value: [%s]",
				namespaceUri, localName, rawValue)

			fmt.Printf(
				"Could not parse char-data under node [%s] [%s] value: [%s]\n",
				namespaceUri, localName, rawValue)

			return nil
		}

		log.Panic(err)
	}

	// TODO(dustin): !! Still need to store this value somewhere.
	fmt.Printf("Parsed CHAR-DATA [%s] [%s] [%s] -> [%s] [%v]\n", namespaceUri, localName, rawValue, reflect.TypeOf(parsedValue), parsedValue)
	parsedValue = parsedValue

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

	xpi = newXmpPropertyIndex(XmlName{})

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
