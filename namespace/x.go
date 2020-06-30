package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
)

const (
	// XUri is the 'x' namespace URI made a constant to support testing.
	XUri = "adobe:ns:meta/"
)

var (
	XNamespace = xmpregistry.Namespace{
		Uri:             XUri,
		PreferredPrefix: "x",
	}
)

func init() {
	xmpregistry.Register(XNamespace)
}
