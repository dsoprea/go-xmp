package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpGUri is the 'xmpG' namespace URI made a constant to support
	// testing.
	XmpGUri = "http://ns.adobe.com/xap/1.0/g/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpGUri,
		PreferredPrefix: "xmpG",
		Fields: map[string]interface{}{
			"A":          xmptype.IntegerFieldType{},
			"B":          xmptype.IntegerFieldType{},
			"L":          xmptype.RealFieldType{},
			"black":      xmptype.RealFieldType{},
			"cyan":       xmptype.RealFieldType{},
			"magenta":    xmptype.RealFieldType{},
			"yellow":     xmptype.RealFieldType{},
			"blue":       xmptype.IntegerFieldType{},
			"green":      xmptype.IntegerFieldType{},
			"red":        xmptype.IntegerFieldType{},
			"mode":       xmptype.ClosedChoiceFieldType{},
			"swatchName": xmptype.TextFieldType{},
			"type":       xmptype.ClosedChoiceFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
