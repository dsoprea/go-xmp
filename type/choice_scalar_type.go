package xmptype

import (
	"errors"
)

var (
	ErrChoicesNotSet = errors.New("choices not overridden by type")
)

type OpenChoiceFieldValue struct {
	choices []string
	raw     string
}

func (ocfv OpenChoiceFieldValue) Parse() (parsed interface{}, err error) {
	return ocfv.raw, nil
}

type ClosedChoiceFieldValue struct {
	choices []string
	raw     string
}

func (ccfv ClosedChoiceFieldValue) Parse() (parsed interface{}, err error) {
	for _, choice := range ccfv.choices {
		if choice == ccfv.raw {
			return ccfv.raw, nil
		}
	}

	return nil, ErrValueNotValid
}

// OpenChoiceFieldType represents a string that could or should (depending
// on the definition of the field in the standard) be taken from a list of
// defined choices. The defined choices are only suggested candidate values.
type OpenChoiceFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (ocft OpenChoiceFieldType) GetValueParser(raw string) ScalarValueParser {
	return OpenChoiceFieldValue{
		raw: raw,
	}
}

// Choices returns the set of possible choices. This should be overridden by a
// purpose-specific type.
func (ocft OpenChoiceFieldType) Choices() []string {

	panic(ErrChoicesNotSet)

	return nil
}

// ClosedChoiceFieldType represents a string that could or should (depending on
// the definition of the field in the standard) be taken from a list of defined
// choices. The value must be one of the defined choices.
type ClosedChoiceFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (ccft ClosedChoiceFieldType) GetValueParser(raw string) ScalarValueParser {
	return ClosedChoiceFieldValue{
		raw: raw,
	}
}

// Choices returns the set of required choices. This should be overridden by a
// purpose-specific type.
func (ccft ClosedChoiceFieldType) Choices() []string {

	panic(ErrChoicesNotSet)

	return nil
}
