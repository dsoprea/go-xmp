package xmptype

// NOTE(dustin): The URI might be misspelled (an extra space) in the standard.

const (
	// StEvtUri is the 'stEvt' namespace URI made a constant to support
	// testing.
	StEvtUri = "http://ns.adobe.com/xap/1.0/sType/ResourceEvent#"
)

var (
	resourceEventFields = map[string]interface{}{
		"action":        TextFieldType{},
		"changed":       TextFieldType{},
		"instanceID":    GuidFieldType{},
		"parameters":    TextFieldType{},
		"softwareAgent": AgentNameFieldType{},
		"when":          DateFieldType{},
	}

	resourceEventNamespace = Namespace{
		Uri:             StEvtUri,
		PreferredPrefix: "stEvt",
	}
)

// ResourceEventFieldType encapsulates a set of fields that refer to a
// resource-event.
type ResourceEventFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (reft ResourceEventFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := resourceEventFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (reft ResourceEventFieldType) Namespace() Namespace {
	return resourceEventNamespace
}
