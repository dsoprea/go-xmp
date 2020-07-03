package xmp

import (
	"reflect"
	"testing"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/namespace"
	"github.com/dsoprea/go-xmp/registry"
	"github.com/dsoprea/go-xmp/type"
)

var (
	rdfLiName = xml.Name{
		Space: xmpnamespace.RdfUri,
		Local: "li",
	}
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

func constructLiItem(value string) (sln ScalarLeafNode) {
	sln.Name = rdfLiName
	sln.ParsedValue = value

	return sln
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

	actual := values[0].(ScalarLeafNode)

	expected := constructLiItem("Der Goalie bin ig")

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Stored leaf values not correct:\n  Actual: [%s] [%v]\nExpected: [%s] [%v]", reflect.TypeOf(actual), actual, reflect.TypeOf(expected), expected)
	}
}

func TestXmpPropertyIndex_add(t *testing.T) {

	// We process arrays differently than this test implies, though the purpose
	// of this test is to inject several scalars and then retrieve them. Pay not
	// mind to the actual values being pushed.

	xpi := getTestIndex()

	// Make sure the first one is loaded correctly in the index hierarchy.

	if xpi.Count() != 6 {
		t.Fatalf("Count not correct: (%d)", xpi.Count())
	}

	checkFirstLoadedProperty(t, xpi)

	// 1

	actual, err := xpi.Get([]string{"[x]xmpmeta", "[dc]title", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expectedValue := constructLiItem("Der Goalie bin ig")
	expected := []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (1).")
	}

	// 2

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]description", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("Der Goalie bin ig")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (2).")
	}

	// 3

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]creator", "[rdf]Seq", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("CREDIT")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (3).")
	}

	// 4

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]subject", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (4).")
	}

	// 5

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordXMP", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (5).")
	}

	// 6

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordIPTC", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

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

	expectedValue := constructLiItem("Der Goalie bin ig")
	expected := []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (1).")
	}

	// 2

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]description", "[rdf]Alt", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("Der Goalie bin ig")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (2).")
	}

	// 3

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]creator", "[rdf]Seq", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("CREDIT")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (3).")
	}

	// 4

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[dc]subject", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (4).")
	}

	// 5

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordXMP", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Result not correct (5).")
	}

	// 6

	actual, err = xpi.Get([]string{"[x]xmpmeta", "[MicrosoftPhoto]LastKeywordIPTC", "[rdf]Bag", "[rdf]li"})
	log.PanicIf(err)

	expectedValue = constructLiItem("tag")
	expected = []interface{}{expectedValue}

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

func TestXmpPropertyIndex_addValue_One_OneLevel(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xn := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_name",
	}

	xpi := newXmpPropertyIndex(xn)

	if len(xpi.leaves) != 0 {
		t.Fatalf("Expected no initial leaves.")
	}

	xpn := xmpregistry.XmpPropertyName{xn}
	someValue := 55

	err := xpi.addValue(xpn, someValue)
	log.PanicIf(err)

	if len(xpi.leaves) != 1 {
		t.Fatalf("Expected one leave.")
	}

	expected := map[string][]interface{}{
		"[?]test_name": {someValue},
	}

	if reflect.DeepEqual(xpi.leaves, expected) != true {
		t.Fatalf("Index not correct.")
	}
}

func TestXmpPropertyIndex_addValue_One_MultipleLevel(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	rootXn := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "root_name",
	}

	xpi := newXmpPropertyIndex(rootXn)

	if len(xpi.subindices) != 0 {
		t.Fatalf("Expected zero subindices.")
	} else if len(xpi.leaves) != 0 {
		t.Fatalf("Expected no initial leaves.")
	}

	xn1 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_name1",
	}

	xn2 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_name2",
	}

	xpn := xmpregistry.XmpPropertyName{xn1, xn2}
	someValue := 55

	err := xpi.addValue(xpn, someValue)
	log.PanicIf(err)

	if len(xpi.subindices) != 1 {
		t.Fatalf("Expected one subindex.")
	}

	if len(xpi.leaves) != 0 {
		t.Fatalf("Expected zero leaves.")
	}

	// Check subindex.

	subindex := xpi.subindices["[?]test_name1"]

	if subindex == nil {
		t.Fatalf("Subindex was not found.")
	}

	if len(subindex.subindices) != 0 {
		t.Fatalf("Expected zero subindices.")
	}

	if len(subindex.leaves) != 1 {
		t.Fatalf("Expected one leave.")
	}

	expected := map[string][]interface{}{
		"[?]test_name2": {someValue},
	}

	if reflect.DeepEqual(subindex.leaves, expected) != true {
		t.Fatalf("Index not correct.")
	}
}

