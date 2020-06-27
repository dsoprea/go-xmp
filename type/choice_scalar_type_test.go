package xmptype

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

// Open-choices tests

func TestNewOpenChoiceFieldValue(t *testing.T) {
	choices := []string{"aa", "bb"}

	ocfv := NewOpenChoiceFieldValue("def", choices)

	if ocfv.Raw() != "def" {
		t.Fatalf("Raw value not correct.")
	} else if reflect.DeepEqual(ocfv.choices, choices) != true {
		t.Fatalf("Choices not correct.")
	}
}

func TestOpenChoiceFieldValue_Parse_Hit(t *testing.T) {
	ocfv := OpenChoiceFieldValue{
		raw: "bb",
	}

	parsed, err := ocfv.Parse()
	log.PanicIf(err)

	if parsed != "bb" {
		t.Fatalf("Expected result same as argument.")
	}
}

func TestOpenChoiceFieldValue_Parse_Miss(t *testing.T) {
	ocfv := OpenChoiceFieldValue{
		raw: "cc",
	}

	parsed, err := ocfv.Parse()
	log.PanicIf(err)

	if parsed != "cc" {
		t.Fatalf("Expected result same as argument.")
	}
}

func TestOpenChoiceFieldValue_Raw(t *testing.T) {
	ocfv := OpenChoiceFieldValue{
		raw: "abc",
	}

	if ocfv.Raw() != "abc" {
		t.Fatalf("Returned value not correct: [%s]", ocfv.Raw())
	}
}

type testOpenChoicesType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to parse
// a specific string. It must be overridden for each XMP type that needs it.
func (testOpenChoicesType) GetValueParser(raw string) ScalarValueParser {
	return OpenChoiceFieldValue{
		raw: raw,
		choices: []string{
			"aa",
			"bb",
		},
	}
}

func TestOpenChoiceFieldType_GetValueParser_Hit(t *testing.T) {
	ocft := testOpenChoicesType{}
	svp := ocft.GetValueParser("aa")

	parsed, err := svp.Parse()
	log.PanicIf(err)

	if parsed != "aa" {
		t.Fatalf("Parse failed: [%s]", parsed)
	}
}

func TestOpenChoiceFieldType_GetValueParser_Miss(t *testing.T) {
	ocft := testOpenChoicesType{}
	svp := ocft.GetValueParser("cc")

	parsed, err := svp.Parse()
	log.PanicIf(err)

	if parsed != "cc" {
		t.Fatalf("Parse failed: [%s]", parsed)
	}
}

func TestOpenChoiceFieldType_GetValueParser_NotImplemented(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if err != ErrChoicesNotOverridden {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected panic if choices not overridden.")
		}
	}()

	ccft := OpenChoiceFieldType{}
	ccft.GetValueParser("")
}

// Closed-choices tests

func TestNewClosedChoiceFieldValue(t *testing.T) {
	choices := []string{"aa", "bb"}
	ccfv := NewClosedChoiceFieldValue("def", choices)

	if ccfv.Raw() != "def" {
		t.Fatalf("Raw value not correct.")
	} else if reflect.DeepEqual(ccfv.choices, choices) != true {
		t.Fatalf("Choices not correct.")
	}
}

func TestClosedChoiceFieldValue_Parse_Hit(t *testing.T) {
	choices := []string{"aa", "bb"}
	ocfv := NewClosedChoiceFieldValue("bb", choices)

	parsed, err := ocfv.Parse()
	log.PanicIf(err)

	if parsed != "bb" {
		t.Fatalf("Expected result same as argument.")
	}
}

func TestClosedChoiceFieldValue_Parse_Miss(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if log.Is(err, ErrValueNotValid) != true {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected value failure.")
		}
	}()

	choices := []string{"aa", "bb"}
	ocfv := NewClosedChoiceFieldValue("cc", choices)

	_, err := ocfv.Parse()
	log.PanicIf(err)
}

func TestClosedChoiceFieldValue_Raw(t *testing.T) {
	ocfv := OpenChoiceFieldValue{
		raw: "abc",
	}

	if ocfv.Raw() != "abc" {
		t.Fatalf("Returned value not correct: [%s]", ocfv.Raw())
	}
}

func TestClosedChoiceFieldType_GetValueParser_NotImplemented(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if err != ErrChoicesNotOverridden {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected panic if choices not overridden.")
		}
	}()

	ccft := ClosedChoiceFieldType{}
	ccft.GetValueParser("")
}

type testClosedChoicesType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to parse
// a specific string. It must be overridden for each XMP type that needs it.
func (testClosedChoicesType) GetValueParser(raw string) ScalarValueParser {
	return ClosedChoiceFieldValue{
		raw: raw,
		choices: []string{
			"aa",
			"bb",
		},
	}
}

func TestClosedChoiceFieldType_GetValueParser_Hit(t *testing.T) {
	ccft := testClosedChoicesType{}

	svp := ccft.GetValueParser("aa")

	parsed, err := svp.Parse()
	log.PanicIf(err)

	if parsed != "aa" {
		t.Fatalf("Parse failed: [%s]", parsed)
	}
}

func TestClosedChoiceFieldType_GetValueParser_Miss(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)
			if log.Is(err, ErrValueNotValid) != true {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected panic if choices not overridden.")
		}
	}()

	ccft := testClosedChoicesType{}
	svp := ccft.GetValueParser("cc")

	_, err := svp.Parse()
	log.PanicIf(err)
}
