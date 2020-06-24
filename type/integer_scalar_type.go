package xmptype

import (
	"strconv"
)

type IntegerFieldValue struct {
	raw string
}

func (ifv IntegerFieldValue) Parse() (parsed interface{}, err error) {
	n, err := strconv.ParseInt(ifv.raw, 10, 64)
	if err != nil {
		return nil, ErrValueNotValid
	}

	return n, nil
}

// IntegerFieldType represents an integer value.
type IntegerFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (ift IntegerFieldType) GetValueParser(raw string) ScalarValueParser {
	return IntegerFieldValue{
		raw: raw,
	}
}
