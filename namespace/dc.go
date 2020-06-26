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
			"contributor": xmptype.ProperNameFieldType{},
			"coverage":    xmptype.TextFieldType{},
			"creator":     xmptype.ProperNameFieldType{},
			"date":        xmptype.DateFieldType{},
			"description": xmptype.LanguageAlternativeArrayFieldType{},
			"format":      xmptype.MimeTypeFieldType{},
			"identifier":  xmptype.TextFieldType{},
			"language":    xmptype.LocaleFieldType{},
			"publisher":   xmptype.ProperNameFieldType{},
			"relation":    xmptype.TextFieldType{},
			"rights":      xmptype.LanguageAlternativeArrayFieldType{},
			"source":      xmptype.TextFieldType{},
			"subject":     xmptype.TextFieldType{},
			"title":       xmptype.LanguageAlternativeArrayFieldType{},
			"type":        xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
