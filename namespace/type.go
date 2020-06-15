package xmpnamespace

// FieldType describes the type of value that a field is expected to have. (0)
// is not a valid value.
type FieldType int

const (
	// BooleanFieldType represents a boolean value.
	BooleanFieldType FieldType = iota + 1

	// DateFieldType represents a date values.
	DateFieldType FieldType = iota + 1

	// IntegerFieldType represents an integer value.
	IntegerFieldType FieldType = iota + 1

	// RealFieldType represents a float value.
	RealFieldType FieldType = iota + 1

	// TextFieldType represents a text value.
	TextFieldType FieldType = iota + 1

	// AgentNameFieldType represents the name of an XMP processor.
	AgentNameFieldType FieldType = iota + 1

	// ChoiceFieldType represents a string tht could or should (depending on
	// the definition of the field in the standard) be taken from a list of
	// defined choices.
	ChoiceFieldType FieldType = iota + 1

	// GuidFieldType represents a GUID.
	GuidFieldType FieldType = iota + 1

	// LanguageAlternativeFieldType represents a value from a list of defined
	// choices, chosen based on the desired language.
	LanguageAlternativeFieldType FieldType = iota + 1

	// LocaleFieldType represents a RFC 3066 language-code.
	LocaleFieldType FieldType = iota + 1

	// MimeTypeFieldType represents a MIME type.
	MimeTypeFieldType FieldType = iota + 1

	// ProperNameFieldType represents the name of a person or organization.
	ProperNameFieldType FieldType = iota + 1

	// RenditionClassFieldType denoting how the resource wil be used/presented.
	RenditionClassFieldType FieldType = iota + 1

	// ResourceRefFieldType encapsulates a set of fields that refer to a
	// resource (this has limited utility in go-xmp given how it manages
	// values).
	ResourceRefFieldType FieldType = iota + 1

	// UriFieldType represents an RFC 3986 URI.
	UriFieldType FieldType = iota + 1

	// UrlFieldType represents a URL as defined in
	// http://www.w3.org/TR/2001/NOTE-uri-clarification-20010921 .
	UrlFieldType FieldType = iota + 1
)
