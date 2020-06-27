package xmptype

// TextFieldValue parses the encapsulated value to a string.
type TextFieldValue struct {
	raw string
}

// Parse parses the raw value to a string value on-the-fly.
func (tfv TextFieldValue) Parse() (parsed interface{}, err error) {
	return tfv.raw, nil
}

// TextFieldType represents a string value.
type TextFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (tft TextFieldType) GetValueParser(raw string) ScalarValueParser {
	return TextFieldValue{
		raw: raw,
	}
}
