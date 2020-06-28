package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestProperNameFieldType_GetValueParser(t *testing.T) {
	pnft := ProperNameFieldType{}
	scp := pnft.GetValueParser("test_text")

	pnfv := scp.(ProperNameFieldValue)

	parsed, err := pnfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
