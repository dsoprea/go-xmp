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

type arrayInstance struct {
	name      xml.Name
	collected []interface{}
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

	lastCharData *string
	lastToken    xml.Token

	unknownNamespaces map[string]struct{}

	unfinishedArrayLayers [][]interface{}
}

// NewParser returns a new Parser struct.
func NewParser(r io.Reader) *Parser {
	xd := xml.NewDecoder(r)

	nameStack := make([]XmlName, 0)

	unknownNamespaces := make(map[string]struct{})

	unfinishedArrayLayers := make([][]interface{}, 0)

	return &Parser{
		xd:                xd,
		nameStack:         nameStack,
		unknownNamespaces: unknownNamespaces,

		unfinishedArrayLayers: unfinishedArrayLayers,
	}
}

func (xp *Parser) isArrayNode(name xml.Name) (flag bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	nodeNamespaceUri := name.Space
	nodeLocalName := name.Local

	nodeNamespace, err := xmpnamespace.Get(nodeNamespaceUri)
	if err != nil {
		if err == xmpnamespace.ErrNamespaceNotFound {
			if _, found := xp.unknownNamespaces[nodeNamespaceUri]; found == false {
				parseLogger.Warningf(
					nil,
					"Namespace [%s] for node [%s] is not known. Skipping array check.",
					nodeNamespaceUri, nodeLocalName)
			}

			return false, nil
		} else {
			log.Panic(err)
		}
	}

	flag, err = isArrayType(nodeNamespace, nodeLocalName)
	if err != nil {
		if err == ErrChildFieldNotFound {
			parseLogger.Warningf(
				nil,
				"Namespace [%s] for node [%s] is not known to have child [%s]. Skipping array check.",
				nodeNamespaceUri, nodeLocalName, nodeLocalName)

			return false, nil
		} else {
			log.Panic(err)
		}
	}

	return true, nil
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

	// Try to lookup and parse attributes.

	for _, attribute := range t.Attr {
		attributeNamespaceUri := attribute.Name.Space
		attributeLocalName := attribute.Name.Local
		attributeRawValue := attribute.Value

		attributeNamespace, err := xmpnamespace.Get(attributeNamespaceUri)
		if err != nil {
			if err == xmpnamespace.ErrNamespaceNotFound {
				if _, found := xp.unknownNamespaces[attributeNamespaceUri]; found == false {
					parseLogger.Warningf(
						nil,
						"Namespace [%s] for attribute [%s] is not known. Skipping.",
						attributeNamespaceUri, attributeLocalName)
				}

				continue
			}

			log.Panic(err)
		}

		parsedValue, err := parseValue(attributeNamespace, attributeLocalName, attributeRawValue)
		if err != nil {
			if err == ErrChildFieldNotFound || err == xmptype.ErrValueNotValid {
				parseLogger.Warningf(
					nil,
					"Could not parse attribute [%s] [%s] value: [%s]",
					attributeNamespaceUri, attributeLocalName, attributeRawValue)

				continue
			}

			log.Panic(err)
		}

		// TODO(dustin): !! Still need to store this value somewhere.
		// fmt.Printf("Parsed ATTRIBUTE [%s] [%s] [%s] -> [%s] [%v]\n", namespaceUri, localName, rawValue, reflect.TypeOf(parsedValue), parsedValue)
		parsedValue = parsedValue
	}

	// Determine if the current node is known to have an array underneath it.

	isArray, err := xp.isArrayNode(t.Name)
	log.PanicIf(err)

	if isArray == true {
		// We've encountered a new array.

		xpn := XmpPropertyName(xp.nameStack)
		fmt.Printf("Starting array: %s\n", xpn)

		// TODO(dustin): !! Might want to capture the parsed attributes here.

		xp.unfinishedArrayLayers = append(xp.unfinishedArrayLayers, make([]interface{}, 0))
	} else if len(xp.unfinishedArrayLayers) > 0 {
		// We've not encountered a new array but are currently inside a higher
		// one. Append the current node to it.

		xpn := XmpPropertyName(xp.nameStack)
		fmt.Printf("Collecting within array: %s\n", xpn)

		// TODO(dustin): !! Might want to capture the parsed attributes here.

		currentLayerNumber := len(xp.unfinishedArrayLayers) - 1
		xp.unfinishedArrayLayers[currentLayerNumber] = append(xp.unfinishedArrayLayers[currentLayerNumber], t)
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

	// Process the array-end if this node was known to be an array.

	isArray, err := xp.isArrayNode(t.Name)
	log.PanicIf(err)

	if isArray == true {
		// We've encountered a new array.

		currentUnfinishedLayerNumber := len(xp.unfinishedArrayLayers) - 1

		ai := arrayInstance{
			name:      t.Name,
			collected: xp.unfinishedArrayLayers[currentUnfinishedLayerNumber],
		}

		xp.unfinishedArrayLayers = xp.unfinishedArrayLayers[:currentUnfinishedLayerNumber]

		if len(xp.unfinishedArrayLayers) > 0 {
			newUnfinishedLayerNumber := len(xp.unfinishedArrayLayers) - 1
			xp.unfinishedArrayLayers[newUnfinishedLayerNumber] = append(xp.unfinishedArrayLayers[newUnfinishedLayerNumber], ai)
		}

		xpn := XmpPropertyName(xp.nameStack)
		fmt.Printf("Finished array: %s\n", xpn)

		// TODO(dustin): !! Flatten and finish the array-instance into the index. Note that arrays within arrays will appear as separate instances: The lower arrays are directly indexable as well as able to be found within the items of the higher arrays.

	} else if len(xp.unfinishedArrayLayers) > 0 {
		// We've not closed an array but are currently inside a higher one.
		// Append the current node to it.

		xpn := XmpPropertyName(xp.nameStack)
		fmt.Printf("Collecting within array: %s\n", xpn)

		currentLayerNumber := len(xp.unfinishedArrayLayers) - 1
		xp.unfinishedArrayLayers[currentLayerNumber] = append(xp.unfinishedArrayLayers[currentLayerNumber], t)

		// TODO(dustin): !! We probably want to capture the parsed char-data, here.
	}

	// Go already validates that the tags are balanced.
	xp.nameStack = xp.nameStack[:len(xp.nameStack)-1]

	return nil
}

// parseCharData parses the char-data that exists in leaf-nodes (not in nodes
// that have child-nodes).
func (xp *Parser) parseCharData(nodeName xml.Name, rawValue string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// This is an array item. This is handled in the tag-close handler.
	if nodeName.Local == "li" {
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
