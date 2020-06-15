package xmp

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestParser_Parse(t *testing.T) {
	data := GetTestData()
	b := bytes.NewBuffer(data)
	xp := NewParser(b)

	xpi, err := xp.Parse()
	log.PanicIf(err)

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
	} else if len(p.stack) != 0 {
		t.Fatalf("Stack not initialized or not empty.")
	}
}
