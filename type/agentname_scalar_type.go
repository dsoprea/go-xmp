package xmptype

// AgentNameFieldValue knows how to parse an agent-name value.
type AgentNameFieldValue struct {
	TextFieldValue
}

// AgentNameFieldType represents the name of an XMP processor.
type AgentNameFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (anft AgentNameFieldType) GetValueParser(raw string) ScalarValueParser {
	return AgentNameFieldValue{
		TextFieldValue: TextFieldValue{
			raw: raw,
		},
	}
}
