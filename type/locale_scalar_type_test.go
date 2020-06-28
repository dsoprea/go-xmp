package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestLocaleFieldType_GetValueParser(t *testing.T) {
	lft := LocaleFieldType{}
	scp := lft.GetValueParser("test_text")

	lfv := scp.(LocaleFieldValue)

	parsed, err := lfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
