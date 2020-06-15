package xmpnamespace

const (
	// XUri is the 'x' namespace URI made a constant to support testing.
	XUri = "adobe:ns:meta/"
)

func init() {
	namespace := Namespace{
		Uri:             XUri,
		PreferredPrefix: "x",
	}

	register(namespace)
}
