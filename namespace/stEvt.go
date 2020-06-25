package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

// NOTE(dustin): The URI might be misspelled (an extra space) in the standard.

const (
	// StEvtUri is the 'stEvt' namespace URI made a constant to support
	// testing.
	StEvtUri = "http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             StEvtUri,
		PreferredPrefix: "stEvt",
		Fields: map[string]interface{}{
			"action":        xmptype.TextFieldType{},
			"changed":       xmptype.TextFieldType{},
			"instanceID":    xmptype.GuidFieldType{},
			"parameters":    xmptype.TextFieldType{},
			"softwareAgent": xmptype.AgentNameFieldType{},
			"when":          xmptype.DateFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
