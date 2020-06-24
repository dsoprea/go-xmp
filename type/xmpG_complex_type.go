package xmptype

const (
	// XmpGUri is the 'xmpG' namespace URI made a constant to support
	// testing.
	XmpGUri = "http://ns.adobe.com/xap/1.0/g/"
)

var (
	colorantFields = map[string]interface{}{
		"A":          IntegerFieldType{},
		"B":          IntegerFieldType{},
		"L":          RealFieldType{},
		"black":      RealFieldType{},
		"cyan":       RealFieldType{},
		"magenta":    RealFieldType{},
		"yellow":     RealFieldType{},
		"blue":       IntegerFieldType{},
		"green":      IntegerFieldType{},
		"red":        IntegerFieldType{},
		"mode":       ClosedChoiceFieldType{},
		"swatchName": TextFieldType{},
		"type":       ClosedChoiceFieldType{},
	}

	colorantNamespace = Namespace{
		Uri:             XmpGUri,
		PreferredPrefix: "xmpG",
	}
)

// ColorantFieldType encapsulates a set of fields that refer to a
// swatch.
type ColorantFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (cft ColorantFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := colorantFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (cft ColorantFieldType) Namespace() Namespace {
	return colorantNamespace
}

func init() {
	registerComplex(ColorantFieldType{})
}
