package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpidqUri is the 'xmpidq' namespace URI made a constant to support
	// testing.
	XmpidqUri = "http://ns.adobe.com/xmp/Identifier/qual/1.0/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpidqUri,
		PreferredPrefix: "xmpidq",
		Fields: map[string]interface{}{
			"Scheme": xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
