package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestPartFieldType_GetValueParser(t *testing.T) {
	pft := PartFieldType{}
	scp := pft.GetValueParser("test_text")

	pfv := scp.(PartFieldValue)

	parsed, err := pfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
