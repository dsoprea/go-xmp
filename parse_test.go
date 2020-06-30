package xmp

import (
	"bytes"
	"reflect"
	"testing"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/registry"
)

var (
	xmpLabelName = xml.Name{
		Space: xmpnamespace.XmpUri,
		Local: "Label",
	}
)

func TestParser_Parse_Complex(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)
			t.Fatalf("Test failed.")
		}
	}()

	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.XNamespace)
	xmpregistry.Register(xmpnamespace.XmpMmNamespace)
	xmpregistry.Register(xmpnamespace.StRefNamespace)

	data := GetTestData()
	b := bytes.NewBuffer(data)
	xp := NewParser(b)

	xpi, err := xp.Parse()
	log.PanicIf(err)

	results, err := xpi.Get([]string{"[x]xmpmeta", "[xmpMM]DerivedFrom"})
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("Expected one result: (%d)", len(results))
	}

	cln, ok := results[0].(ComplexLeafNode)
	if ok != true {
		t.Fatalf("Result is not a complex leaf node.")
	}

	value, found := cln.Get(xmpnamespace.StRefUri, "documentID")
	if found != true {
		t.Fatalf("Could not find attribute in result.")
	}

	if value != "xmp.did:146E0D5C4520681181ACE0A302384436" {
		t.Fatalf("Value is not correct: [%s]", value)
	}
}

func TestRawAttributeAssignment_parse_doubleQuotes(t *testing.T) {
	phrase := `aa="bb"`

	name, value := rawAttributeAssignment(phrase).parse()

	if name != "aa" {
		t.Fatalf("Name not correct: [%s]", name)
	} else if value != "bb" {
		t.Fatalf("Value not correct: [%s]", value)
	}
}

func TestRawAttributeAssignment_parse_singleQuotes(t *testing.T) {
	phrase := `aa='bb'`

	name, value := rawAttributeAssignment(phrase).parse()

	if name != "aa" {
		t.Fatalf("Name not correct: [%s]", name)
	} else if value != "bb" {
		t.Fatalf("Value not correct: [%s]", value)
	}
}

func TestRawAttributeAssignment_parse_emptyValue(t *testing.T) {
	phrase := `aa=""`

	name, value := rawAttributeAssignment(phrase).parse()

	if name != "aa" {
		t.Fatalf("Name not correct: [%s]", name)
	} else if value != "" {
		t.Fatalf("Value not correct: [%s]", value)
	}
}

func TestRawAttributeAssignment_parse_invalid(t *testing.T) {
	phrase := `aa=`

	name, value := rawAttributeAssignment(phrase).parse()

	if name != "" {
		t.Fatalf("Name should have been unparseable: [%s]", name)
	} else if value != "" {
		t.Fatalf("Value should have been unparseable: [%s]", value)
	}
}

func TestNewParser(t *testing.T) {
	p := NewParser(nil)

	if p.xd == nil {
		t.Fatalf("XML decoder not assigned.")
	} else if len(p.nameStack) != 0 {
		t.Fatalf("Name stack not initialized or not empty.")
	}
}

func TestParser_parseStartElementToken(t *testing.T) {
	xp := NewParser(nil)

	se := xml.StartElement{Name: xmpnamespace.RdfTag}

	err := xp.parseStartElementToken(nil, se)
	log.PanicIf(err)

	if xp.rdfIsOpen != true {
		t.Fatalf("RDF document did not register as open.")
	}
}

func TestParser_parseEndElementToken(t *testing.T) {
	xp := NewParser(nil)

	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	name := xml.Name{
		Space: "aa",
		Local: "bb",
	}

	err := xp.parseStartElementToken(nil, xml.StartElement{Name: name})
	log.PanicIf(err)

	if len(xp.nameStack) != 1 {
		t.Fatalf("Stack should have one item.")
	}

	err = xp.parseEndElementToken(nil, xml.EndElement{Name: name})
	log.PanicIf(err)

	if len(xp.nameStack) != 0 {
		t.Fatalf("Stack should be empty.")
	}
}

func TestParser_parseCharDataToken(t *testing.T) {
	xp := NewParser(nil)

	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	charData := xml.CharData("some data")

	err := xp.parseCharDataToken(nil, charData, nil)
	log.PanicIf(err)

	if *xp.lastCharData != "some data" {
		t.Fatalf("Character data not stashed correctly: [%s]", *xp.lastCharData)
	}
}

func TestParser_parseEndElementToken_pushToIndex(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.XmpNamespace)

	xp := NewParser(nil)

	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	name := xml.Name{
		Space: xmpnamespace.XmpUri,
		Local: "Label",
	}

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	err := xp.parseStartElementToken(xpi, xml.StartElement{Name: name})
	log.PanicIf(err)

	charData := xml.CharData("some data")

	err = xp.parseCharDataToken(xpi, charData, nil)
	log.PanicIf(err)

	if xpi.Count() != 0 {
		t.Fatalf("XPI should be empty prior to close.")
	}

	err = xp.parseEndElementToken(xpi, xml.EndElement{Name: name})
	log.PanicIf(err)

	if xpi.Count() != 1 {
		t.Fatalf("XPI should have one item after close: (%d)", xpi.Count())
	}

	results, err := xpi.Get([]string{"[xmp]Label"})
	log.PanicIf(err)

	sln := ScalarLeafNode{
		Name:        xmpLabelName,
		ParsedValue: "some data",
	}

	expected := []interface{}{
		sln,
	}

	if reflect.DeepEqual(results, expected) != true {
		t.Fatalf("Results not correct: %v", results)
	}
}

