package xmptype

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestAgentNameFieldType_GetValueParser(t *testing.T) {
	anft := AgentNameFieldType{}
	scp := anft.GetValueParser("test_text")

	anfv := scp.(AgentNameFieldValue)

	parsed, err := anfv.Parse()
	log.PanicIf(err)

	if parsed != "test_text" {
		t.Fatalf("Parse is not correct: [%s]", parsed)
	}
}
