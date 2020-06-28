package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestUrlFieldType_GetValueParser(t *testing.T) {
	uft := UrlFieldType{}
	scp := uft.GetValueParser("test_text")

	ufv := scp.(UrlFieldValue)

	parsed, err := ufv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
