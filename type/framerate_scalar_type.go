package xmptype

// FrameRateFieldValue knows how to parse frame-count expressions.
type FrameRateFieldValue struct {
	TextFieldValue
}

// FrameRateFieldType represents a frame-count specification. It is handled
// as a string.
type FrameRateFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (fcft FrameRateFieldType) GetValueParser(raw string) ScalarValueParser {
	return FrameRateFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
