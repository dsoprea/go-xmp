package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// StRefUri is the 'stRef' namespace URI made a constant to support
	// testing.
	StRefUri = "http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
)

func init() {
	namespace := Namespace{
		Uri:             StRefUri,
		PreferredPrefix: "stRef",
		Fields: map[string]interface{}{
			"alternatePaths":  xmptype.OrderedUriArrayType{},
			"documentID":      xmptype.UriFieldType{},
			"filePath":        xmptype.UriFieldType{},
			"fromPart":        xmptype.PartFieldType{},
			"instanceID":      xmptype.UriFieldType{},
			"lastModifyDate":  xmptype.DateFieldType{},
			"manager":         xmptype.AgentNameFieldType{},
			"managerVariant":  xmptype.TextFieldType{},
			"manageTo":        xmptype.UriFieldType{},
			"manageUI":        xmptype.UriFieldType{},
			"maskMarkers":     xmptype.ClosedChoiceFieldType{},
			"partMapping":     xmptype.TextFieldType{},
			"renditionClass":  xmptype.RenditionClassFieldType{},
			"renditionParams": xmptype.TextFieldType{},
			"toPart":          xmptype.PartFieldType{},
			"versionID":       xmptype.TextFieldType{},
		},
	}

	register(namespace)
}
