package xmpregistry

import (
	"fmt"
	"strings"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
)

var (
	typeLogger = log.NewLogger("xmpregistry.type")
)

var (
	cachedPrefixes = make(map[string]string)
)

func ClearCachedPrefixes() {
	cachedPrefixes = make(map[string]string)
}

// XmlName is a localized version of xml.Name with a String() method attached.
type XmlName xml.Name

func (xn XmlName) Prefix() string {
	prefix, found := cachedPrefixes[xn.Space]
	if found == true {
		return prefix
	}

	ns, err := Get(xn.Space)
	if err != nil {
		// They should notify us of the unknown namespace so that we
		// can register it and they can handle it properly.

		typeLogger.Warningf(nil, "Namespace [%s] is not registered.", xn.Space)

		prefix = "?"
	} else {
		prefix = ns.PreferredPrefix
	}

	cachedPrefixes[xn.Space] = prefix

	return prefix
}

func (xn XmlName) String() string {
	prefix := xn.Prefix()

	return fmt.Sprintf("[%s]%s", prefix, xn.Local)
}

// XmpPropertyName is a series of constituent parts comprising a property's
// fully-qualified name.
type XmpPropertyName []XmlName

// Parts returns a slice of stringifications of the constituent names.
func (xpn XmpPropertyName) Parts() (parts []string) {
	parts = make([]string, len(xpn))
	for i, tag := range xpn {
		parts[i] = tag.String()
	}

	return parts
}

// String returns a string-representation of the name slice.
func (xpn XmpPropertyName) String() string {
	parts := xpn.Parts()
	return strings.Join(parts, ".")
}

func InlineAttributes(attributes map[xml.Name]interface{}) string {
	parts := make([]string, 0, len(attributes))
	for name, parsedValue := range attributes {
		xn := XmlName(name)
		parts = append(parts, fmt.Sprintf("%s=[%v]", xn, parsedValue))
	}

	return strings.Join(parts, " ")
}
