package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestIntegerFieldType_GetValueParser(t *testing.T) {
	ift := IntegerFieldType{}
	scp := ift.GetValueParser("5")

	ifv := scp.(IntegerFieldValue)

	parsed, err := ifv.Parse()
	log.PanicIf(err)

	if parsed != int64(5) {
		t.Fatalf("Parse is not correct: [%v]", parsed)
	}
}
