package xmptype

type BooleanFieldValue struct {
	raw string
}

func (bfv BooleanFieldValue) Parse() (parsed interface{}, err error) {
	if bfv.raw == "True" {
		return true, nil
	} else if bfv.raw == "False" {
		return false, nil
	}

	return nil, ErrValueNotValid
}

// BooleanFieldType represents a boolean value.
type BooleanFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (bft *BooleanFieldType) GetValueParser(raw string) ScalarValueParser {
	return &BooleanFieldValue{
		raw: raw,
	}
}
