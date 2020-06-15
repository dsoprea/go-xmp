package xmp

const (
	// These are constants in order to support testing.

	rdfNamespaceUri            = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	xNamespaceUri              = "adobe:ns:meta/"
	dcNamespaceUri             = "http://purl.org/dc/elements/1.1/"
	microsoftPhotoNamespaceUri = "http://ns.microsoft.com/photo/1.0/"
)

var (
	// knownPreferredNamespacePrefixes describes all of the recommended prefixes
	// for the standard (and then some) XMP namespaces. This is largely found in
	// the XMP Specification parts 1 and 2.
	knownPreferredNamespacePrefixes = map[string]string{
		// NOTE(dustin): gofmt odd alignment here
		rdfNamespaceUri:                                  "rdf",
		xNamespaceUri:                                    "x",
		"http://ns.adobe.com/photoshop/1.0/":             "photoshop",
		"http://ns.adobe.com/xap/1.0/mm/":                "xmpMM",
		dcNamespaceUri:                                   "dc",
		"http://ns.adobe.com/xap/1.0/":                   "xmp",
		microsoftPhotoNamespaceUri:                       "MicrosoftPhoto",
		"http://www.elpical.com/claro/synt1.0/":          "claro",
		"http://ns.adobe.com/xap/1.0/bj/":                "xmpBJ",
		"http://ns.adobe.com/xap/1.0/t/pg/":              "xmpTPg",
		"http://ns.adobe.com/xmp/1.0/DynamicMedia/":      "xmpDM",
		"http://ns.adobe.com/pdf/1.3/":                   "pdf",
		"http://ns.adobe.com/camera-raw-settings/1.0/":   "crs",
		"http://ns.adobe.com/xap/1.0/sType/ResourceRef#": "stRef",
		"http://ns.adobe.com/xap/1.0/rights/":            "xmpRights",
		"http://ns.adobe.com/xmp/Identifier/qual/1.0/":   "xmpidq",
	}
)

// LookupPreferredNamespacePrefix looks for the common prefix associated with
// the given namespace (namespaces are very strictly named, so normalization is
// not required). Returns an empty-string if not found.
func LookupPreferredNamespacePrefix(namespace string) string {
	return knownPreferredNamespacePrefixes[namespace]
}
