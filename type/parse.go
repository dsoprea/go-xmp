package xmptype

import (
	"errors"
	"fmt"
	"reflect"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

var (
	parseLogger = log.NewLogger("xmptype.parse")
)

var (
	// ErrChildFieldNotValid indicates that the given complex type has no child
	// field with the given name.
	ErrChildFieldNotFound = errors.New("field not found")
)

func ParseValue(namespace xmpregistry.Namespace, fieldName string, rawValue string) (parsedValue interface{}, err error) {
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
	sft, ok := ft.(ScalarFieldType)
	if ok == false {
		log.Panicf("scalar value field did not return a scalar parser: NS=[%s] FIELD=[%s] TYPE=[%v]", namespaceUri, fieldName, reflect.TypeOf(ft))
	}

	parser := sft.GetValueParser(rawValue)

	parsedValue, err = parser.Parse()
	if err != nil {
		parseLogger.Warningf(nil, "Could not parse value: NS=[%s] FIELD=[%s] VALUE=[%s]", namespaceUri, fieldName, rawValue)

		fmt.Printf("Could not parse value: NS=[%s] FIELD=[%s] VALUE=[%s]\n", namespaceUri, fieldName, rawValue)

		return nil, err
	}

	return parsedValue, nil
}

func IsArrayType(namespace xmpregistry.Namespace, fieldName string) (flag bool, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ft, found := namespace.Fields[fieldName]
	if found == false {
		return false, ErrChildFieldNotFound
	}

	_, ok := ft.(ArrayType)

	return ok, nil
}

func ParseAttributes(se xml.StartElement) (attributes map[xml.Name]interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	attributes = make(map[xml.Name]interface{})

	for _, attribute := range se.Attr {
		attributeNamespaceUri := attribute.Name.Space
		attributeLocalName := attribute.Name.Local
		attributeRawValue := attribute.Value

		attributeNamespace, err := xmpregistry.Get(attributeNamespaceUri)
		if err != nil {
			if err == xmpregistry.ErrNamespaceNotFound {
				// TODO(dustin): Add a package-level method to get namespaces or warn once if not available.

				// if _, found := xp.unknownNamespaces[attributeNamespaceUri]; found == false {
				parseLogger.Warningf(
					nil,
					"Namespace [%s] for attribute [%s] is not known. Skipping.",
					attributeNamespaceUri, attributeLocalName)

				// xp.unknownNamespaces[attributeNamespaceUri] = struct{}{}
				// }

				continue
			}

			log.Panic(err)
		}

		parsedValue, err := ParseValue(attributeNamespace, attributeLocalName, attributeRawValue)
		if err != nil {
			if err == ErrChildFieldNotFound || err == ErrValueNotValid {
				parseLogger.Warningf(
					nil,
					"Could not parse attribute [%s] [%s] value: [%s]",
					attributeNamespaceUri, attributeLocalName, attributeRawValue)

				continue
			}

			log.Panic(err)
		}

		attributes[attribute.Name] = parsedValue
	}

	return attributes, nil
}
