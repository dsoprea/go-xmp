package xmp

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestNewParser(t *testing.T) {
	// Not much else we can test at this juncture.
	NewParser(nil)
}

func TestParser_Parse(t *testing.T) {
	data := GetTestData()
	b := bytes.NewBuffer(data)
	xp := NewParser(b)

	xpi, err := xp.Parse()
	log.PanicIf(err)

	actual, err := xpi.get([]string{"[x]xmpmeta", "[claro]Logging", "[rdf]Seq", "[rdf]li"})
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
