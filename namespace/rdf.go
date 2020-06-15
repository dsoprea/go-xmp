package xmpnamespace

const (
	// RdfUri is the 'rdf' namespace URI made a constant to support testing.
	RdfUri = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

func init() {
	namespace := Namespace{
		Uri:             RdfUri,
		PreferredPrefix: "rdf",
	}

	register(namespace)
}
