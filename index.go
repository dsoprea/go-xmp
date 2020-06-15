package xmp

import (
	"errors"
	"fmt"
	"strings"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrPropertyNotFound represents an error for a get operation that produced
	// no results.
	ErrPropertyNotFound = errors.New("property not found")
)

// XmlName is a localized version of xml.Name with a String() method attached.
type XmlName xml.Name

func (xn XmlName) String() string {
	prefix := LookupPreferredNamespacePrefix(xn.Space)
	if prefix == "" {
		// They should notify us of the unknown namespace so that we
		// can register it and they can handle it properly.
		prefix = "?"
	}

	return fmt.Sprintf("[%s]%s", prefix, xn.Local)
}

// XmpPropertyName is a series of constituent parts comprising a property's
// fully-qualified name.
type XmpPropertyName []XmlName

// Parts returns a slice of stringifications of the constituent names.
func (xpn XmpPropertyName) Parts() (parts []string) {
	parts = make([]string, len(xpn))
	for i, tag := range xpn {
		prefix := LookupPreferredNamespacePrefix(tag.Space)
		if prefix == "" {
			// They should notify us of the unknown namespace so that we
			// can register it and they can handle it properly.
			prefix = "?"
		}

		parts[i] = fmt.Sprintf("[%s]%s", prefix, tag.Local)
	}

	return parts
}

// String returns a string-representation of the name slice.
func (xpn XmpPropertyName) String() string {
	parts := xpn.Parts()
	return strings.Join(parts, ".")
}

// XmpPropertyIndex allows for lookups and browsing of found properties.
type XmpPropertyIndex struct {
	subindices map[string]*XmpPropertyIndex
	leaves     map[string][]interface{}
}

func newXmpPropertyIndex() *XmpPropertyIndex {
	subindices := make(map[string]*XmpPropertyIndex)
	leaves := make(map[string][]interface{})

	xpi := &XmpPropertyIndex{
		subindices: subindices,
		leaves:     leaves,
	}

	return xpi
}

func (xpi *XmpPropertyIndex) add(name XmpPropertyName, value interface{}) {
	currentNodeName := name[0]
	currentNodeNamePhrase := currentNodeName.String()

	if len(name) > 1 {
		subindex, found := xpi.subindices[currentNodeNamePhrase]

		if found == false {
			subindex = newXmpPropertyIndex()
		}

		subindex.add(name[1:], value)

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
			return nil, ErrPropertyNotFound
		}

		values, err := subindex.Get(namePhraseSlice[1:])
		if err != nil {
			if err == ErrPropertyNotFound {
				return nil, err
			}

			log.Panic(err)
		}

		return values, nil
	}

	// If we get here, we are expecting to find a leaf-node.

	if values, found := xpi.leaves[currentNodeNamePhrase]; found == true {
		return values, nil
	}

	return nil, ErrPropertyNotFound
}

func (xpi *XmpPropertyIndex) dump(prefix []string) {
	for name, subindex := range xpi.subindices {
		subindex.dump(append(prefix, name))
	}

	for name, values := range xpi.leaves {
		fqName := append(prefix, name)
		fqNamePhrase := strings.Join(fqName, ".")

		for _, value := range values {
			fmt.Printf("%s = [%s]\n", fqNamePhrase, value)
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
