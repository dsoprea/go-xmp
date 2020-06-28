package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestRenditionClassFieldType_GetValueParser(t *testing.T) {
	rcft := RenditionClassFieldType{}
	scp := rcft.GetValueParser("low-res")

	rcfv := scp.(RenditionClassFieldValue)

	parsed, err := rcfv.Parse()
	log.PanicIf(err)

	if parsed != "low-res" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
