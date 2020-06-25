package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// StVerUri is the 'stVer' namespace URI made a constant to support
	// testing.
	StVerUri = "http://ns.adobe.com/xap/1.0/sType/Version#"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             StVerUri,
		PreferredPrefix: "stVer",
		Fields: map[string]interface{}{
			"comments": xmptype.TextFieldType{},
			// "event":    xmptype.ResourceEventFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
