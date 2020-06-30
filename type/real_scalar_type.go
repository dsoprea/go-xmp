package xmptype

import (
	"strconv"
)

// RealFieldValue knows how to parse a float.
type RealFieldValue struct {
	raw string
}

// Parse parses the raw string.
func (rfv RealFieldValue) Parse() (parsed interface{}, err error) {
	f, err := strconv.ParseFloat(rfv.raw, 64)
	if err != nil {
		return nil, ErrValueNotValid
	}

	return f, nil
}

// RealFieldType represents an integer value.
type RealFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (rft RealFieldType) GetValueParser(raw string) ScalarValueParser {
	return RealFieldValue{
		raw: raw,
	}
}
