package xmp

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/type"
)

var (
	valueLogger = log.NewLogger("xmp.value")
)

var (
	// ErrChildFieldNotValid indicates that the given complex type has no child
	// field with the given name.
	ErrChildFieldNotFound = errors.New("field not found")
)

func parseValue(namespace xmpnamespace.Namespace, fieldName string, rawValue string) (parsedValue interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ft, found := namespace.Fields[fieldName]
	if found == false {
		fmt.Printf("Could not find field [%s] in namespace [%s] (parseValue)\n", fieldName, namespace.PreferredPrefix)
		return nil, ErrChildFieldNotFound
	}

	namespaceUri := namespace.Uri

	// TODO(dustin): !! Drop support for complex types and field-mapping type to be map[string]ScalarFieldType.
	sft, ok := ft.(xmptype.ScalarFieldType)
	if ok == false {
		log.Panicf("scalar value field did not return a scalar parser: NS=[%s] FIELD=[%s] TYPE=[%v]", namespaceUri, fieldName, reflect.TypeOf(ft))
	}

	parser := sft.GetValueParser(rawValue)

	parsedValue, err = parser.Parse()
	if err != nil {
		valueLogger.Warningf(nil, "Could not parse value: NS=[%s] FIELD=[%s] VALUE=[%s]", namespaceUri, fieldName, rawValue)

		fmt.Printf("Could not parse value: NS=[%s] FIELD=[%s] VALUE=[%s]\n", namespaceUri, fieldName, rawValue)

		return nil, err
	}

	return parsedValue, nil
}

func isArrayType(namespace xmpnamespace.Namespace, fieldName string) (flag bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ft, found := namespace.Fields[fieldName]
	if found == false {
		return false, ErrChildFieldNotFound
	}

	_, ok := ft.(xmptype.ArrayType)

	return ok, nil
}
