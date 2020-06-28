package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestBooleanFieldType_GetValueParser_True(t *testing.T) {
	bft := BooleanFieldType{}
	scp := bft.GetValueParser("True")

	bfv := scp.(BooleanFieldValue)

	parsed, err := bfv.Parse()
	log.PanicIf(err)

	if parsed != true {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}

func TestBooleanFieldType_GetValueParser_False(t *testing.T) {
	bft := BooleanFieldType{}
	scp := bft.GetValueParser("False")

	bfv := scp.(BooleanFieldValue)

	parsed, err := bfv.Parse()
	log.PanicIf(err)

	if parsed != false {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
