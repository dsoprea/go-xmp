package xmptype

const (
	// StRefUri is the 'stRef' namespace URI made a constant to support
	// testing.
	StRefUri = "http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
)

var (
	resourceRefFields = map[string]interface{}{
		"alternatePaths":  UriFieldType{},
		"documentID":      UriFieldType{},
		"filePath":        UriFieldType{},
		"fromPart":        PartFieldType{},
		"instanceID":      UriFieldType{},
		"lastModifyDate":  DateFieldType{},
		"manager":         AgentNameFieldType{},
		"managerVariant":  TextFieldType{},
		"manageTo":        UriFieldType{},
		"manageUI":        UriFieldType{},
		"maskMarkers":     ClosedChoiceFieldType{},
		"partMapping":     TextFieldType{},
		"renditionClass":  RenditionClassFieldType{},
		"renditionParams": TextFieldType{},
		"toPart":          PartFieldType{},
		"versionID":       TextFieldType{},
	}

	resourceRefNamespace = Namespace{
		Uri:             StRefUri,
		PreferredPrefix: "stRef",
	}
)

// ResourceRefFieldType encapsulates a set of fields that refer to a
// resource.
type ResourceRefFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (rrft ResourceRefFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := resourceRefFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (rrft ResourceRefFieldType) Namespace() Namespace {
	return resourceRefNamespace
}

func init() {
	registerComplex(ResourceRefFieldType{})
}
