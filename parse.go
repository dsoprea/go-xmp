package xmp

import (
	"io"
	"strings"

	"encoding/binary"
	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

var (
	parseLogger = log.NewLogger("xmp.parse")
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

// EmbeddedArray describes an array found in another array.
type EmbeddedArray struct {
	name      xml.Name
	collected []interface{}
}

// Parser parses an XMP document.
type Parser struct {
	xd *xml.Decoder

	// TODO(dustin): !! Add an accessor for this. Investigate whether we even have this value in our test-data.
	bomEncoding  bom.Encoding
	bomByteOrder binary.ByteOrder

	packetIsOpen         bool
	rdfIsOpen            bool
	rdfDescriptionIsOpen bool

	// nameStack is a stack comprised of xml.Name structs.
	nameStack []xmpregistry.XmlName

	lastCharData *string
	lastToken    xml.Token

	unfinishedArrayLayers [][]interface{}
}

// NewParser returns a new Parser struct.
func NewParser(r io.Reader) *Parser {
	xd := xml.NewDecoder(r)

	nameStack := make([]xmpregistry.XmlName, 0)

	unfinishedArrayLayers := make([][]interface{}, 0)

	return &Parser{
		xd:                    xd,
		nameStack:             nameStack,
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

	nodeNamespace, err := xmpregistry.Get(nodeNamespaceUri)
	if err != nil {
		if err == xmpregistry.ErrNamespaceNotFound {
			return false, nil
		}

		log.Panic(err)
	}

	flag, err = xmptype.IsArrayType(nodeNamespace, nodeLocalName)
	if err != nil {
		if err == xmptype.ErrChildFieldNotFound {
			return false, nil
		}

		log.Panic(err)
	}

	return flag, nil
}

func (xp *Parser) isInArray() bool {

	// TODO(dustin): Add test

	return len(xp.unfinishedArrayLayers) > 0
}

func (xp *Parser) collectForCurrentArray(element interface{}) {

	// TODO(dustin): Add test

	if len(xp.unfinishedArrayLayers) == 0 {
		log.Panicf("no array currently open")
	}

	currentLayerNumber := len(xp.unfinishedArrayLayers) - 1

	xp.unfinishedArrayLayers[currentLayerNumber] =
		append(
			xp.unfinishedArrayLayers[currentLayerNumber],
			element)
}

func (xp *Parser) parseStartElementToken(xpi *XmpPropertyIndex, t xml.StartElement) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Expand the unit-tests.

	xp.lastCharData = nil

	// Note that there are other tags that may be outside of the RDF/description
	// nodes, but that the XMP data is *within* them.

	if t.Name == xmpnamespace.RdfTag {
		if xp.rdfIsOpen == true {
			log.Panicf("RDF is already open")
		}

		xp.rdfIsOpen = true

		return nil
	} else if t.Name == xmpnamespace.RdfDescriptionTag {
		if xp.rdfDescriptionIsOpen == true {
			log.Panicf("RDF description is already open")
		}

		xp.rdfDescriptionIsOpen = true

		return nil
	}

	nodeName := xmpregistry.XmlName(t.Name)
	xp.nameStack = append(xp.nameStack, nodeName)

	// Determine if the current node is known to have an array underneath it.

	isArray, err := xp.isArrayNode(t.Name)
	log.PanicIf(err)

	if isArray == true {
		// We've encountered a new array. None of the known RDF array types has
		// attributes on the start-tag, so we won't gather them.

		xp.unfinishedArrayLayers = append(xp.unfinishedArrayLayers, make([]interface{}, 0))
	} else if xp.isInArray() == true {
		// We've not encountered a new array but are currently inside a higher
		// one. Append the current node to it. Since any attributes may be
		// encapsulated, we defer to our array-management to extract them.

		// Any attributes will be extracted by our array management.

		xp.collectForCurrentArray(t)
	} else {
		// This is a simple leaf node that is not an array nor underneath an
		// array. If it has tangible attributes, we'll represent it as a
		// complex-node type and push to the index.

		attributes, err := xmptype.ParseAttributes(t)
		log.PanicIf(err)

		if len(attributes) > 0 {
			xpn := xmpregistry.XmpPropertyName(xp.nameStack)

			err := xpi.addComplexValue(xpn, attributes)
			log.PanicIf(err)
		}
	}

	return nil
}

func (xp *Parser) parseEndElementToken(xpi *XmpPropertyIndex, t xml.EndElement) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Expand the unit-tests.

	isRootRdfTag, err := xp.processNodeCloseRootRdfTags(xpi, t.Name)
	log.PanicIf(err)

	if isRootRdfTag == true {
		// We're sitting on an XMP document RDF boundary tag. Skip.

		return nil
	} else if xp.rdfDescriptionIsOpen == false {
		// We're outside of the XMP space.

		return nil
	}

	if _, ok := xp.lastToken.(xml.StartElement); ok == true {
		// If a self-closing tag, the parser sublimes directly from the start-
		// tag to the end-tag and entirely skip the char-data token.
		// Artificially add it so that we can uniformly process array (where we
		// need to enumerate all of the inner nodes and it'll help if each item
		// is three tokens: start, char-data, end).

		emptyCharData := ""
		xp.lastCharData = &emptyCharData
	}

	if xp.lastCharData != nil {
		err := xp.processNodeCloseCharData(xpi, t.Name, *xp.lastCharData)
		log.PanicIf(err)

		xp.lastCharData = nil
	}

	// Process the array-end if this node was known to be an array.

	nodeNamespaceUri := t.Name.Space
	nodeLocalName := t.Name.Local

	// If the current node is an array-type, get the struct that represents it.

	var arrayType xmptype.ArrayFieldType

	if nodeNamespace, err := xmpregistry.Get(nodeNamespaceUri); err == nil {
		if ft, found := nodeNamespace.Fields[nodeLocalName]; found == true {
			if t, ok := ft.(xmptype.ArrayFieldType); ok == true {
				arrayType = t
			}
		}
	} else if err != xmpregistry.ErrNamespaceNotFound {
		log.Panic(err)
	}

	if arrayType != nil {
		// We're closing an array.

		err := xp.processNodeCloseArrayClose(xpi, t.Name, arrayType)
		log.PanicIf(err)
	} else if xp.isInArray() == true {
		// We've not closed an array but are currently inside a higher one.
		// Append the current node to it.

		xp.collectForCurrentArray(t)
	}

	// Go already validates that the tags are balanced.
	xp.nameStack = xp.nameStack[:len(xp.nameStack)-1]

	return nil
}

func (xp *Parser) processNodeCloseRootRdfTags(xpi *XmpPropertyIndex, nodeName xml.Name) (hit bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	if nodeName == xmpnamespace.RdfTag {
		if xp.rdfIsOpen == false {
			log.Panicf("RDF is not open")
		} else if xp.rdfDescriptionIsOpen != false {
			log.Panicf("RDF document node closed before description node closed")
		}

		xp.rdfIsOpen = false

		return true, nil
	} else if nodeName == xmpnamespace.RdfDescriptionTag {
		if xp.rdfDescriptionIsOpen == false {
			log.Panicf("RDF description is not open")
		}

		xp.rdfDescriptionIsOpen = false

		return true, nil
	}

	return false, nil
}

func (xp *Parser) processNodeCloseCharData(xpi *XmpPropertyIndex, nodeName xml.Name, charData string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	charData = strings.Trim(charData, " \t\r\n")

	if nodeName == xmpnamespace.RdfLiTag {
		if xp.isInArray() == true {
			// This is an array item. Since the value is encapsulated, it will
			// need to be extracted before parsing. So, we just push the raw
			// element and defer to our array management to do that.

			xp.collectForCurrentArray(charData)
		} else {
			// We encountered an array item under a namespace that we don't
			// recognize.

			xpn := xmpregistry.XmpPropertyName(xp.nameStack)

			parseLogger.Warningf(
				nil,
				"We encountered an array item that wasn't in an array, likely because it is in an unregistered namespace (or because its parent types should have an array type and do not): [%s]",
				xpn)
		}
	} else {
		err := xp.parseCharData(xpi, nodeName, charData)
		log.PanicIf(err)
	}

	return nil
}

func (xp *Parser) processNodeCloseArrayClose(xpi *XmpPropertyIndex, nodeName xml.Name, arrayType xmptype.ArrayFieldType) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	currentUnfinishedLayerNumber := len(xp.unfinishedArrayLayers) - 1
	finishedArray := xp.unfinishedArrayLayers[currentUnfinishedLayerNumber]

	ea := EmbeddedArray{
		name:      nodeName,
		collected: finishedArray,
	}

	xp.unfinishedArrayLayers = xp.unfinishedArrayLayers[:currentUnfinishedLayerNumber]

	if xp.isInArray() == true {
		newUnfinishedLayerNumber := len(xp.unfinishedArrayLayers) - 1
		xp.unfinishedArrayLayers[newUnfinishedLayerNumber] = append(
			xp.unfinishedArrayLayers[newUnfinishedLayerNumber],
			ea)
	}

	xpn := xmpregistry.XmpPropertyName(xp.nameStack)

	wrappedArray := arrayType.New(xpn, finishedArray)

	err = xpi.addArrayValue(xpn, wrappedArray)
	log.PanicIf(err)

	return nil
}

