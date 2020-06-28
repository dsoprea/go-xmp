package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestFrameRateFieldType_GetValueParser(t *testing.T) {
	fft := FrameRateFieldType{}
	scp := fft.GetValueParser("test_text")

	ffv := scp.(FrameRateFieldValue)

	parsed, err := ffv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
