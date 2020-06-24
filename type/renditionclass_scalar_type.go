package xmptype

type RenditionClassFieldValue struct {
	OpenChoiceFieldValue
}

// OpenChoiceFieldType represents a string that could or should (depending
// on the definition of the field in the standard) be taken from a list of
// defined choices. The defined choices are only suggested candidate values.
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
