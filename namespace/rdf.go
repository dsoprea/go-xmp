package xmpnamespace

import (
	"encoding/xml"

	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// RdfUri is the 'rdf' namespace URI made a constant to support testing.
	RdfUri = xmptype.RdfUri
)

var (
	// These assist us during parsing.

	RdfTag = xml.Name{
		Space: RdfUri,
		Local: "RDF",
	}

	RdfDescriptionTag = xml.Name{
		Space: RdfUri,
		Local: "Description",
	}

	RdfLiTag = xml.Name{
		Space: RdfUri,
		Local: "li",
	}
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             RdfUri,
		PreferredPrefix: "rdf",
	}

	xmpregistry.Register(namespace)
}
