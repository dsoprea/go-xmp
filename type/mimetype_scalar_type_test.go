package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestMimeTypeFieldType_GetValueParser(t *testing.T) {
	mft := MimeTypeFieldType{}
	scp := mft.GetValueParser("test_text")

	mfv := scp.(MimeTypeFieldValue)

	parsed, err := mfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
