package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestGuidFieldType_GetValueParser(t *testing.T) {
	gft := GuidFieldType{}
	scp := gft.GetValueParser("test_text")

	gfv := scp.(GuidFieldValue)

	parsed, err := gfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
