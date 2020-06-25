package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmlUri is the 'xml' namespace URI made a constant to support testing.
	XmlUri = "http://www.w3.org/XML/1998/namespace"
)

// We only define this type so that we parse xml:lang attributes.

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmlUri,
		PreferredPrefix: "xml",
		Fields: map[string]interface{}{
			"lang": xmptype.TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
