package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
)

const (
	// RdfUri is the 'rdf' namespace URI made a constant to support testing.
	RdfUri = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             RdfUri,
		PreferredPrefix: "rdf",
	}

	xmpregistry.Register(namespace)
}
