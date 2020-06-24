package xmpnamespace

import (
	"github.com/dsoprea/go-xmp/type"
)

const (
	// XmpDmUri is the 'xmpDM' namespace URI made a constant to support testing.
	XmpDmUri = "http://ns.adobe.com/xmp/1.0/DynamicMedia/"
)

func init() {
	namespace := Namespace{
		Uri:             XmpDmUri,
		PreferredPrefix: "xmpDM",
		Fields: map[string]interface{}{
			"riseInDecibel": xmptype.RealFieldType{},

			// Not a scalar type. Irrelevant here.
			// "riseInTimeDuration": TimeFieldType,

			"useFileBeatsMarker": xmptype.BooleanFieldType{},
			"key":                xmptype.TextFieldType{},

			// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
			//
			// In XMP Specification (Part 2)...
			// - 1.2.6.2 text
			// - 1.2.6.9 integer
			//
			// "value":   TextFieldType{},
			"comment": xmptype.TextFieldType{},

			// Not a scalar type. Irrelevant here.
			// "cuePointParams":,

			"cuePointType": xmptype.TextFieldType{},

			// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
			//
			// In XMP Specification (Part 2)...
			// - 1.2.6.5 FrameCount
			// - 1.2.6.6 Time
			// - 2.5 Time
			//
			// "duration":     FrameCountFieldType{},
			"location":    xmptype.UriFieldType{},
			"name":        xmptype.TextFieldType{},
			"probability": xmptype.RealFieldType{},
			"speaker":     xmptype.TextFieldType{},

			// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
			//
			// In XMP Specification (Part 2)...
			// - 1.2.6.5 FrameCount
			// - 1.2.6.6 Time
			//
			// "startTime":    FrameCountFieldType{},
			"target":       xmptype.TextFieldType{},
			"type":         xmptype.OpenChoiceFieldType{},
			"managed":      xmptype.BooleanFieldType{},
			"path":         xmptype.UriFieldType{},
			"track":        xmptype.TextFieldType{},
			"webStatement": xmptype.UriFieldType{},

			// TODO(dustin): We don't yet have enough information to handle this
			// due to a contradiction.
			//
			// In XMP Specification (Part 2)...
			// - 1.2.6.5, it's defined as a "open choice" or "comma-delimited list".
			// - 1.2.6.7, it's defined as a "closed choice".
			//
			// If this can have different types of values based on context, it's not immediately clear how to distinguish, and such contradictions are rare occurrences. Therefore, we will just ignore it and come back later if this is actually used and someone complains.
			//
			// "type":                       ??{},

			"scale": xmptype.RationalFieldType{},
			// "value":                      IntegerFieldType{},
			"timeFormat":                 xmptype.ClosedChoiceFieldType{},
			"timeValue":                  xmptype.TextFieldType{},
			"frameOverlappingPercentage": xmptype.RealFieldType{},
			"frameSize":                  xmptype.RealFieldType{},
			"quality":                    xmptype.ClosedChoiceFieldType{},
			"frameRate":                  xmptype.FrameRateFieldType{},

			// Not a scalar type. Irrelevant here.
			// "markers":,
			"trackName": xmptype.TextFieldType{},

			// TODO(dustin): Looks like this should be an inline list parsed and checked against an actual list. It'd help to see an example.
			// "trackType": ChoiceFieldType{},
		},
	}

	register(namespace)
}
