package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpMmUri is the 'xmpMM' namespace URI made a constant to support
	// testing.
	XmpMmUri = "http://ns.adobe.com/xap/1.0/mm/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpMmUri,
		PreferredPrefix: "xmpMM",
		Fields: map[string]interface{}{
			// "DerivedFrom": xmptype.ResourceRefFieldType{},
			"DocumentID": xmptype.GuidFieldType{},

			// TODO(dustin): ResourceEventFieldType type is not current implemented. Return to this.
			"History": xmptype.OrderedResourceEventArrayType{},

			// "Ingredients":    xmptype.ResourceRefFieldType{},
			// "ManagedFrom":    xmptype.ResourceRefFieldType{},
			"Manager":        xmptype.AgentNameFieldType{},
			"ManageTo":       xmptype.UriFieldType{},
			"ManageUI":       xmptype.UriFieldType{},
			"ManagerVariant": xmptype.TextFieldType{},

			"InstanceID":         xmptype.GuidFieldType{},
			"OriginalDocumentID": xmptype.GuidFieldType{},

			// Not implemented due to non-strict nature.
			// "Pantry":,

			"RenditionClass":  xmptype.RenditionClassFieldType{},
			"RenditionParams": xmptype.TextFieldType{},
			"VersionID":       xmptype.TextFieldType{},

			// Not implemented due to irrelevancy because of how we handle values.
			// "Versions":
		},
	}

	register(namespace)
}
