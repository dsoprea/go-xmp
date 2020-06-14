package xmp

var (
	knownPreferredNamespacePrefixes = map[string]string{
		// NOTE(dustin): gofmt odd alignment here
		"http://www.w3.org/1999/02/22-rdf-syntax-ns#": "rdf",
		"adobe:ns:meta/":                                 "x",
		"http://ns.adobe.com/photoshop/1.0/":             "photoshop",
		"http://ns.adobe.com/xap/1.0/mm/":                "xmpMM",
		"http://purl.org/dc/elements/1.1/":               "dc",
		"http://ns.adobe.com/xap/1.0/":                   "xmp",
		"http://ns.microsoft.com/photo/1.0/":             "MicrosoftPhoto",
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

func LookupPreferredNamespacePrefix(namespace string) string {
	return knownPreferredNamespacePrefixes[namespace]
}
