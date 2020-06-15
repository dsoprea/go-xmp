package xmpnamespace

const (
	// XmpidqUri is the 'xmpidq' namespace URI made a constant to support
	// testing.
	XmpidqUri = "http://ns.adobe.com/xmp/Identifier/qual/1.0/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpidqUri,
		PreferredPrefix: "xmpidq",
		Fields: map[string]FieldType{
			"Scheme": TextFieldType,
		},
	}

	register(namespace)
}
