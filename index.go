package xmp

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

var (
	indexLogger = log.NewLogger("xmp.index")
)

var (
	// ErrFieldNotFound represents an error for a get operation that produced
	// no results.
	ErrFieldNotFound = errors.New("field not found in namespace")
)

type ValueParser interface {
	Parse() (parsed interface{}, err error)
}

// XmpPropertyIndex allows for lookups and browsing of found properties.
type XmpPropertyIndex struct {
	nodeName   xmpregistry.XmlName
	subindices map[string]*XmpPropertyIndex
	leaves     map[string][]interface{}
}

func newXmpPropertyIndex(nodeName xmpregistry.XmlName) *XmpPropertyIndex {
	subindices := make(map[string]*XmpPropertyIndex)
	leaves := make(map[string][]interface{})

	xpi := &XmpPropertyIndex{
		nodeName:   nodeName,
		subindices: subindices,
		leaves:     leaves,
	}

	return xpi
}

func (xpi *XmpPropertyIndex) isRoot() bool {
	return xpi.nodeName.Local == ""
}

func (xpi *XmpPropertyIndex) name() xmpregistry.XmlName {
	return xpi.nodeName
}

func (xpi *XmpPropertyIndex) addValue(xpn xmpregistry.XmpPropertyName, value interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	currentNodeName := xpn[0]
	currentNodeNamePhrase := currentNodeName.String()

	if len(xpn) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			subindex = newXmpPropertyIndex(currentNodeName)
		}

		err := subindex.addValue(xpn[1:], value)
		log.PanicIf(err)

		if found == false {
			xpi.subindices[currentNodeNamePhrase] = subindex
		}
	} else {
		if currentLeaves, found := xpi.leaves[currentNodeNamePhrase]; found == true {
			xpi.leaves[currentNodeNamePhrase] = append(currentLeaves, value)
		} else {
			xpi.leaves[currentNodeNamePhrase] = []interface{}{value}
		}
	}

	return nil
}

func (xpi *XmpPropertyIndex) addArrayValue(xpn xmpregistry.XmpPropertyName, array xmptype.ArrayValue) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = xpi.addValue(xpn, array)
	log.PanicIf(err)

	return nil
}

type ScalarLeafNode struct {
	Namespace   xmpregistry.Namespace
	FieldName   string
	ParsedValue interface{}
}

func (xpi *XmpPropertyIndex) addScalarValue(xpn xmpregistry.XmpPropertyName, parsedValue interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	currentNodeName := xpn[len(xpn)-1]

	namespace, err := xmpregistry.Get(currentNodeName.Space)
	log.PanicIf(err)

	sln := ScalarLeafNode{
		Namespace:   namespace,
		FieldName:   currentNodeName.Local,
		ParsedValue: parsedValue,
	}

	err = xpi.addValue(xpn, sln)
	log.PanicIf(err)

	return nil
}

type ComplexLeafNode map[xml.Name]interface{}

func (xpi *XmpPropertyIndex) addComplexValue(xpn xmpregistry.XmpPropertyName, attributes map[xml.Name]interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	cn := ComplexLeafNode(attributes)

	err = xpi.addValue(xpn, cn)
	log.PanicIf(err)

	return nil
}

// Get searches the index for the property with the name represented by the
// string slice.
func (xpi *XmpPropertyIndex) Get(namePhraseSlice []string) (results []interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	currentNodeNamePhrase := namePhraseSlice[0]

	if len(namePhraseSlice) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			return nil, ErrFieldNotFound
		}

		values, err := subindex.Get(namePhraseSlice[1:])
		if err != nil {
			if err == ErrFieldNotFound {
				return nil, err
			}

			log.Panic(err)
		}

		return values, nil
	}

	// If we get here, we are expecting to find a leaf-node.

	if wrappedValues, found := xpi.leaves[currentNodeNamePhrase]; found == true {
		return wrappedValues, nil
	}

	return nil, ErrFieldNotFound
}

func (xpi *XmpPropertyIndex) dump(prefix []string) {
	for name, subindex := range xpi.subindices {
		subindex.dump(append(prefix, name))
	}

	for name, values := range xpi.leaves {
		fqName := append(prefix, name)
		fqNamePhrase := strings.Join(fqName, ".")

		for _, value := range values {
			if sl, ok := value.(xmptype.ArrayStringValueLister); ok == true {
				items, err := sl.Items()
				log.PanicIf(err)

				fmt.Printf("%s:\n\n  ARRAY [%s]\n  COUNT (%d)\n", fqNamePhrase, reflect.TypeOf(sl), len(items))
				fmt.Printf("\n")

				for i, item := range items {
					fmt.Printf("  Item (%d): [%s] %v\n", i, reflect.TypeOf(item), item)
				}

				fmt.Printf("\n")

			} else if sln, ok := value.(ScalarLeafNode); ok == true {
				fmt.Printf("%s:\n\n   SCALAR\n", fqNamePhrase)
				fmt.Printf("\n")

				fmt.Printf("  [%s]%s = [%s] [%v]\n", sln.Namespace.PreferredPrefix, sln.FieldName, reflect.TypeOf(sln.ParsedValue), sln.ParsedValue)

				fmt.Printf("\n")
			} else if cln, ok := value.(ComplexLeafNode); ok == true {
				fmt.Printf("%s:\n\n  COMPLEX\n", fqNamePhrase)
				fmt.Printf("\n")

				for name, value := range cln {
					fmt.Printf("  %s: [%s] [%v]\n", xmpregistry.XmlName(name), reflect.TypeOf(value), value)
				}

				fmt.Printf("\n")
			} else {
				log.Panicf("can not dump unhandled value: [%v]", reflect.TypeOf(value))
			}
		}
	}
}

// Dump prints all of the properties in the index.
func (xpi *XmpPropertyIndex) Dump() {
	xpi.dump([]string{})
}

// Count returns the number of entries.
func (xpi *XmpPropertyIndex) Count() (count int) {
	for _, subindex := range xpi.subindices {
		count += subindex.Count()
	}

	for _, values := range xpi.leaves {
		count += len(values)
	}

	return count
}
