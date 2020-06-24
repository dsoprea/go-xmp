package xmptype

const (
	// StJobUri is the 'stJob' namespace URI made a constant to support
	// testing.
	StJobUri = "http://ns.adobe.com/xap/1.0/sType/Job#"
)

var (
	jobFields = map[string]interface{}{
		"id":   TextFieldType{},
		"name": TextFieldType{},
		"url":  UrlFieldType{},
	}

	jobNamespace = Namespace{
		Uri:             StJobUri,
		PreferredPrefix: "stJob",
	}
)

// JobFieldType encapsulates a set of fields that refer to a job.
type JobFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (jft JobFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := jobFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (jft JobFieldType) Namespace() Namespace {
	return jobNamespace
}

func init() {
	registerComplex(JobFieldType{})
}
