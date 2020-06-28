package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestFrameCountFieldType_GetValueParser(t *testing.T) {
	fft := FrameCountFieldType{}
	scp := fft.GetValueParser("test_text")

	ffv := scp.(FrameCountFieldValue)

	parsed, err := ffv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
