package xmptype

// PartFieldValue knows how to parse part values.
type PartFieldValue struct {
	TextFieldValue
}

// PartFieldType represents a path specification. It is handled as a string.
type PartFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (pft PartFieldType) GetValueParser(raw string) ScalarValueParser {
	return PartFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
