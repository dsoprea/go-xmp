package xmp

const (
	rdfNamespaceUri            = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	xNamespaceUri              = "adobe:ns:meta/"
	photoshopNamespaceUri      = "http://ns.adobe.com/photoshop/1.0/"
	xmpmmNamespaceUri          = "http://ns.adobe.com/xap/1.0/mm/"
	dcNamespaceUri             = "http://purl.org/dc/elements/1.1/"
	xmpNamespaceUri            = "http://ns.adobe.com/xap/1.0/"
	microsoftphotoNamespaceUri = "http://ns.microsoft.com/photo/1.0/"
	claroNamespaceUri          = "http://www.elpical.com/claro/synt1.0/"
	xmpbjNamespaceUri          = "http://ns.adobe.com/xap/1.0/bj/"
	xmptpgNamespaceUri         = "http://ns.adobe.com/xap/1.0/t/pg/"
	xmpdmNamespaceUri          = "http://ns.adobe.com/xmp/1.0/DynamicMedia/"
	pdfNamespaceUri            = "http://ns.adobe.com/pdf/1.3/"
	crsNamespaceUri            = "http://ns.adobe.com/camera-raw-settings/1.0/"
	strefNamespaceUri          = "http://ns.adobe.com/xap/1.0/sType/ResourceRef#"
	xmprightsNamespaceUri      = "http://ns.adobe.com/xap/1.0/rights/"
	xmpidqNamespaceUri         = "http://ns.adobe.com/xmp/Identifier/qual/1.0/"
)

var (
	// knownPreferredNamespacePrefixes describes all of the recommended prefixes
	// for the standard (and then some) XMP namespaces. This is largely found in
	// the XMP Specification parts 1 and 2.
	knownPreferredNamespacePrefixes = map[string]string{
		// NOTE(dustin): gofmt odd alignment here
		rdfNamespaceUri:            "rdf",
		xNamespaceUri:              "x",
		photoshopNamespaceUri:      "photoshop",
		xmpmmNamespaceUri:          "xmpMM",
		dcNamespaceUri:             "dc",
		xmpNamespaceUri:            "xmp",
		microsoftphotoNamespaceUri: "MicrosoftPhoto",
		claroNamespaceUri:          "claro",
		xmpbjNamespaceUri:          "xmpBJ",
		xmptpgNamespaceUri:         "xmpTPg",
		xmpdmNamespaceUri:          "xmpDM",
		pdfNamespaceUri:            "pdf",
		crsNamespaceUri:            "crs",
		strefNamespaceUri:          "stRef",
		xmprightsNamespaceUri:      "xmpRights",
		xmpidqNamespaceUri:         "xmpidq",
	}
)

// LookupPreferredNamespacePrefix looks for the common prefix associated with
// the given namespace (namespaces are very strictly named, so normalization is
// not required). Returns an empty-string if not found.
func LookupPreferredNamespacePrefix(namespace string) string {
	return knownPreferredNamespacePrefixes[namespace]
}
