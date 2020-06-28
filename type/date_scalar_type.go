package xmptype

import (
	"time"
)

type DateFieldValue struct {
	raw string
}

var (
	timeLayouts = []string{
		"2006-01-02T15:04:05.999999999Z-07:00",
		"2006-01-02T15:04:05Z-07:00",
		"2006-01-02T15:04Z-07:00",
		"2006-01-02",
		"2006-01",
		"2006",
	}
)

func (dfv DateFieldValue) Parse() (parsed interface{}, err error) {
	for _, layout := range timeLayouts {
		t, err := time.Parse(layout, dfv.raw)
		if err == nil {
			return t, nil
		}
	}

	return nil, ErrValueNotValid
}

// DateFieldType represents a date values.
type DateFieldType struct {
}

// GetValueParser returns an instance of ScalarValueParser initialized to
// parse a specific string.
func (dft DateFieldType) GetValueParser(raw string) ScalarValueParser {
	return DateFieldValue{
		raw: raw,
	}
}
