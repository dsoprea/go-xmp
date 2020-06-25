package xmp

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/registry"
)

func TestNewXmpPropertyIndex(t *testing.T) {
	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})
	if xpi.subindices == nil {
		t.Fatalf("subindices not initialized.")
	} else if xpi.leaves == nil {
		t.Fatalf("leaves not initialized.")
	}
}

func getTestIndex() *XmpPropertyIndex {
	xpi := newXmpPropertyIndex(xmpregistry.XmlName{})

	microsoftphotoNamespaceUri := "http://ns.microsoft.com/photo/1.0/"

	name := xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {xmpnamespace.DcUri, "title"}, {xmpnamespace.RdfUri, "Alt"}, {xmpnamespace.RdfUri, "li"}}
	value := "Der Goalie bin ig"

	xpi.addScalarValue(name, value)

	name = xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {xmpnamespace.DcUri, "description"}, {xmpnamespace.RdfUri, "Alt"}, {xmpnamespace.RdfUri, "li"}}
	value = "Der Goalie bin ig"

	xpi.addScalarValue(name, value)

	name = xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {xmpnamespace.DcUri, "creator"}, {xmpnamespace.RdfUri, "Seq"}, {xmpnamespace.RdfUri, "li"}}
	value = "CREDIT"

	xpi.addScalarValue(name, value)

	name = xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {xmpnamespace.DcUri, "subject"}, {xmpnamespace.RdfUri, "Bag"}, {xmpnamespace.RdfUri, "li"}}
	value = "tag"

	xpi.addScalarValue(name, value)

	name = xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {microsoftphotoNamespaceUri, "LastKeywordXMP"}, {xmpnamespace.RdfUri, "Bag"}, {xmpnamespace.RdfUri, "li"}}
	value = "tag"

	xpi.addScalarValue(name, value)

	name = xmpregistry.XmpPropertyName{{xmpnamespace.XUri, "xmpmeta"}, {microsoftphotoNamespaceUri, "LastKeywordIPTC"}, {xmpnamespace.RdfUri, "Bag"}, {xmpnamespace.RdfUri, "li"}}
	value = "tag"

	xpi.addScalarValue(name, value)

	return xpi
}

func TestXmpPropertyIndex_Count(t *testing.T) {
	xpi := getTestIndex()

	if xpi.Count() != 6 {
		t.Fatalf("Count not correct: (%d)", xpi.Count())
	}
}

func checkFirstLoadedProperty(t *testing.T, xpi *XmpPropertyIndex) {
	if len(xpi.subindices) != 1 {
		t.Fatalf("Subindices at level 0 not correct.")
	}

	if len(xpi.leaves) != 0 {
		t.Fatalf("Leaves at level 0 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].subindices) != 6 {
		t.Fatalf("Subindices at level 1 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].leaves) != 0 {
		t.Fatalf("Leaves at level 1 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].subindices["[dc]title"].subindices) != 1 {
		t.Fatalf("Subindices at level 2 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].subindices["[dc]title"].leaves) != 0 {
		t.Fatalf("Leaves at level 2 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].subindices["[dc]title"].subindices["[rdf]Alt"].subindices) != 0 {
		t.Fatalf("Subindices at level 3 not correct.")
	}

	if len(xpi.subindices["[x]xmpmeta"].subindices["[dc]title"].subindices["[rdf]Alt"].leaves) != 1 {
		t.Fatalf("Leaves at level 3 not correct.")
	}

	values := xpi.subindices["[x]xmpmeta"].subindices["[dc]title"].subindices["[rdf]Alt"].leaves["[rdf]li"]

	if len(values) != 1 {
		t.Fatalf("Final leaves not correct: %v", values)
	}

	expected := []interface{}{"Der Goalie bin ig"}

	if reflect.DeepEqual(values, expected) != true {
		t.Fatalf("Stored leaf values not correct: %v", values)
	}
}

func TestXmpPropertyIndex_add(t *testing.T) {
	xpi := getTestIndex()

	// Make sure the first one is loaded correctly in the index hierarchy.

	if xpi.Count() != 6 {
		t.Fatalf("Count not correct: (%d)", xpi.Count())
	}

	checkFirstLoadedProperty(t, xpi)

	// 1

	actual, err := xpi.Get([]string{"[x]xmpmeta", "[dc]title", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expected := []interface{}{"Der Goalie bin ig"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (1).")
	}

	// 2

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]description", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"Der Goalie bin ig"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (2).")
	}

	// 3

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]creator", "[rdf]Seq", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"CREDIT"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (3).")
	}

	// 4

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]subject", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (4).")
	}

	// 5

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordXMP", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (5).")
	}

	// 6

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordIPTC", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (6).")
	}
}

func TestXmpPropertyIndex_Get(t *testing.T) {
	xpi := getTestIndex()

	if xpi.Count() != 6 {
		t.Fatalf("Count not correct: (%d)", xpi.Count())
	}

	// 1

	actual, err := xpi.Get([]string{"[x]xmpmeta", "[dc]title", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expected := []interface{}{"Der Goalie bin ig"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (1).")
	}

	// 2

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]description", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"Der Goalie bin ig"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (2).")
	}

	// 3

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]creator", "[rdf]Seq", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"CREDIT"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (3).")
	}

	// 4

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]subject", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (4).")
	}

	// 5

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordXMP", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (5).")
	}

	// 6

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordIPTC", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expected = []interface{}{"tag"}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (6).")
	}
}

func TestXmpPropertyIndex_Dump(t *testing.T) {
	xpi := getTestIndex()
	xpi.Dump()
}

func TestXmpPropertyIndex_dump(t *testing.T) {
	xpi := getTestIndex()
	xpi.dump([]string{})
}
