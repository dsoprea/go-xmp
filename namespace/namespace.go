package xmpnamespace

import (
	"errors"

	"github.com/dsoprea/go-logging"
)

var (
	// ErrNamespaceNotFound indicates that a namespace was requested that is
	// not registered.
	ErrNamespaceNotFound = errors.New("namespace not found")
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
	Fields map[string]FieldType
}

var (
	namespaces map[string]Namespace
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
		return namespace, ErrNamespaceNotFound
	}

	return namespace, nil
}

func init() {
	namespaces = make(map[string]Namespace)
}