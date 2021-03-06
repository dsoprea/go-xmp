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
	mainLogger = log.NewLogger("xmp.main")
)

var (
	// ErrFieldNotFound represents an error for a get operation that produced
	// no results.
	ErrFieldNotFound = errors.New("node not found in document")
)

// ValueParser knows how to parse raw values.
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

func (xpi *XmpPropertyIndex) exportValue(value interface{}, doPrintSimplified bool) (encoded interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	if ail, ok := value.(xmptype.ArrayItemLister); ok == true {
		items, err := ail.Items()
		log.PanicIf(err)

		ail = ail

		distilled := make([]interface{}, len(items))
		for i, ai := range items {
			attributes := make(map[string]interface{})

			for name, value := range ai.Attributes {
				namePhrase := xmpregistry.XmlName(name).String()
				attributes[namePhrase] = value
			}

			if len(attributes) > 0 {
				// If there are attributes, then only include the char-data if
				// non-empty. Very frequently, values expressed as attributes
				// are not paired with char-data. So, that just adds pollution
				// to the output.

				encodedComplex := map[string]interface{}{
					"Attributes": attributes,
				}

				if ai.CharData != "" {
					encodedComplex["CharData"] = ai.CharData

					distilled[i] = encodedComplex
				} else {
					if doPrintSimplified == true {
						distilled[i] = attributes
					} else {
						distilled[i] = encodedComplex
					}
				}
			} else {
				if doPrintSimplified == true {
					distilled[i] = ai.CharData
				} else {
					distilled[i] = map[string]interface{}{
						"CharData": ai.CharData,
					}
				}
			}
		}

		return distilled, nil
	} else if sl, ok := value.(xmptype.ArrayStringValueLister); ok == true {
		items, err := sl.StringItems()
		log.PanicIf(err)

		return items, nil
	} else if sln, ok := value.(ScalarLeafNode); ok == true {
		return sln.ParsedValue, nil
	} else if cln, ok := value.(ComplexLeafNode); ok == true {
		exported := make(map[string]interface{})

		for name, value := range cln {
			namePhrase := xmpregistry.XmlName(name).String()
			exported[namePhrase] = value
		}

		return exported, nil
	}

	log.Panicf("can not dump unhandled value: [%v]", reflect.TypeOf(value))
	panic(nil)
}

func (xpi *XmpPropertyIndex) export(xmp xmpregistry.XmpPropertyName, doPrintSimplified bool) (exported map[string]interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	exported = make(map[string]interface{})

	for _, subindex := range xpi.subindices {
		currentXpn := append(xmp, xpi.nodeName)

		currentExported, err := subindex.export(currentXpn, doPrintSimplified)
		log.PanicIf(err)

		exportedKey := subindex.nodeName.String()

		exported[exportedKey] = currentExported
	}

	for key, values := range xpi.leaves {
		encodedValues := make([]interface{}, len(values))
		for i, value := range values {
			encodedValue, err := xpi.exportValue(value, doPrintSimplified)
			if err != nil {
				mainLogger.Errorf(nil, err, "%s: Had trouble enumerating array items under leaf (%d).", key, i)
				log.Panic(err)
			}

			encodedValues[i] = encodedValue
		}

		exported[key] = encodedValues
	}

	return exported, nil
}

// Export returns the indexed XMP values as a simple data-structure. It may be
// used to print structured output (e.g. JSON).
//
// SIMPLIFICATION
//
// The values will be maps that may have an "Attributes" field, "CharData"
// field, or both, but this will be further flattened if `doPrintSimplified` is
// true: if attributes are present but not char-data, then just replace this
// with the attributes map; if char-data is present but not attributes then
// just replace this with the char-data string. More often than not, the
// simplified structure for a given particular namespace should be consistent
// enough that tools can make assumptions about this structure, though they may
// need to be tolerant of rare inconsistencies. Still, this makes for a more
// digestable format for both humans and machine.
func (xpi *XmpPropertyIndex) Export(doPrintSimplified bool) (exported map[string]interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	// TODO(dustin): Add test

	rootXpn := make(xmpregistry.XmpPropertyName, 0)

	exported, err = xpi.export(rootXpn, doPrintSimplified)
	log.PanicIf(err)

	return exported, nil
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

	if len(xpn) == 0 {
		log.Panicf("array value must have non-empty property-name")
	}

	err = xpi.addValue(xpn, array)
	log.PanicIf(err)

	return nil
}

// ScalarLeafNode describes a node having a value stored in the index.
type ScalarLeafNode struct {
	Name        xml.Name
	ParsedValue interface{}
}

func (xpi *XmpPropertyIndex) addScalarValue(xpn xmpregistry.XmpPropertyName, parsedValue interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if len(xpn) == 0 {
		log.Panicf("scalar value must have non-empty property-name")
	}

	currentNodeName := xpn[len(xpn)-1]

	sln := ScalarLeafNode{
		Name:        xml.Name(currentNodeName),
		ParsedValue: parsedValue,
	}

	err = xpi.addValue(xpn, sln)
	log.PanicIf(err)

	return nil
}

// ComplexLeafNode describes the attributes of a leaf node.
type ComplexLeafNode map[xml.Name]interface{}

// Get retrieves the given complex node from the index.
func (cln ComplexLeafNode) Get(uri string, local string) (value interface{}, found bool) {
	name := xml.Name{
		Space: uri,
		Local: local,
	}

	value, found = cln[name]

	return value, found
}

func (xpi *XmpPropertyIndex) addComplexValue(xpn xmpregistry.XmpPropertyName, attributes map[xml.Name]interface{}) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	if len(xpn) == 0 {
		log.Panicf("complex value must have non-empty property-name")
	}

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

		for i, value := range values {
			if sl, ok := value.(xmptype.ArrayStringValueLister); ok == true {
				items, err := sl.StringItems()
				if err != nil {
					fmt.Printf("%s: Had trouble enumerating array items under leaf (%d).\n\n", fqNamePhrase, i)
					log.Panic(err)
				}

				fmt.Printf("%s:\n\n  ARRAY [%s]\n  COUNT (%d)\n", fqNamePhrase, reflect.TypeOf(sl), len(items))
				fmt.Printf("\n")

				for i, item := range items {
					fmt.Printf("  Item (%d): [%s] %v\n", i, reflect.TypeOf(item), item)
				}

				fmt.Printf("\n")

			} else if sln, ok := value.(ScalarLeafNode); ok == true {
				fmt.Printf("%s:\n\n   SCALAR\n", fqNamePhrase)
				fmt.Printf("\n")

				namePhrase := xmpregistry.XmlName(sln.Name).String()

				fmt.Printf("  %s = [%s] [%v]\n", namePhrase, reflect.TypeOf(sln.ParsedValue), sln.ParsedValue)

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
