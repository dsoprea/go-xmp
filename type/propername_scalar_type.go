package xmptype

type ProperNameFieldValue struct {
	TextFieldValue
}

// ProperNameFieldType represents the name of a person or thing.
type ProperNameFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (pnft *ProperNameFieldType) GetValueParser(raw string) ScalarValueParser {
	return &ProperNameFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
