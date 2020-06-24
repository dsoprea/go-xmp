package xmptype

const (
	// StFntUri is the 'stFnt' namespace URI made a constant to support
	// testing.
	StFntUri = "http:ns.adobe.com/xap/1.0/sType/Font#"
)

// Specification has both "String" and "Text" fields. "String" fields are not
// defined, so using "Text" instead.

var (
	fontFields = map[string]interface{}{
		"childFontFiles": TextFieldType{},
		"composite":      BooleanFieldType{},
		"fontFace":       TextFieldType{},
		"fontFamily":     TextFieldType{},
		"fontFileName":   TextFieldType{},
		"fontName":       TextFieldType{},
		"fontType":       TextFieldType{},
		"versionString":  TextFieldType{},
	}

	fontNamespace = Namespace{
		Uri:             StFntUri,
		PreferredPrefix: "stFnt",
	}
)

// FontFieldType encapsulates a set of fields that refer to a
// font.
type FontFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (fft FontFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := fontFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (fft FontFieldType) Namespace() Namespace {
	return fontNamespace
}

func init() {
	registerComplex(FontFieldType{})
}
