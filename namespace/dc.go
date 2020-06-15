package xmpnamespace

const (
	// DcUri is the 'dc' URI made a constant to support testing.
	DcUri = "http://purl.org/dc/elements/1.1/"
)

func init() {
	namespace := Namespace{
		Uri:             DcUri,
		PreferredPrefix: "dc",
		Fields: map[string]FieldType{
			"contributor": ProperNameFieldType,
			"coverage":    TextFieldType,
			"creator":     ProperNameFieldType,
			"date":        DateFieldType,
			"description": LanguageAlternativeFieldType,
			"format":      MimeTypeFieldType,
			"identifier":  TextFieldType,
			"language":    LocaleFieldType,
			"publisher":   ProperNameFieldType,
			"relation":    TextFieldType,
			"rights":      LanguageAlternativeFieldType,
			"source":      TextFieldType,
			"subject":     TextFieldType,
			"title":       LanguageAlternativeFieldType,
			"type":        TextFieldType,
		},
	}

	register(namespace)
}
