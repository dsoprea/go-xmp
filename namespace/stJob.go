package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// StJobUri is the 'stJob' namespace URI made a constant to support
	// testing.
	StJobUri = "http://ns.adobe.com/xap/1.0/sType/Job#"
)

func init() {
	namespace := Namespace{
		Uri:             StJobUri,
		PreferredPrefix: "stJob",
		Fields: map[string]interface{}{
			"id":   xmptype.TextFieldType{},
			"name": xmptype.TextFieldType{},
			"url":  xmptype.UrlFieldType{},
		},
	}

	register(namespace)
}
