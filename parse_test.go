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

	name := xml.Name{
		Space: "aa",
		Local: "bb",
	}

	err := xp.parseStartElementToken(nil, xml.StartElement{Name: name})
	log.PanicIf(err)

	expected := []xmpregistry.XmlName{xmpregistry.XmlName(name)}

	if reflect.DeepEqual(xp.nameStack, expected) != true {
		t.Fatalf("Stack not correct.")
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
