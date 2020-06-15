package xmpnamespace

const (
	// XmpRightsUri is the 'xmpRights' namespace URI made a constant to support
	// testing.
	XmpRightsUri = "http://ns.adobe.com/xap/1.0/rights/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpRightsUri,
		PreferredPrefix: "xmpRights",
		Fields: map[string]FieldType{
			"Certificate":  TextFieldType,
			"Marked":       BooleanFieldType,
			"Owner":        ProperNameFieldType,
			"UsageTerms":   LanguageAlternativeFieldType,
			"WebStatement": TextFieldType,
		},
	}

	register(namespace)
}
