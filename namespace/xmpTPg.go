package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpTPgUri is the 'xmpTPg' namespace URI made a constant to support
	// testing.
	XmpTPgUri = "http://ns.adobe.com/xap/1.0/t/pg/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpTPgUri,
		PreferredPrefix: "xmpTPg",
		Fields: map[string]interface{}{
			// NOTE(dustin): Not implemented
			// "Colorants":,
			// "Fonts":,
			// "MaxPageSize":
			"NPages":               xmptype.IntegerFieldType{},
			"PlateNames":           xmptype.TextFieldType{},
			"absPeakAudioFilePath": xmptype.UriFieldType{},
			"album":                xmptype.TextFieldType{},
			"altTapeName":          xmptype.TextFieldType{},
			// NOTE(dustin): Not implemented
			// "AltTimecode":,
			"artist":           xmptype.TextFieldType{},
			"audioChannelType": xmptype.TextFieldType{},
			"audioCompressor":  xmptype.TextFieldType{},
			"audioSampleRate":  xmptype.IntegerFieldType{},
			"audioSampleType":  xmptype.TextFieldType{},
			// NOTE(dustin): Not implemented
			// "beatSpliceParams":,
			"cameraAngle": xmptype.TextFieldType{},
			"cameraLabel": xmptype.TextFieldType{},
			"cameraModel": xmptype.TextFieldType{},
			"cameraMove":  xmptype.TextFieldType{},
			"client":      xmptype.TextFieldType{},
			"comment":     xmptype.TextFieldType{},
			"composer":    xmptype.TextFieldType{},
			// NOTE(dustin): Not implemented
			// "contributedMedia":,
			"director":            xmptype.TextFieldType{},
			"directorPhotography": xmptype.TextFieldType{},
			// NOTE(dustin): Not implemented
			// "duration":            TimeFieldType,
			"engineer":     xmptype.TextFieldType{},
			"fileDataRate": xmptype.RationalFieldType{},
			"genre":        xmptype.TextFieldType{},
			"good":         xmptype.BooleanFieldType{},
			"instrument":   xmptype.TextFieldType{},
		},
	}

	register(namespace)
}
