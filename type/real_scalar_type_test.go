package xmptype

import (
	"math"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestRealFieldValue_Parse(t *testing.T) {
	rfv := RealFieldValue{
		raw: "123.456",
	}

	parsed, err := rfv.Parse()
	log.PanicIf(err)

	f := parsed.(float64)

	if f < 123.456 || f >= math.Nextafter(f, f+1) {
		t.Fatalf("Parse is not correct: [%6.4f]", f)
	}
}

func TestRealFieldType_GetValueParser(t *testing.T) {
	rft := RealFieldType{}
	scp := rft.GetValueParser("123.456")

	rfv := scp.(RealFieldValue)

	parsed, err := rfv.Parse()
	log.PanicIf(err)

	f := parsed.(float64)

	if f < 123.456 || f >= math.Nextafter(f, f+1) {
		t.Fatalf("Parse is not correct: [%6.4f]", f)
	}
}
