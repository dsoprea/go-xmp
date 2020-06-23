package xmptype

const (
	// StVerUri is the 'stVer' namespace URI made a constant to support
	// testing.
	StVerUri = "http://ns.adobe.com/xap/1.0/sType/Version#"
)

var (
	versionFields = map[string]interface{}{
		"comments": TextFieldType{},
		"event":    ResourceEventFieldType{},
	}

	versionNamespace = Namespace{
		Uri:             StVerUri,
		PreferredPrefix: "stVer",
	}
)

// VersionFieldType encapsulates a set of fields that refer to a
// version.
type VersionFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (vft VersionFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := versionFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (vft VersionFieldType) Namespace() Namespace {
	return versionNamespace
}
