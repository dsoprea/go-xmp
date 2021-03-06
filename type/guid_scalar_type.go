package xmptype

// GuidFieldValue knows how to parse a GUID.
type GuidFieldValue struct {
	TextFieldValue
}

// GuidFieldType represents a GUID.
type GuidFieldType struct {
}

// A XMP GUID is an opaque string that may or may not look like a URI according
// to the specification.

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (gft GuidFieldType) GetValueParser(raw string) ScalarValueParser {
	return GuidFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
