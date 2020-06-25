package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpGImageUri is the 'xmpGImg' namespace URI made a constant to support
	// testing.
	XmpGImageUri = "http://ns.adobe.com/xap/1.0/g/img/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpGImageUri,
		PreferredPrefix: "xmpGImg",
		Fields: map[string]interface{}{
			"format": xmptype.TextFieldType{},
			"height": xmptype.IntegerFieldType{},
			"width":  xmptype.IntegerFieldType{},
			"image":  xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
