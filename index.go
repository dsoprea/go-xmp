package xmp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/namespace"
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
	nodeName   XmlName
	subindices map[string]*XmpPropertyIndex
	leaves     map[string][]ValueParser
}

func newXmpPropertyIndex(nodeName XmlName) *XmpPropertyIndex {
	subindices := make(map[string]*XmpPropertyIndex)
	leaves := make(map[string][]ValueParser)

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

func (xpi *XmpPropertyIndex) name() XmlName {
	return xpi.nodeName
}

type scalarLeafNode struct {
	namespace   xmpnamespace.Namespace
	fieldName   string
	parsedValue interface{}
}

func (sln scalarLeafNode) Parse() (parsed interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ft, found := sln.namespace.Fields[sln.fieldName]
	if found == false {
		fmt.Printf("ErrFieldNotFound (1): [%s]\n", sln.namespace)
		return nil, ErrFieldNotFound
	}

	sft, ok := ft.(xmptype.ScalarFieldType)
	if ok == false {
		fmt.Printf("ErrFieldNotFound (2)\n")
		return nil, ErrFieldNotFound
	}

	svp := sft.GetValueParser(sln.rawValue)

	parsed, err = svp.Parse()
	if err != nil {
		if err == xmptype.ErrValueNotValid {
			indexLogger.Warningf(nil, "Could not parse SCALAR attribute [%s] [%s]: [%s]", sln.namespace, sln.fieldName, sln.rawValue)
			return nil, err
		}

		log.Panic(err)
	}

	return parsed, nil
}

func (xpi *XmpPropertyIndex) addArrayValue(name XmpPropertyName, array xmptype.ArrayValue) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): !! Finish

	return nil
}

func (xpi *XmpPropertyIndex) addScalarValue(name XmpPropertyName, parsedValue interface{}) (err error) {
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
		namespace, err := xmpnamespace.Get(currentNodeName.Space)
		log.PanicIf(err)

		sln := scalarLeafNode{
			namespace:   namespace,
			fieldName:   currentNodeName.Local,
			parsedValue: parsedValue,
		}

		if currentLeaves, found := xpi.leaves[currentNodeNamePhrase]; found == true {
			xpi.leaves[currentNodeNamePhrase] = append(currentLeaves, sln)
		} else {
			xpi.leaves[currentNodeNamePhrase] = []ValueParser{sln}
		}
	}

	return nil
}

// Get searches the index for the property with the name represented by the
// string slice.
func (xpi *XmpPropertyIndex) Get(namePhraseSlice []string) (results []ValueParser, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	currentNodeNamePhrase := namePhraseSlice[0]

	if len(namePhraseSlice) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			fmt.Printf("ErrFieldNotFound (3)\n")
			return nil, ErrFieldNotFound
		}

		values, err := subindex.Get(namePhraseSlice[1:])
		if err != nil {
			if err == ErrFieldNotFound {
				fmt.Printf("ErrFieldNotFound (4)\n")
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

	fmt.Printf("ErrFieldNotFound (5)\n")

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
			fmt.Printf("%s\n", fqNamePhrase)
			parsed, err := value.Parse()
			if err != nil {
				fmt.Printf("- Error: [%s]\n", err.Error())

				if err == ErrFieldNotFound || err == xmptype.ErrValueNotValid {
					indexLogger.Warningf(nil, "Not dumping value for [%s]: [%s]", fqNamePhrase, err.Error())
					continue
				}

				log.Panic(err)
			}

			fmt.Printf("%s = [%s]\n", fqNamePhrase, parsed)
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