func TestXmpPropertyIndex_addValue_Multiple_OneLevel(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xn := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_name",
	}

	xpi := newXmpPropertyIndex(xn)

	if len(xpi.leaves) != 0 {
		t.Fatalf("Expected no initial leaves.")
	}

	// 1

	someValue1 := 55

	xn1 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node1",
	}

	xpn1 := xmpregistry.XmpPropertyName{xn1}

	err := xpi.addValue(xpn1, someValue1)
	log.PanicIf(err)

	// 2

	someValue2 := 65

	xn2 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node2",
	}

	xpn2 := xmpregistry.XmpPropertyName{xn2}

	err = xpi.addValue(xpn2, someValue2)
	log.PanicIf(err)

	// 3

	someValue3 := 75

	xn3 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node3",
	}

	xpn3 := xmpregistry.XmpPropertyName{xn3}

	err = xpi.addValue(xpn3, someValue3)
	log.PanicIf(err)

	expected := map[string][]interface{}{
		"[?]node1": {someValue1},
		"[?]node2": {someValue2},
		"[?]node3": {someValue3},
	}

	if reflect.DeepEqual(xpi.leaves, expected) != true {
		t.Fatalf("Index not correct.")
	}
}

func TestXmpPropertyIndex_addValue_Multiple_MultipleLevels_RandomAccess(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xn := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_name",
	}

	xpi := newXmpPropertyIndex(xn)

	if len(xpi.leaves) != 0 {
		t.Fatalf("Expected no initial leaves.")
	}

	// 1

	someValue11 := 55

	xn11 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node11",
	}

	xpn11 := xmpregistry.XmpPropertyName{xn11}

	err := xpi.addValue(xpn11, someValue11)
	log.PanicIf(err)

	someValue12 := 65

	xn12 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node12",
	}

	xpn12 := xmpregistry.XmpPropertyName{xn12}

	err = xpi.addValue(xpn12, someValue12)
	log.PanicIf(err)

	if xpi.Count() != 2 {
		t.Fatalf("Index count is not valid: (%d)", xpi.Count())
	}

	// 2

	someValue21 := 85

	xn21 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node21",
	}

	xpn21 := xmpregistry.XmpPropertyName{xn11, xn21}

	err = xpi.addValue(xpn21, someValue21)
	log.PanicIf(err)

	someValue22 := 95

	xn22 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node22",
	}

	xpn22 := xmpregistry.XmpPropertyName{xn12, xn22}

	err = xpi.addValue(xpn22, someValue22)
	log.PanicIf(err)

	if xpi.Count() != 4 {
		t.Fatalf("Index count is not valid: (%d)", xpi.Count())
	}

	// 3

	someValue31 := 115

	xn31 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node31",
	}

	xpn31 := xmpregistry.XmpPropertyName{xn11, xn21, xn31}

	err = xpi.addValue(xpn31, someValue31)
	log.PanicIf(err)

	someValue32 := 125

	xn32 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "node32",
	}

	xpn32 := xmpregistry.XmpPropertyName{xn12, xn22, xn32}

	err = xpi.addValue(xpn32, someValue32)
	log.PanicIf(err)

	if xpi.Count() != 6 {
		t.Fatalf("Index count is not valid: (%d)", xpi.Count())
	}

	// Check

	results, err := xpi.Get([]string{"[?]node12", "[?]node22", "[?]node32"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{125}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}

	results, err = xpi.Get([]string{"[?]node11", "[?]node21", "[?]node31"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{115}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}

	results, err = xpi.Get([]string{"[?]node12", "[?]node22"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{95}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}

	results, err = xpi.Get([]string{"[?]node11", "[?]node21"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{85}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}

	results, err = xpi.Get([]string{"[?]node12"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{65}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}

	results, err = xpi.Get([]string{"[?]node11"})
	log.PanicIf(err)

	if reflect.DeepEqual(results, []interface{}{55}) != true {
		t.Fatalf("query 1 not correct: %v", results)
	}
}

func TestComplexLeafNode_Get_Hit(t *testing.T) {
	xn := xml.Name{
		Space: "space/uri",
		Local: "test_name",
	}

	cln := ComplexLeafNode{
		xn: 55,
	}

	value, found := cln.Get("space/uri", "test_name")
	if found != true {
		t.Fatalf("Not found.")
	}

	if value != 55 {
		t.Fatalf("Value not correct: (%d)", value)
	}
}

func TestComplexLeafNode_Get_Miss(t *testing.T) {
	xn := xml.Name{
		Space: "space/uri",
		Local: "test_name",
	}

	cln := ComplexLeafNode{
		xn: 55,
	}

	_, found := cln.Get("space/uri", "unregistered_key")
	if found != false {
		t.Fatalf("Expected not found.")
	}
}

func TestXmpPropertyIndex_addArrayValue(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xnRoot := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "root_node",
	}

	xpi := newXmpPropertyIndex(xnRoot)

	oreaft := xmptype.OrderedResourceEventArrayFieldType{}

	xn1 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_node1",
	}

	xpn1 := xmpregistry.XmpPropertyName{xnRoot, xn1}

	items := []interface{}{
		11,
		22,
		33,
	}

	avOriginal := oreaft.New(xpn1, items)

	err := xpi.addArrayValue(xpn1, avOriginal)
	log.PanicIf(err)

	// Check.

	results, err := xpi.Get([]string{"[?]root_node", "[?]test_node1"})
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("Did not find exactly one result: (%d)", len(results))
	}

	avRecovered := results[0].(xmptype.ArrayValue)

	if reflect.DeepEqual(avRecovered, avOriginal) != true {
		t.Fatalf("Recovered array not correct.")
	}
}

func TestXmpPropertyIndex_addScalarValue(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xnRoot := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "root_node",
	}

	xpi := newXmpPropertyIndex(xnRoot)

	xn1 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_node1",
	}

	xpn1 := xmpregistry.XmpPropertyName{xnRoot, xn1}

	err := xpi.addScalarValue(xpn1, 55)
	log.PanicIf(err)

	// Check.

	results, err := xpi.Get([]string{"[?]root_node", "[?]test_node1"})
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("Did not find exactly one result: (%d)", len(results))
	}

	sln := ScalarLeafNode{
		Name:        xml.Name(xn1),
		ParsedValue: 55,
	}

	if reflect.DeepEqual(results[0], sln) != true {
		t.Fatalf("Recovered scalar not correct.")
	}
}

func TestXmpPropertyIndex_addComplexValue(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xnRoot := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "root_node",
	}

	xpi := newXmpPropertyIndex(xnRoot)

	xn1 := xmpregistry.XmlName{
		Space: "space/uri",
		Local: "test_node1",
	}

	xpn1 := xmpregistry.XmpPropertyName{xnRoot, xn1}

	xn2 := xml.Name{
		Space: "space/uri",
		Local: "test_value_node1",
	}

	attributes := map[xml.Name]interface{}{
		xn2: 55,
	}

	err := xpi.addComplexValue(xpn1, attributes)
	log.PanicIf(err)

	// Check.

	results, err := xpi.Get([]string{"[?]root_node", "[?]test_node1"})
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("Did not find exactly one result: (%d)", len(results))
	}

	cvn := results[0].(ComplexLeafNode)

	if reflect.DeepEqual(cvn, ComplexLeafNode(attributes)) != true {
		t.Fatalf("Recovered complex not correct.")
	}
}
