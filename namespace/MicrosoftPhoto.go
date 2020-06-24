package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// MicrosoftPhotoUri is the 'MicrosoftPhoto' namespace URI made a constant
	// to support testing.
	MicrosoftPhotoUri = "http://ns.microsoft.com/photo/1.0/"
)

func init() {
	namespace := Namespace{
		Uri:             MicrosoftPhotoUri,
		PreferredPrefix: "MicrosoftPhoto",
		Fields: map[string]interface{}{
			"CameraSerialNumber": xmptype.TextFieldType{},
			"DateAcquired":       xmptype.DateFieldType{},
			"FlashManufacturer":  xmptype.TextFieldType{},
			"FlashModel":         xmptype.TextFieldType{},
			"LastKeywordIPTC":    xmptype.TextFieldType{},
			"LastKeywordXMP":     xmptype.TextFieldType{},
			"LensManufacturer":   xmptype.TextFieldType{},
			"LensModel":          xmptype.TextFieldType{},
			"Rating":             xmptype.DateFieldType{},
		},
	}

	register(namespace)
}