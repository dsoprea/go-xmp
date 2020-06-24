package xmp

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-unicode-byteorder"
)

func TestParser_Parse(t *testing.T) {
	data := GetTestData()
	b := bytes.NewBuffer(data)
	xp := NewParser(b)

	xpi, err := xp.Parse()
	log.PanicIf(err)

	fmt.Printf("Dumping\n")

	xpi.Dump()

	return

	actual, err := xpi.Get([]string{"[x]xmpmeta", "[claro]Logging", "[rdf]Seq", "[rdf]li"})
	log.PanicIf(err)

	expected := []interface{}{
		"20141001 11:47:23 Channel 'WebRGB_Crop' processing file: 36253.jpg (36253.xml)",
		"20141001 11:47:23 IMGINFO size 6.103851 MPixel 150 DPI RGB, used assumed profile sRGB IEC61966-2.1",
		"20141001 11:47:24 WARNING invalid CROP parameter: Height limited:2969.814464569092>2953",
		"20141001 11:47:24 CROP with: xCropStart=27 yCropStart=25 width=2022 height=2927",
		"20141001 11:47:25 ADJUST : resized with factor 0.15921041 (dpi set to 300.0):  width =322  height=466  size=150052 pixels",
		"20141001 11:47:25 ADJUST : set 300.0 DPI",
		"20141001 11:47:26 IMPROVE Sharpening 100% radius:0.7849731545638932 threshold: 0",
		"20141001 11:47:26 CONVERT EMBED profile  sRGB.icc (RelativeColorimetric BPC rendering intent), Out Of Gamut in conversion: 0%",
	}

	if reflect.DeepEqual(actual, expected) != true {
		fmt.Printf("Actual:\n")
		fmt.Printf("\n")

		for _, line := range actual {
			fmt.Printf("[%s]\n", line)
		}

		fmt.Printf("\n")

		fmt.Printf("Expected:\n")
		fmt.Printf("\n")

		for _, line := range expected {
			fmt.Printf("[%s]\n", line)
		}

		fmt.Printf("\n")

		t.Fatalf("Values not correct.")
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

	expected := []XmlName{XmlName(name)}

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
	xp := NewParser(nil)

	xp.rdfIsOpen = true
	xp.rdfDescriptionIsOpen = true

	name := xml.Name{
		Space: "aa",
		Local: "bb",
	}

	xpi := newXmpPropertyIndex(XmlName{})

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
		t.Fatalf("XPI should have one item after close.")
	}

	results, err := xpi.Get([]string{"[?]bb"})
	log.PanicIf(err)

	expected := []interface{}{
		"some data",
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
