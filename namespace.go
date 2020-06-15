package xmp

import (
	"github.com/dsoprea/go-xmp/namespace"
)

const (
	photoshopNamespaceUri      = "http://ns.adobe.com/photoshop/1.0/"
	microsoftphotoNamespaceUri = "http://ns.microsoft.com/photo/1.0/"
	claroNamespaceUri          = "http://www.elpical.com/claro/synt1.0/"
	xmptpgNamespaceUri         = "http://ns.adobe.com/xap/1.0/t/pg/"
	xmpdmNamespaceUri          = "http://ns.adobe.com/xmp/1.0/DynamicMedia/"
	pdfNamespaceUri            = "http://ns.adobe.com/pdf/1.3/"
	crsNamespaceUri            = "http://ns.adobe.com/camera-raw-settings/1.0/"
	strefNamespaceUri          = "http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
)

var (
	// knownPreferredNamespacePrefixes describes all of the recommended prefixes
	// for the standard (and then some) XMP namespaces. This is largely found in
	// the XMP Specification parts 1 and 2.
	knownPreferredNamespacePrefixes = map[string]string{
		photoshopNamespaceUri:      "photoshop",
		microsoftphotoNamespaceUri: "MicrosoftPhoto",
		claroNamespaceUri:          "claro",
		xmptpgNamespaceUri:         "xmpTPg",
		xmpdmNamespaceUri:          "xmpDM",
		pdfNamespaceUri:            "pdf",
		crsNamespaceUri:            "crs",
		strefNamespaceUri:          "stRef",
	}
)

// LookupPreferredNamespacePrefix looks for the common prefix associated with
// the given namespace (namespaces are very strictly named, so normalization is
// not required). Returns an empty-string if not found.
func LookupPreferredNamespacePrefix(uri string) string {

	// TODO(dustin): !! This is a bridge to allow us to incrementally move namespaces from here to the xmlnamespace package. After we're done with this, we should remove this function (and file).
	namespace, err := xmpnamespace.Get(uri)
	if err == nil {
		return namespace.PreferredPrefix
	}

	return knownPreferredNamespacePrefixes[uri]
}
