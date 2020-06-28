package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestTextFieldType_GetValueParser(t *testing.T) {
	tft := TextFieldType{}

	s := "test_string"

	actual := tft.GetValueParser(s)
	expected := TextFieldValue{raw: s}

	if actual != expected {
		t.Fatalf("Value parser not correct.")
	}
}

func TestTextFieldValue_Parse(t *testing.T) {
	tfv := TextFieldValue{
		raw: "abc",
	}

	parsed, err := tfv.Parse()
	log.PanicIf(err)

	if parsed != "abc" {
		t.Fatalf("Value not parsed correct: [%v]", parsed)
	}
}
