package xmptype

import (
	"errors"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

var (
	typeLogger = log.NewLogger("xmp.type")
)

var (
	// // ErrChildFieldNotValid indicates that the given complex type has no child
	// // field with the given name.
	// ErrChildFieldNotValid = errors.New("child field not valid")

	// ErrValueNotValid indicates that the value was not valid for the type
	// prescribed in the specification.
	ErrValueNotValid = errors.New("value not valid/allowed")

// // ErrComplexTypeNotFound indicates that the complex-type is not known.
// ErrComplexTypeNotFound = errors.New("complex-type not known")
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

// Array types

type ArrayValue interface {
	FullName() xmpregistry.XmpPropertyName
	Count() int
}

type ArrayType interface {
	New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue
}

type BaseArrayValue struct {
	fullName  xmpregistry.XmpPropertyName
	collected []interface{}
}

func (bav BaseArrayValue) FullName() xmpregistry.XmpPropertyName {
	return bav.fullName
}

func (bav BaseArrayValue) Count() int {
	return len(bav.collected)
}

// Ordered array semantics

// TODO(dustin): Ordered array yet-to-implement: CuePointParam, Marker, ResourceEvent, Version, Colorant, Marker, Layer, "point" (?)

type OrderedArrayValue struct {
	BaseArrayValue
}

type OrderedArrayType struct {
}

func (oat OrderedArrayType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	return OrderedArrayValue{
		BaseArrayValue: BaseArrayValue{
			fullName:  fullName,
			collected: collected,
		},
	}
}

type OrderedTextArrayType struct {
	OrderedArrayType
}

type OrderedUriArrayType struct {
	OrderedArrayType
}

type OrderedResourceEventArrayType struct {
	OrderedArrayType
}

// Unordered array semantics

// TODO(dustin): Unordered array yet-to-implement: XPath, ResourceRef, "struct" (?), Job, Font, Media, Track, Ancestor

type UnorderedArrayValue struct {
	BaseArrayValue
}

type UnorderedArrayType struct {
}

func (uat UnorderedArrayType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	return UnorderedArrayValue{
		BaseArrayValue: BaseArrayValue{
			fullName:  fullName,
			collected: collected,
		},
	}
}

type UnorderedTextArrayType struct {
	UnorderedArrayType
}

// Alternatives array semantics

type AlternativeArrayType struct {
	BaseArrayValue
}

func (aat AlternativeArrayType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	return AlternativeArrayType{
		BaseArrayValue: BaseArrayValue{
			fullName:  fullName,
			collected: collected,
		},
	}
}
