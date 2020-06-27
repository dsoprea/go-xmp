package xmptype

import (
	"errors"
)

var (
	// ErrChoicesNotOverridden indicates that a particular type is not correctly
	// overridden.
	ErrChoicesNotOverridden = errors.New("choices type method must be overridden")
)

// OpenChoiceFieldValue encapsulates the choices and the value found in the
// document. Passes even if value does not appear among choices.
type OpenChoiceFieldValue struct {
	raw     string
	choices []string
}

// NewOpenChoiceFieldValue returns a new struct.
func NewOpenChoiceFieldValue(raw string, choices []string) OpenChoiceFieldValue {
	return OpenChoiceFieldValue{
		raw:     raw,
		choices: choices,
	}
}

// Raw returns the original text to be parsed.
func (ocfv OpenChoiceFieldValue) Raw() string {
	return ocfv.raw
}

// Parse returns the string value.
func (ocfv OpenChoiceFieldValue) Parse() (parsed interface{}, err error) {
	return ocfv.Raw(), nil
}

// OpenChoiceFieldType represents a string that could or should (depending
// on the definition of the field in the standard) be taken from a list of
// defined choices. The defined choices are only suggested candidate values.
type OpenChoiceFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to parse
// a specific string. It must be overridden for each XMP type that needs it.
func (ocft OpenChoiceFieldType) GetValueParser(raw string) ScalarValueParser {

	panic(ErrChoicesNotOverridden)

}

// ClosedChoiceFieldValue encapsulates the choices and the value found in the
// document. Fails if the value does not appear among the choices.
type ClosedChoiceFieldValue struct {
	raw     string
	choices []string
}

// NewClosedChoiceFieldValue returns a new struct.
func NewClosedChoiceFieldValue(raw string, choices []string) ClosedChoiceFieldValue {
	return ClosedChoiceFieldValue{
		raw:     raw,
		choices: choices,
	}
}

// Raw returns the original text to be parsed.
func (ccfv ClosedChoiceFieldValue) Raw() string {
	return ccfv.raw
}

// Parse returns the string value if it matches one of the available choices, or
// returns ErrValueNotValid.
func (ccfv ClosedChoiceFieldValue) Parse() (parsed interface{}, err error) {
	for _, choice := range ccfv.choices {
		if choice == ccfv.raw {
			return ccfv.raw, nil
		}
	}

	return nil, ErrValueNotValid
}

// ClosedChoiceFieldType represents a string that could or should (depending on
// the definition of the field in the standard) be taken from a list of defined
// choices. The value must be one of the defined choices.
type ClosedChoiceFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to parse
// a specific string. It must be overridden for each XMP type that needs it.
func (ccft ClosedChoiceFieldType) GetValueParser(raw string) ScalarValueParser {

	panic(ErrChoicesNotOverridden)

}
