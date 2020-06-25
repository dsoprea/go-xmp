package xmptype

import (
	"errors"
)

var (
	// ErrValueNotValid indicates that the value was not valid for the type
	// prescribed in the specification.
	ErrValueNotValid = errors.New("value not valid/allowed")
)

// ScalarFieldType represents a factory for ScalarValueParser types.
type ScalarFieldType interface {
	// GetValueParser returns an instance of ScalarValueParser initialized to
	// parse a specific string.
	GetValueParser(raw string) ScalarValueParser
}

// ScalarValueParser knows how to parse a value encoded to a string.
type ScalarValueParser interface {
	// Parse parses the wrapped string to a specific type.
	Parse() (interface{}, error)
}
