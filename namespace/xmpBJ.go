package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
)

const (
	// XmpBJUri is the 'xmpBJ' namespace URI made a constant to support
	// testing.
	XmpBJUri = "http://ns.adobe.com/xap/1.0/bj/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpBJUri,
		PreferredPrefix: "xmpBJ",
		// Fields:          map[string]FieldType{
		// 	// NOTE(dustin): Not implementing due to irrelevancy to how we process values.
		// 	// "JobRef":,
		// },
	}

	xmpregistry.Register(namespace)
}
