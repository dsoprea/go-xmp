package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// DcUri is the 'dc' (Dublin Core) namespace URI made a constant to support
	// testing.
	DcUri = "http://purl.org/dc/elements/1.1/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             DcUri,
		PreferredPrefix: "dc",
		Fields: map[string]interface{}{
			"contributor": xmptype.UnorderedTextArrayFieldType{},
			"coverage":    xmptype.TextFieldType{},
			"creator":     xmptype.OrderedTextArrayFieldType{},
			//			"date":        xmptype.DateFieldType{},
			"description": xmptype.LanguageAlternativeArrayFieldType{},
			"format":      xmptype.MimeTypeFieldType{},
			"identifier":  xmptype.TextFieldType{},
			"language":    xmptype.UnorderedTextArrayFieldType{},
			"publisher":   xmptype.UnorderedTextArrayFieldType{},
			"relation":    xmptype.UnorderedTextArrayFieldType{},
			"rights":      xmptype.LanguageAlternativeArrayFieldType{},
			"source":      xmptype.TextFieldType{},
			"subject":     xmptype.UnorderedTextArrayFieldType{},
			"title":       xmptype.LanguageAlternativeArrayFieldType{},
			"type":        xmptype.UnorderedTextArrayFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
