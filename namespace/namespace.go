package xmpnamespace

import (
	"errors"

	"encoding/xml"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrNamespaceNotFound indicates that a namespace was requested that is
	// not registered.
	ErrNamespaceNotFound = errors.New("namespace not found")

	// ErrFieldNotFound indicates that a field was not found for a specific
	// namespace.
	ErrFieldNotFound = errors.New("field not found")
)

// Namespace describes the information about a single namespace.
type Namespace struct {
	// Uri is the URI of a namespace (it should be regarded as a string only;
	// XML namespaces are not necssarily valid Internet resources).
	Uri string

	// PreferredPrefix is the preferred naming-prefix prescribed by the
	// governing standard of this namespace.
	PreferredPrefix string

	// Fields is a mapping of field names to types.
	Fields map[string]interface{}
}

var (
	namespaces = make(map[string]Namespace)
)

func register(namespace Namespace) {
	if _, found := namespaces[namespace.Uri]; found == true {
		log.Panicf("namespace already registered: [%s]", namespace.Uri)
	}

	namespaces[namespace.Uri] = namespace
}

// Get returns the namespace registration for the given URI. Since namespaces
// URIs are strictly defined, no normalization is required.
func Get(uri string) (namespace Namespace, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	namespace, found := namespaces[uri]

	if found == false {
		return Namespace{}, ErrNamespaceNotFound
	}

	return namespace, nil
}

// MustGet returns the Namespace struct associated with the given URI. It panics
// if not known.
func MustGet(uri string) (namespace Namespace) {
	namespace, err := Get(uri)
	if err != nil {
		panic(err)
	}

	return namespace
}

var (
	cachedLookups = make(map[xml.Name]interface{})
)

// GetFieldType returns the field-type for a specific `xml.Name`.
func GetFieldType(name xml.Name) (ft interface{}, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	ft, found := cachedLookups[name]
	if found == true {
		return ft, nil
	}

	namespace, err := Get(name.Space)
	if err != nil {
		if err == ErrNamespaceNotFound {
			return 0, err
		}

		log.Panic(err)
	}

	ft, found = namespace.Fields[name.Local]
	if found == false {
		return ft, ErrFieldNotFound
	}

	cachedLookups[name] = ft

	return ft, nil
}

// MustGetFieldType returns the field-type for a specific `xml.Name`. It panics
// if not known.
func MustGetFieldType(name xml.Name) (ft interface{}) {
	ft, err := GetFieldType(name)
	log.PanicIf(err)

	return ft
}
