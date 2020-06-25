package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// StFntUri is the 'stFnt' namespace URI made a constant to support
	// testing.
	StFntUri = "http:ns.adobe.com/xap/1.0/sType/Font#"
)

// Specification has both "String" and "Text" fields. "String" fields are not
// defined, so using "Text" instead.

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             StFntUri,
		PreferredPrefix: "stFnt",
		Fields: map[string]interface{}{
			"childFontFiles": xmptype.OrderedTextArrayType{},
			"composite":      xmptype.BooleanFieldType{},
			"fontFace":       xmptype.TextFieldType{},
			"fontFamily":     xmptype.TextFieldType{},
			"fontFileName":   xmptype.TextFieldType{},
			"fontName":       xmptype.TextFieldType{},
			"fontType":       xmptype.TextFieldType{},
			"versionString":  xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
