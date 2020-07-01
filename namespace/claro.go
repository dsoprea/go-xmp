package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// ClaroUri is the 'claro' namespace URI made a constant to support testing.
	ClaroUri = "http://www.elpical.com/claro/synt1.0/"
)

var (
	ClaroNamespace = xmpregistry.Namespace{
		Uri:             ClaroUri,
		PreferredPrefix: "claro",
		Fields: map[string]interface{}{
			"Logging": xmptype.OrderedTextArrayFieldType{},
		},
	}
)

func init() {
	xmpregistry.Register(ClaroNamespace)
}
