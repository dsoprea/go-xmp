package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// DcUri is the 'dc' (Dublin Core) namespace URI made a constant to support
	// testing.
	DcUri = "http://purl.org/dc/elements/1.1/"
)

func init() {
	namespace := Namespace{
		Uri:             DcUri,
		PreferredPrefix: "dc",
		Fields: map[string]interface{}{
			"contributor": xmptype.ProperNameFieldType{},
			"coverage":    xmptype.TextFieldType{},
			"creator":     xmptype.ProperNameFieldType{},
			"date":        xmptype.DateFieldType{},

			// TODO(dustin): Let's revisit once we can write a unit-test for it.
			// "description": xmptype.LanguageAlternativeFieldType{},

			"format":     xmptype.MimeTypeFieldType{},
			"identifier": xmptype.TextFieldType{},
			"language":   xmptype.LocaleFieldType{},
			"publisher":  xmptype.ProperNameFieldType{},
			"relation":   xmptype.TextFieldType{},

			// TODO(dustin): Let's revisit once we can write a unit-test for it.
			// "rights":      xmptype.LanguageAlternativeFieldType{},

			"source":  xmptype.TextFieldType{},
			"subject": xmptype.TextFieldType{},

			// TODO(dustin): Let's revisit once we can write a unit-test for it.
			// "title":       xmptype.LanguageAlternativeFieldType{},

			"type": xmptype.TextFieldType{},
		},
	}

	register(namespace)
}
