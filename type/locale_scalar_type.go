package xmptype

type LocaleFieldValue struct {
	TextFieldValue
}

// LocaleFieldType represents a locale.
type LocaleFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (lft LocaleFieldType) GetValueParser(raw string) ScalarValueParser {
	return LocaleFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
