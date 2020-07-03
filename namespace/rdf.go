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

	// RdfTag is the name for the "RDF" tag.
	RdfTag = xml.Name{
		Space: RdfUri,
		Local: "RDF",
	}

	// RdfDescriptionTag is the name for the "Description" tag.
	RdfDescriptionTag = xml.Name{
		Space: RdfUri,
		Local: "Description",
	}

	// RdfLiTag is the name for the "li" tag.
	RdfLiTag = xml.Name{
		Space: RdfUri,
		Local: "li",
	}

	// RdfNamespace is the namespace descriptor for "rdf". We do not define any
	// fields for it because it defined no leaf nodes [that we have encountered]
	// and therefore we require no parsing and no knowledge of types.
	RdfNamespace = xmpregistry.Namespace{
		Uri:             RdfUri,
		PreferredPrefix: "rdf",
	}
)

func init() {
	xmpregistry.Register(RdfNamespace)
}
