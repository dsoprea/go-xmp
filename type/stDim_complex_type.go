package xmptype

const (
	// StDimUri is the 'stDim' namespace URI made a constant to support
	// testing.
	StDimUri = "http://ns.adobe.com/xap/1.0/sType/Dimensions#"
)

var (
	dimensionsFields = map[string]interface{}{
		"h":    RealFieldType{},
		"w":    RealFieldType{},
		"unit": OpenChoiceFieldType{},
	}

	dimensionsNamespace = Namespace{
		Uri:             StDimUri,
		PreferredPrefix: "stDim",
	}
)

// DimensionsFieldType encapsulates a set of fields that refer to the
// dimensions.
type DimensionsFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (dft DimensionsFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := dimensionsFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (dft DimensionsFieldType) Namespace() Namespace {
	return dimensionsNamespace
}
