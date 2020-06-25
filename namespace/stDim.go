package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// StDimUri is the 'stDim' namespace URI made a constant to support
	// testing.
	StDimUri = "http://ns.adobe.com/xap/1.0/sType/Dimensions#"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             StDimUri,
		PreferredPrefix: "stDim",
		Fields: map[string]interface{}{
			"h":    xmptype.RealFieldType{},
			"w":    xmptype.RealFieldType{},
			"unit": xmptype.OpenChoiceFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
