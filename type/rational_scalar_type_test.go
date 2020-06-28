package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestRationalFieldValue_Parse(t *testing.T) {
	rfv := RationalFieldValue{
		raw: "11/22",
	}

	parsed, err := rfv.Parse()
	log.PanicIf(err)

	r := Rational{
		Numerator:   int64(11),
		Denominator: int64(22),
	}

	if parsed != r {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}

func TestRationalFieldType_GetValueParser(t *testing.T) {
	rft := RationalFieldType{}
	scp := rft.GetValueParser("11/22")

	rfv := scp.(RationalFieldValue)

	parsed, err := rfv.Parse()
	log.PanicIf(err)

	r := Rational{
		Numerator:   int64(11),
		Denominator: int64(22),
	}

	if parsed != r {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
