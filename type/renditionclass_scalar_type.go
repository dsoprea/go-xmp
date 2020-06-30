package xmptype

// RenditionClassFieldValue knows how to parse a rendition-class string.
type RenditionClassFieldValue struct {
	OpenChoiceFieldValue
}

// RenditionClassFieldType describes a rendition-class value.
type RenditionClassFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (rcft RenditionClassFieldType) GetValueParser(raw string) ScalarValueParser {
	return RenditionClassFieldValue{
		OpenChoiceFieldValue: OpenChoiceFieldValue{
			raw: raw,
			choices: []string{
				"default",
				"draft",
				"low-res",
				"proof",
				"screen",
				"thumbnail",
			},
		},
	}
}
