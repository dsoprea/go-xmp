package xmptype

// MimeTypeFieldValue knows how to parse a mime-type.
type MimeTypeFieldValue struct {
	TextFieldValue
}

// MimeTypeFieldType represents a MIME-type.
type MimeTypeFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (mft MimeTypeFieldType) GetValueParser(raw string) ScalarValueParser {
	return MimeTypeFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
