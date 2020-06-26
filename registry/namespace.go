package xmpregistry

import (
	"errors"
	"fmt"

	"github.com/dsoprea/go-logging"
)

var (
	namespaceLogger = log.NewLogger("xmpregistry.namespace")
)

var (
	// ErrNamespaceNotFound indicates that a namespace was requested that is
	// not registered.
	ErrNamespaceNotFound = errors.New("namespace not found")
)

var (
	// namespaces contains all of the namespace registrations.
	namespaces = make(map[string]Namespace)

	// unknownNamespaces indicates which namespaces have been looked up that we
	// don't have registrations for. This allows us to log warnings once and
	// only once.
	unknownNamespaces = make(map[string]struct{})
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

// String returns a string representation of the namespace.
func (namespace Namespace) String() string {
	return fmt.Sprintf("Namespace<URI=[%s] PREFIX=[%s]>", namespace.Uri, namespace.PreferredPrefix)
}

// Register registers a namespace for access during parsing and indexing.
func Register(namespace Namespace) {
	if _, found := namespaces[namespace.Uri]; found == true {
		log.Panicf("namespace already registered: [%s]", namespace.Uri)
	}

	namespaces[namespace.Uri] = namespace
}

// Clear removes all namespace registrations. Supports testing.
func Clear() {
	namespaces = make(map[string]Namespace)
	unknownNamespaces = make(map[string]struct{})
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
		if err == ErrNamespaceNotFound {
			if _, found := unknownNamespaces[uri]; found == false {
				namespaceLogger.Warningf(
					nil,
					"Namespace [%s] was requested but is not known.",
					uri)

				unknownNamespaces[uri] = struct{}{}
			}
		}

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
