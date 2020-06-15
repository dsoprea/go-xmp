package xmpnamespace

const (
	// XmpUri is the 'xmp' namespace URI made a constant to support testing.
	XmpUri = "http://ns.adobe.com/xap/1.0/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpUri,
		PreferredPrefix: "xmp",
		Fields: map[string]FieldType{
			"BaseURL":      UrlFieldType,
			"CreateDate":   DateFieldType,
			"CreatorTool":  AgentNameFieldType,
			"Identifier":   TextFieldType,
			"Label":        TextFieldType,
			"MetadataDate": DateFieldType,
			"ModifyDate":   DateFieldType,
			"Nickname":     TextFieldType,
			"Rating":       RealFieldType,

			// NOTE(dustin): It's unclear how to implemented an "alternative array". Come back to this.
			//			"Thumbnails":
		},
	}

	register(namespace)
}