// parseCharData parses the char-data that exists in leaf-nodes (not in nodes
// that have child-nodes).
func (xp *Parser) parseCharData(xpi *XmpPropertyIndex, nodeName xml.Name, rawValue string) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Expand the unit-tests.

	xpn := xmpregistry.XmpPropertyName(xp.nameStack)

	// Parse a normal node.

	namespaceUri := nodeName.Space
	localName := nodeName.Local

	namespace, err := xmpregistry.Get(namespaceUri)
	if err != nil {
		if err == xmpregistry.ErrNamespaceNotFound {
			return nil
		}

		log.Panic(err)
	}

	// Since we ensure that all leaf nodes have char-data we'll periodically end-
	// up with char-data that is empty for nodes in namespaces that don't
	// identify that node with a type. In this case, just silently skip.
	if rawValue == "" && namespace.Fields[localName] == nil {
		return nil
	}

	parsedValue, err := xmptype.ParseValue(namespace, localName, rawValue)
	if err != nil {
		if err == xmptype.ErrChildFieldNotFound {
			parseLogger.Warningf(
				nil,
				"Could not parse char-data under node [%s] [%s] value (field not found): [%s]",
				namespaceUri, localName, rawValue)

			return nil
		} else if err == xmptype.ErrValueNotValid {
			parseLogger.Warningf(
				nil,
				"Could not parse char-data under node [%s] [%s] value (value not valid): [%s]",
				namespaceUri, localName, rawValue)

			return nil
		}

		log.Panic(err)
	}

	if xp.isInArray() == true {
		// We're currently collecting items for an array. Append the char-data
		// to the collector slice.

		xp.collectForCurrentArray(parsedValue)
	} else {
		// This is a non-array-item value-node.

		err := xpi.addScalarValue(xpn, parsedValue)
		log.PanicIf(err)
	}

	return nil
}

// parseCharDataToken processes a char-data token encountered by the parser.
// Note that tangible processing is deferred until we see an end-token in order
// to ignore char-data on nodes that have child nodes ("tree nodes" versus "
// value nodes", both of which will have char-data, which is likely always blank
// for tree nodes).
func (xp *Parser) parseCharDataToken(xpi *XmpPropertyIndex, t xml.CharData, lastToken xml.Token) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Expand the unit-tests.

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

	switch t := token.(type) {
	case xml.StartElement:
		err := xp.parseStartElementToken(xpi, t)
		log.PanicIf(err)

	case xml.EndElement:
		err := xp.parseEndElementToken(xpi, t)
		log.PanicIf(err)

	case xml.CharData:
		err := xp.parseCharDataToken(xpi, t, xp.lastToken)
		log.PanicIf(err)

	case xml.ProcInst:
		err := xp.parseProcInstToken(xpi, t)
		log.PanicIf(err)
	}

	xp.lastToken = token

	return nil
}

// Parse parses the XMP document.
func (xp *Parser) Parse() (xpi *XmpPropertyIndex, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	xpi = newXmpPropertyIndex(xmpregistry.XmlName{})

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
