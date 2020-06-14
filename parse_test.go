package xmp

import (
	"bytes"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestNewParser(t *testing.T) {
	// Not much else we can test at this juncture.
	NewParser(nil)
}

func TestParser_Parse(t *testing.T) {
	data := GetTestData()
	b := bytes.NewBuffer(data)
	xp := NewParser(b)

	err := xp.Parse()
	log.PanicIf(err)

	// TODO(dustin): !! Add validation.
}
