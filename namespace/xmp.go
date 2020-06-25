package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpUri is the 'xmp' namespace URI made a constant to support testing.
	XmpUri = "http://ns.adobe.com/xap/1.0/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             XmpUri,
		PreferredPrefix: "xmp",
		Fields: map[string]interface{}{
			"BaseURL":      xmptype.UrlFieldType{},
			"CreateDate":   xmptype.DateFieldType{},
			"CreatorTool":  xmptype.AgentNameFieldType{},
			"Identifier":   xmptype.UnorderedTextArrayType{},
			"Label":        xmptype.TextFieldType{},
			"MetadataDate": xmptype.DateFieldType{},
			"ModifyDate":   xmptype.DateFieldType{},
			"Nickname":     xmptype.TextFieldType{},
			"Rating":       xmptype.RealFieldType{},

			// NOTE(dustin): It's unclear how to implemented an "alternative array". Come back to this.
			//			"Thumbnails":
		},
	}

	xmpregistry.Register(namespace)
}
