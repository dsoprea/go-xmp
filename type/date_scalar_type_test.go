package xmptype

import (
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
)

var (
	testTimezone = time.FixedZone("UTC-5", -5*60*60)
)

func TestDataFieldType_GetValueParser(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	if parsed != time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("Parse is not correct: [%s]", parsed.(time.Time).Format(time.RFC3339Nano))
	}
}

func TestDataFieldType_Parse_Format1(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	if parsed != time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("Parse is not correct: [%s]", parsed.(time.Time).Format(time.RFC3339Nano))
	}
}

func TestDataFieldType_Parse_Format2(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019-05")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	if parsed != time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("Parse is not correct: [%s]", parsed.(time.Time).Format(time.RFC3339Nano))
	}
}

func TestDataFieldType_Parse_Format3(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019-05-07")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	if parsed != time.Date(2019, 5, 7, 0, 0, 0, 0, time.UTC) {
		t.Fatalf("Parse is not correct: [%s]", parsed.(time.Time).Format(time.RFC3339Nano))
	}
}

func TestDataFieldType_Parse_Format4(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019-05-07T12:34Z-05:00")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	actual := parsed.(time.Time)

	expected := time.Date(2019, 5, 7, 12, 34, 0, 0, testTimezone)

	if actual.Equal(expected) != true {
		t.Fatalf("Parse is not correct: [%s] (%d) != [%s] (%d)", actual.Format(time.RFC3339Nano), actual.Nanosecond(), expected.Format(time.RFC3339Nano), expected.Nanosecond())
	}
}

func TestDataFieldType_Parse_Format5(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019-05-07T12:34:56Z-05:00")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	actual := parsed.(time.Time)

	expected := time.Date(2019, 5, 7, 12, 34, 56, 0, testTimezone)

	if actual.Equal(expected) != true {
		t.Fatalf("Parse is not correct: [%s] (%d) != [%s] (%d)", actual.Format(time.RFC3339Nano), actual.Nanosecond(), expected.Format(time.RFC3339Nano), expected.Nanosecond())
	}
}

func TestDataFieldType_Parse_Format6(t *testing.T) {
	dft := DateFieldType{}
	scp := dft.GetValueParser("2019-05-07T12:34:56.000000123Z-05:00")

	dfv := scp.(DateFieldValue)

	parsed, err := dfv.Parse()
	log.PanicIf(err)

	actual := parsed.(time.Time)

	expected := time.Date(2019, 5, 7, 12, 34, 56, 123, testTimezone)

	if actual.Equal(expected) != true {
		t.Fatalf("Parse is not correct: [%s] (%d) != [%s] (%d)", actual.Format(time.RFC3339Nano), actual.Nanosecond(), expected.Format(time.RFC3339Nano), expected.Nanosecond())
	}
}
