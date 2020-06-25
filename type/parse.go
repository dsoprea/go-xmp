package xmptype

import (
	"errors"
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
		return nil, ErrChildFieldNotFound
	}

	namespaceUri := namespace.Uri

	sft, ok := ft.(ScalarFieldType)
	if ok == false {
		log.Panicf("scalar value field did not return a scalar parser: NS=[%s] FIELD=[%s] TYPE=[%v]", namespaceUri, fieldName, reflect.TypeOf(ft))
	}

	parser := sft.GetValueParser(rawValue)

	parsedValue, err = parser.Parse()
	if err != nil {
		parseLogger.Warningf(nil, "Could not parse value: NS=[%s] FIELD=[%s] VALUE=[%s]", namespaceUri, fieldName, rawValue)
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