func TestParser_parseProcInstToken(t *testing.T) {
	xp := NewParser(nil)

	err := xp.parseProcInstToken(nil, xml.ProcInst{Target: "xpacket", Inst: []byte("begin=\"\uFEFF\" id=\"W5M0MpCehiHzreSzNTczkc9d\"")})
	log.PanicIf(err)

	if xp.packetIsOpen != true {
		t.Fatalf("Expected packetIsOpen to be true.")
	} else if xp.bomEncoding != bom.Utf8Encoding {
		t.Fatalf("bomEncoding not correct: 0x%x", xp.bomEncoding)
	} else if xp.bomByteOrder != nil {
		t.Fatalf("Byte-order not correct: [%v]", xp.bomByteOrder)
	}
}

func TestParser_parseToken_ProcInst(t *testing.T) {

	xp := NewParser(nil)

	if xp.packetIsOpen != false {
		t.Fatalf("Expected packet to be initially closed.")
	}

	// Test ProcInst.

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	token1 := xml.ProcInst{
		Target: "xpacket",
		Inst:   []byte(`begin="" id="W5M0MpCehiHzreSzNTczkc9d"`),
	}

	err := xp.parseToken(xpi, token1)
	log.PanicIf(err)

	if xp.packetIsOpen != true {
		t.Fatalf("Expected packet to be open.")
	}
}

func TestParser_parseToken_StartElement(t *testing.T) {

	se := xml.StartElement{
		Name: xmpnamespace.RdfTag,
	}

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	xp := NewParser(nil)

	err := xp.parseToken(xpi, se)
	log.PanicIf(err)

	if xp.rdfIsOpen != true {
		t.Fatalf("RDF document did not register as open.")
	}
}

func TestParser_parseToken_EndElement(t *testing.T) {
	ee := xml.EndElement{
		Name: xmpnamespace.RdfTag,
	}

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	xp := NewParser(nil)
	xp.rdfIsOpen = true

	err := xp.parseToken(xpi, ee)
	log.PanicIf(err)

	if xp.rdfIsOpen != false {
		t.Fatalf("RDF document was not closed as expected.")
	}
}

func TestParser_parseToken_CharData(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.XmpNamespace)

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	xp := NewParser(nil)
	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	xp.lastToken = xml.Name{
		Space: xmpnamespace.XmpUri,
		Local: "Label",
	}

	cd := xml.CharData("abcdef")

	err := xp.parseToken(xpi, cd)
	log.PanicIf(err)

	if *xp.lastCharData != "abcdef" {
		t.Fatalf("Stashed last char-data is not correct: [%s]", *xp.lastCharData)
	}
}

func TestParser_parseCharData(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			log.PrintError(err)
			t.Fatalf("Test failed.")
		}
	}()

	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.XmpNamespace)

	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	xp := NewParser(nil)
	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	lastToken := xml.Name{
		Space: xmpnamespace.XmpUri,
		Local: "Label",
	}

	xp.nameStack = append(xp.nameStack, xmpregistry.XmlName(lastToken))

	err := xp.parseCharData(xpi, lastToken, "abcdef")
	log.PanicIf(err)

	results, err := xpi.Get([]string{"[xmp]Label"})
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("Scalar not found.")
	}

	sln := results[0].(ScalarLeafNode)

	if sln.ParsedValue != "abcdef" {
		t.Fatalf("Scalar not correct: [%s]", sln.ParsedValue)
	}
}

func TestParser_isArrayNode_Hit(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.StFntNamespace)

	name := xml.Name{
		Space: xmpnamespace.StFntUri,
		Local: "childFontFiles",
	}

	xp := NewParser(nil)

	flag, err := xp.isArrayNode(name)
	log.PanicIf(err)

	if flag != true {
		t.Fatalf("Expected name to be seen as array.")
	}
}

func TestParser_isArrayNode_Miss(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.StFntNamespace)

	name := xml.Name{
		Space: xmpnamespace.StFntUri,
		Local: "fontFace",
	}

	xp := NewParser(nil)

	flag, err := xp.isArrayNode(name)
	log.PanicIf(err)

	if flag != false {
		t.Fatalf("Expected name to be seen as non-array.")
	}
}

func TestParser_isArrayNode_InvalidName(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpregistry.Register(xmpnamespace.StFntNamespace)

	name := xml.Name{
		Space: xmpnamespace.StFntUri,
		Local: "xyz",
	}

	xp := NewParser(nil)

	found, err := xp.isArrayNode(name)
	log.PanicIf(err)

	if found != false {
		t.Fatalf("Expected invalid child to not be found (and to not be an error).")
	}
}

func TestParser_isArrayNode_UnregisteredNamespace(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	name := xml.Name{
		Space: "invalid/namespace",
		Local: "xyz",
	}

	xp := NewParser(nil)

	found, err := xp.isArrayNode(name)
	log.PanicIf(err)

	if found != false {
		t.Fatalf("Expected isArrayNode to return false if unregistered namespace.")
	}
}
