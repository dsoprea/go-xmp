package xmptype

const (
	// XmpGImageUri is the 'xmpGImg' namespace URI made a constant to support
	// testing.
	XmpGImageUri = "http://ns.adobe.com/xap/1.0/g/img/"
)

var (
	thumbnailFields = map[string]interface{}{
		"format": TextFieldType{},
		"height": IntegerFieldType{},
		"width":  IntegerFieldType{},
		"image":  TextFieldType{},
	}

	thumbnailNamespace = Namespace{
		Uri:             XmpGImageUri,
		PreferredPrefix: "xmpGImg",
	}
)

// ThumbnailFieldType encapsulates a set of fields that refer to a
// thumbnail.
type ThumbnailFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (tft ThumbnailFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := thumbnailFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (tft ThumbnailFieldType) Namespace() Namespace {
	return thumbnailNamespace
}

func init() {
	registerComplex(ThumbnailFieldType{})
}
