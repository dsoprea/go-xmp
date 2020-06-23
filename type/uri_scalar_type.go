package xmptype

type UriFieldValue struct {
	TextFieldValue
}

// UriFieldType represents an RFC 3986 URI.
type UriFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (uft *UriFieldType) GetValueParser(raw string) ScalarValueParser {
	tv := TextFieldValue{
		raw: raw,
	}

	return &UriFieldValue{
		TextFieldValue: tv,
	}
}
