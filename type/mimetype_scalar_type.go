package xmptype

type MimeTypeFieldValue struct {
	TextFieldValue
}

// MimetypeFieldType represents a MIME-type.
type MimeTypeFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (mft *MimeTypeFieldType) GetValueParser(raw string) ScalarValueParser {
	return &MimeTypeFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
