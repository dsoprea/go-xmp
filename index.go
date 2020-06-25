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

type scalarLeafNode struct {
	namespace   xmpregistry.Namespace
	fieldName   string
	parsedValue interface{}
}

func (xpi *XmpPropertyIndex) addArrayValue(name xmpregistry.XmpPropertyName, array xmptype.ArrayValue) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): !! Just for exploration/debugging.
	stringLister, ok := array.(xmptype.ArrayStringValueLister)
	if ok == true {
		items, err := stringLister.Items()
		log.PanicIf(err)

		items = items

		fmt.Printf("Indexing array value: [%s]\n", name)

		for i, value := range items {
			fmt.Printf("(%d)> [%s]\n", i, value)
		}

		fmt.Printf("\n")
	}

	// TODO(dustin): !! Finish

	return nil
}

func (xpi *XmpPropertyIndex) addScalarValue(name xmpregistry.XmpPropertyName, parsedValue interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	currentNodeName := name[0]
	currentNodeNamePhrase := currentNodeName.String()

	if len(name) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			subindex = newXmpPropertyIndex(currentNodeName)
		}

		err := subindex.addScalarValue(name[1:], parsedValue)
		log.PanicIf(err)

		if found == false {
			xpi.subindices[currentNodeNamePhrase] = subindex
		}
	} else {
		namespace, err := xmpregistry.Get(currentNodeName.Space)
		log.PanicIf(err)

		sln := scalarLeafNode{
			namespace:   namespace,
			fieldName:   currentNodeName.Local,
			parsedValue: parsedValue,
		}

		if currentLeaves, found := xpi.leaves[currentNodeNamePhrase]; found == true {
			xpi.leaves[currentNodeNamePhrase] = append(currentLeaves, sln)
		} else {
			xpi.leaves[currentNodeNamePhrase] = []interface{}{sln}
		}
	}

	return nil
}

func (xpi *XmpPropertyIndex) addComplexNode(xpn xmpregistry.XmpPropertyName, attributes map[xml.Name]interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	fmt.Printf("addComplexNode: %s: %s\n", xpn, xmpregistry.InlineAttributes(attributes))

	// TODO(dustin): !! Finish

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
			fmt.Printf("%s = [%s] [%v]\n", fqNamePhrase, reflect.TypeOf(value), value)
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
