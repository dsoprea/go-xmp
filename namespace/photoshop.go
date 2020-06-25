package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

const (
	// PhotoshopUri is the 'photoshop' namespace URI made a constant to support
	// testing.
	PhotoshopUri = "http://ns.adobe.com/photoshop/1.0/"
)

func init() {
	namespace := xmpregistry.Namespace{
		Uri:             PhotoshopUri,
		PreferredPrefix: "photoshop",
		Fields: map[string]interface{}{
			"AncestorID":        xmptype.UriFieldType{},
			"LayerName":         xmptype.TextFieldType{},
			"LayerText":         xmptype.TextFieldType{},
			"AuthorsPosition":   xmptype.TextFieldType{},
			"CaptionWriter":     xmptype.TextFieldType{},
			"Category":          xmptype.TextFieldType{},
			"City":              xmptype.TextFieldType{},
			"ColorMode":         xmptype.ClosedChoiceFieldType{},
			"Country":           xmptype.TextFieldType{},
			"Credit":            xmptype.TextFieldType{},
			"DateCreated":       xmptype.DateFieldType{},
			"DocumentAncestors": xmptype.UnorderedAncestorArrayType{},
			"Headline":          xmptype.TextFieldType{},
			"History":           xmptype.TextFieldType{},
			"ICCProfile":        xmptype.TextFieldType{},
			"Instructions":      xmptype.TextFieldType{},
			"Source":            xmptype.TextFieldType{},
			"State":             xmptype.TextFieldType{},

			"SupplementalCategories": xmptype.UnorderedTextArrayType{},

			// TODO(dustin): We need to finish our implementation of Array
			// "TextLayers":

			"TransmissionReference": xmptype.TextFieldType{},
			"Urgency":               xmptype.IntegerFieldType{},
		},
	}

	xmpregistry.Register(namespace)
}
