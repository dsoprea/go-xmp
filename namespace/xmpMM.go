package xmpnamespace

const (
	// XmpMmUri is the 'xmpMM' namespace URI made a constant to support
	// testing.
	XmpMmUri = "http://ns.adobe.com/xap/1.0/mm/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpMmUri,
		PreferredPrefix: "xmpMM",
		Fields: map[string]FieldType{
			"DerivedFrom": ResourceRefFieldType,
			"DocumentID":  GuidFieldType,

			// TODO(dustin): ResourceEventFieldType type is not current implemented. Return to this.
			// 			"History": ResourceEventFieldType,

			"Ingredients":    ResourceRefFieldType,
			"ManagedFrom":    ResourceRefFieldType,
			"Manager":        AgentNameFieldType,
			"ManageTo":       UriFieldType,
			"ManageUI":       UriFieldType,
			"ManagerVariant": TextFieldType,

			"InstanceID":         GuidFieldType,
			"OriginalDocumentID": GuidFieldType,

			// Not implemented due to non-strict nature.
			// "Pantry":,

			"RenditionClass":  RenditionClassFieldType,
			"RenditionParams": TextFieldType,
			"VersionID":       TextFieldType,

			// Not implemented due to irrelevancy because of how we handle values.
			// "Versions":
		},
	}

	register(namespace)
}
