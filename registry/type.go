package xmpregistry

import (
	"fmt"
	"sort"
	"strings"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
)

var (
	typeLogger = log.NewLogger("xmpregistry.type")
)

var (
	// TODO(dustin): This has questionable savings. It just swaps one lookup with another with maybe only a constant savings.
	cachedPrefixes = make(map[string]string)
)

// ClearCachedPrefixes clears the namespace-prefix cache that is loaded from
// Get().
func ClearCachedPrefixes() {
	cachedPrefixes = make(map[string]string)
}

// XmlName is a localized version of xml.Name with a String() method attached.
type XmlName xml.Name

// Prefix returns the preferred-prefix for the given namespace if registered,
// else "?".
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

// String returns a string representation of the XML name.
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

// InlineAttributes returns all attributes expressed in a single line (for
// logging/dumping values). They are sorted alphabetically to support testing.
func InlineAttributes(attributes map[xml.Name]interface{}) string {
	keys := make(sort.StringSlice, len(attributes))
	mapping := make(map[string]xml.Name)

	i := 0
	for k := range attributes {
		phrase := XmlName(k).String()
		keys[i] = phrase
		mapping[phrase] = k

		i++
	}

	keys.Sort()

	parts := make([]string, 0, len(attributes))
	for _, namePhrase := range keys {
		name := mapping[namePhrase]
		parsedValue := attributes[name]

		xn := XmlName(name)
		parts = append(parts, fmt.Sprintf("%s=[%v]", xn, parsedValue))
	}

	return strings.Join(parts, " ")
}
