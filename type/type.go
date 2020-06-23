package xmptype

import (
	"errors"
)

var (
	// ErrChildFieldNotValid indicates that the given complex type has no child
	// field with the given name.
	ErrChildFieldNotValid = errors.New("child field not valid")

	// ErrValueNotValid indicates that the value was not valid for the type
	// prescribed in the specification.
	ErrValueNotValid = errors.New("value not valid/allowed")
)

// ScalarFieldType represents a factory for ScalarValueParser types.
type ScalarFieldType interface {
	// GetValueParser returns an instance of ScalarValueParser initialized to
	// parse a specific string.
	GetValueParser(raw string) ScalarValueParser
}

// ScalarValueParser knows how to parse a value encoded to a string.
type ScalarValueParser interface {
	// Parse parses the wrapped string to a specific type.
	Parse() (interface{}, error)
}

// Namespace describes an XML namespace.
type Namespace struct {
	// Uri is the URI of a namespace (it should be regarded as a string only;
	// XML namespaces are not necssarily valid Internet resources).
	Uri string

	// PreferredPrefix is the preferred naming-prefix prescribed by the
	// governing standard of this namespace.
	PreferredPrefix string
}

// ComplexFieldType represents a complex value (comprised of child nodes).
type ComplexFieldType interface {
	// ChildFieldType returns the field-type for the immediate child with the
	// given name.
	ChildFieldType(fieldName string) (ft interface{}, err error)

	// Namespace returns the namespace info the node/children of this type.
	Namespace() Namespace
}
