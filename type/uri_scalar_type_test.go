package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestUriFieldType_GetValueParser(t *testing.T) {
	uft := UriFieldType{}
	scp := uft.GetValueParser("test_text")

	ufv := scp.(UriFieldValue)

	parsed, err := ufv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
