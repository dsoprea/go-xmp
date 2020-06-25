package xmpregistry

import (
	"reflect"
	"testing"
)

func TestXmlName_String_Known(t *testing.T) {
	name := XmlName{
		Space: "http://ns.adobe.com/pdf/1.3/",
		Local: "bb",
	}

	if name.String() != "[pdf]bb" {
		t.Fatalf("String not correct: [%s]", name.String())
	}
}

func TestXmlName_String_NotKnown(t *testing.T) {
	name := XmlName{
		Space: "aa",
		Local: "bb",
	}

	if name.String() != "[?]bb" {
		t.Fatalf("String not correct: [%s]", name.String())
	}
}

func TestXmpPropertyName_Parts(t *testing.T) {
	name := XmpPropertyName{
		{"http://ns.adobe.com/pdf/1.3/", "aa"},
		{"http://purl.org/dc/elements/1.1/", "bb"},
	}

	actual := name.Parts()

	expected := []string{
		"[pdf]aa",
		"[dc]bb",
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Parts not correct: %v\n", actual)
	}
}

func TestXmpPropertyName_String(t *testing.T) {
	name := XmpPropertyName{
		{"http://ns.adobe.com/pdf/1.3/", "aa"},
		{"http://purl.org/dc/elements/1.1/", "bb"},
	}

	namePhrase := name.String()

	if namePhrase != "[pdf]aa.[dc]bb" {
		t.Fatalf("String not correct: %s\n", namePhrase)
	}
}
