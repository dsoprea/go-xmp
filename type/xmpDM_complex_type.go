package xmptype

const (
	// XmpDmUri is the 'xmpDM' namespace URI made a constant to support testing.
	XmpDmUri = "http://ns.adobe.com/xmp/1.0/DynamicMedia/"
)

var (
	dynamicMediaFields = map[string]interface{}{
		"riseInDecibel": RealFieldType{},

		// Not a scalar type. Irrelevant here.
		// "riseInTimeDuration": TimeFieldType,

		"useFileBeatsMarker": BooleanFieldType{},
		"key":                TextFieldType{},

		// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
		//
		// In XMP Specification (Part 2)...
		// - 1.2.6.2 text
		// - 1.2.6.9 integer
		//
		// "value":   TextFieldType{},
		"comment": TextFieldType{},

		// Not a scalar type. Irrelevant here.
		// "cuePointParams":,

		"cuePointType": TextFieldType{},

		// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
		//
		// In XMP Specification (Part 2)...
		// - 1.2.6.5 FrameCount
		// - 1.2.6.6 Time
		// - 2.5 Time
		//
		// "duration":     FrameCountFieldType{},
		"location":    UriFieldType{},
		"name":        TextFieldType{},
		"probability": RealFieldType{},
		"speaker":     TextFieldType{},

		// TODO(dustin): Contradiction in specification. We need to wait to observe real data. Similar to other comment, below.
		//
		// In XMP Specification (Part 2)...
		// - 1.2.6.5 FrameCount
		// - 1.2.6.6 Time
		//
		// "startTime":    FrameCountFieldType{},
		"target":       TextFieldType{},
		"type":         OpenChoiceFieldType{},
		"managed":      BooleanFieldType{},
		"path":         UriFieldType{},
		"track":        TextFieldType{},
		"webStatement": UriFieldType{},

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

		"scale": RationalFieldType{},
		// "value":                      IntegerFieldType{},
		"timeFormat":                 ClosedChoiceFieldType{},
		"timeValue":                  TextFieldType{},
		"frameOverlappingPercentage": RealFieldType{},
		"frameSize":                  RealFieldType{},
		"quality":                    ClosedChoiceFieldType{},
		"frameRate":                  FrameRateFieldType{},

		// Not a scalar type. Irrelevant here.
		// "markers":,
		"trackName": TextFieldType{},

		// TODO(dustin): Looks like this should be an inline list parsed and checked against an actual list. It'd help to see an example.
		// "trackType": ChoiceFieldType{},
	}

	dynamicMediaNamespace = Namespace{
		Uri:             XmpDmUri,
		PreferredPrefix: "xmpDM",
	}
)

// DynamicMediaFieldType encapsulates a set of fields that refer to a job.
type DynamicMediaFieldType struct {
}

// ChildFieldType returns the field-type for the immediate child with the
// given name.
func (dmft DynamicMediaFieldType) ChildFieldType(fieldName string) (ft interface{}, err error) {
	ft, found := dynamicMediaFields[fieldName]
	if found == false {
		return nil, ErrChildFieldNotValid
	}

	return ft, nil
}

// Namespace returns the namespace info the node/children of this type.
func (dmft DynamicMediaFieldType) Namespace() Namespace {
	return dynamicMediaNamespace
}

func init() {
	registerComplex(DynamicMediaFieldType{})
}
