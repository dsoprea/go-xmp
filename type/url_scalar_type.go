package xmptype

type UrlFieldValue struct {
	TextFieldValue
}

// UrlFieldType represents an RFC 3986 URI.
type UrlFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (uft UrlFieldType) GetValueParser(raw string) ScalarValueParser {
	return UrlFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
