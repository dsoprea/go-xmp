package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpRightsUri is the 'xmpRights' namespace URI made a constant to support
	// testing.
	XmpRightsUri = "http://ns.adobe.com/xap/1.0/rights/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpRightsUri,
		PreferredPrefix: "xmpRights",
		Fields: map[string]interface{}{
			"Certificate": xmptype.TextFieldType{},
			"Marked":      xmptype.BooleanFieldType{},
			"Owner":       xmptype.ProperNameFieldType{},

			// TODO(dustin): Let's revisit once we can write a unit-test for it.
			// "UsageTerms":   xmptype.LanguageAlternativeFieldType{},

			"WebStatement": xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
