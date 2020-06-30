package xmptype

// FrameCountFieldValue knows how to parse frame-count expressions.
type FrameCountFieldValue struct {
	TextFieldValue
}

// FrameCountFieldType represents a frame-count specification. It is handled
// as a string.
type FrameCountFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (fcft FrameCountFieldType) GetValueParser(raw string) ScalarValueParser {
	return FrameCountFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
