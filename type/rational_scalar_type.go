package xmptype

import (
	"strconv"
	"strings"
)

// RationalFieldValue knows how to parse rationals/fractions.
type RationalFieldValue struct {
	raw string
}

// Rational describes a single rational/fraction.
type Rational struct {
	Numerator   int64
	Denominator int64
}

// Parse parses the raw string value.
func (rfv RationalFieldValue) Parse() (parsed interface{}, err error) {
	parts := strings.Split(rfv.raw, "/")
	if len(parts) != 2 {
		return nil, ErrValueNotValid
	}

	numerator, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, ErrValueNotValid
	}

	denominator, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, ErrValueNotValid
	}

	rational := Rational{
		Numerator:   numerator,
		Denominator: denominator,
	}

	return rational, nil
}

// RationalFieldType represents an integer value.
type RationalFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (rft RationalFieldType) GetValueParser(raw string) ScalarValueParser {
	return RationalFieldValue{
		raw: raw,
	}
}
