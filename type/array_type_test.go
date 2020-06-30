package xmptype

import (
	"fmt"
	"reflect"
	"testing"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

func TestArrayItem_InlineAttributes(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	attributes := map[xml.Name]interface{}{
		{Space: RdfUri, Local: "aa"}: "value1",
		{Space: xmpUri, Local: "bb"}: "value2",
	}

	ai := ArrayItem{
		Attributes: attributes,
	}

	phrase := ai.InlineAttributes()
	if phrase != "[rdf]aa=[value1] [xmp]bb=[value2]" {
		t.Fatalf("Inline attributes not correct: [%s]", phrase)
	}
}

func TestArrayItem_String(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	attributes := map[xml.Name]interface{}{
		{Space: RdfUri, Local: "aa"}: "value1",
		{Space: xmpUri, Local: "bb"}: "value2",
	}

	ai := ArrayItem{
		Name:       xml.Name{Space: RdfUri, Local: "li"},
		Attributes: attributes,
		CharData:   "test-char-data",
	}

	if ai.String() != "ArrayItem<NAME={[rdf]li} ATTR={[rdf]aa=[value1] [xmp]bb=[value2]} CHAR-DATA=[test-char-data]>" {
		t.Fatalf("Strign not correct: [%s]", ai.String())
	}
}

func TestElementTagName_NotTag(t *testing.T) {
	items := []interface{}{
		"value1",
		"value2",
	}

	_, isTag, isOpenTag := elementTagName(items, 1)

	if isTag != false {
		t.Fatalf("Expected a non-tag")
	} else if isOpenTag != false {
		t.Fatalf("Expected a non-tag to return false for an open-tag.")
	}
}

func TestElementTagName_Tag_Open(t *testing.T) {
	expected := xml.Name{Space: RdfUri, Local: "aa"}

	items := []interface{}{
		"value1",
		xml.StartElement{Name: expected},
	}

	actual, isTag, isOpenTag := elementTagName(items, 1)

	if isTag != true {
		t.Fatalf("Expected a tag.")
	} else if actual != expected {
		t.Fatalf("Name not expected.")
	} else if isOpenTag != true {
		t.Fatalf("Expected to be open-tag.")
	}
}

func TestElementTagName_Tag_Close(t *testing.T) {
	expected := xml.Name{Space: RdfUri, Local: "aa"}

	items := []interface{}{
		"value1",
		xml.EndElement{Name: expected},
	}

	actual, isTag, isOpenTag := elementTagName(items, 1)

	if isTag != true {
		t.Fatalf("Expected a tag.")
	} else if actual != expected {
		t.Fatalf("Name not expected.")
	} else if isOpenTag != false {
		t.Fatalf("Expected to be close-tag.")
	}
}

func TestValidateAnchorElements_Hit(t *testing.T) {
	name := xml.Name{Space: RdfUri, Local: "aa"}

	items := []interface{}{
		xml.StartElement{Name: name},
		"some_value",
		xml.EndElement{Name: name},
	}

	err := validateAnchorElements(items, name)
	log.PanicIf(err)
}

func TestValidateAnchorElements_Miss_Balanced(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	haveName := xml.Name{Space: RdfUri, Local: "aa"}
	needName := xml.Name{Space: xmpUri, Local: "bb"}

	items := []interface{}{
		xml.StartElement{Name: haveName},
		"some_value",
		xml.EndElement{Name: haveName},
	}

	err := validateAnchorElements(items, needName)
	if err == nil {
		t.Fatalf("Expected error for incorrect anchor tag names.")
	} else if err.Error() != "expected first element in array to be a [[?]bb] tag: [[?]aa]" {
		log.Panic(err)
	}
}

func TestValidateAnchorElements_Miss_Unbalanced(t *testing.T) {
	name := xml.Name{Space: RdfUri, Local: "aa"}

	items := []interface{}{
		xml.StartElement{Name: name},
		"some_value",
	}

	err := validateAnchorElements(items, name)
	if err == nil {
		t.Fatalf("Expected error for unbalanced anchor tags.")
	} else if err.Error() != "expected last element in array to be a tag" {
		log.Panic(err)
	}
}

func TestNewBaseArrayValue(t *testing.T) {
	items := []interface{}{
		"value1",
		"value2",
	}

	bav := newBaseArrayValue(testPropertyName, items)

	if reflect.DeepEqual(bav.fullName, testPropertyName) != true {
		t.Fatalf("Full-name not correct.")
	} else if reflect.DeepEqual(bav.collected, items) != true {
		t.Fatalf("Collected items not correct.")
	}
}

func TestBaseArrayValue_FullName(t *testing.T) {
	items := []interface{}{
		"value1",
		"value2",
	}

	bav := newBaseArrayValue(testPropertyName, items)

	if reflect.DeepEqual(bav.FullName(), testPropertyName) != true {
		t.Fatalf("FullName() not correct.")
	}
}

func TestBaseArrayValue_Count(t *testing.T) {
	items := []interface{}{
		"value1",
		"value2",
	}

	bav := newBaseArrayValue(testPropertyName, items)

	if bav.Count() != 2 {
		t.Fatalf("Count() not correct.")
	}
}

func TestBaseArrayValue_constructArrayItem_WithChardata(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithChardata()

	actual, err := bav.constructArrayItem(bav.collected[1:4])
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	expected := ArrayItem{
		Name:       itemName,
		Attributes: extractedAttributes1,
		CharData:   "value2",
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("ArrayItem not correct.")
	}
}

func TestBaseArrayValue_constructArrayItem_WithoutChardata(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithoutChardata()

	actual, err := bav.constructArrayItem(bav.collected[1:3])
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	expected := ArrayItem{
		Name:       itemName,
		Attributes: extractedAttributes1,
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("ArrayItem not correct.")
	}
}

func TestBaseArrayValue_innerItems_WithChardata(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithChardata()

	actualItems, err := bav.innerItems(true)
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("innerItems() not correct.")
	}
}

func TestBaseArrayValue_innerItems_WithoutChardata(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithoutChardata()

	actualItems, err := bav.innerItems(false)
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("innerItems() not correct.")
	}
}

// Ordered array

func TestNewOrderedArrayValue(t *testing.T) {
	// Not much that we can test here.
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, nil)
	newOrderedArrayValue(bav)
}

func TestOrderedArrayValue_String(t *testing.T) {
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, []interface{}{"aa", "bb", "cc"})
	oav := newOrderedArrayValue(bav)

	if oav.String() != "OrderedArray<COUNT=(3)>" {
		t.Fatalf("String not correct: [%s]", oav.String())
	}
}

func TestOrderedArrayValue_Items(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithoutChardata()
	oav := newOrderedArrayValue(bav)

	actualItems, err := oav.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

func TestOrderedArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	oaft := OrderedArrayFieldType{}

	items := getTestSequenceItemsWithoutChardata()
	av := oaft.New(testPropertyName, items)

	actualItems, err := av.(OrderedArrayValue).Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

// OrderedResourceEvent

func TestOrderedResourceEventArrayValue_Items(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestSequenceBaseArrayValueWithoutChardata()
	oav := newOrderedArrayValue(bav)

	oreav := OrderedResourceEventArrayValue{
		OrderedArrayValue: oav,
	}

	actual, err := oreav.Items()
	log.PanicIf(err)

	expected := []string{
		"[rdf]item1=[test_value_1] [rdf]item2=[test_value_2]",
		"[rdf]item1=[test_value_4] [rdf]item2=[test_value_3]",
	}

	if reflect.DeepEqual(actual, expected) != true {
		fmt.Printf("Actual:\n")
		fmt.Printf("\n")

		for _, item := range actual {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		fmt.Printf("Expected:\n")
		fmt.Printf("\n")

		for _, item := range expected {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		t.Fatalf("Items not correct.")
	}
}

func TestOrderedResourceEventArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	items := getTestSequenceItemsWithoutChardata()

	oreaft := OrderedResourceEventArrayFieldType{}
	av := oreaft.New(testPropertyName, items)

	oreav := av.(OrderedResourceEventArrayValue)

	// Test base Items() implementation.

	actualItems, err := oreav.OrderedArrayValue.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		for _, ai := range actualItems {
			fmt.Printf("ACTUAL: %s\n", ai)
		}

		for _, ai := range expectedItems {
			fmt.Printf("EXPECTED: %s\n", ai)
		}

		t.Fatalf("Items not correct.")
	}

	// Test type-specific Items() implementation.

	actualStrings, err := oreav.Items()
	log.PanicIf(err)

	expectedStrings := []string{
		"[rdf]item1=[test_value_1] [rdf]item2=[test_value_2]",
		"[rdf]item1=[test_value_4] [rdf]item2=[test_value_3]",
	}

	if reflect.DeepEqual(actualStrings, expectedStrings) != true {
		t.Fatalf("Items not correct.")
	}
}

// Unordered array

func TestNewUnorderedArrayValue(t *testing.T) {
	// Not much that we can test here.
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, nil)
	newUnorderedArrayValue(bav)
}

func TestUnorderedArrayValue_String(t *testing.T) {
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, []interface{}{"aa", "bb", "cc"})
	uav := newUnorderedArrayValue(bav)

	if uav.String() != "UnorderedArray<COUNT=(3)>" {
		t.Fatalf("String not correct: [%s]", uav.String())
	}
}

func TestUnorderedArrayValue_Items(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestBagBaseArrayValueWithChardata()
	uav := newUnorderedArrayValue(bav)

	actualItems, err := uav.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

func TestUnorderedArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	uaft := UnorderedArrayFieldType{}

	items := getTestBagItemsWithChardata()
	av := uaft.New(testPropertyName, items)

	actualItems, err := av.(UnorderedArrayValue).Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

func TestUnorderedAncestorArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	uaaft := UnorderedAncestorArrayFieldType{}

	items := getTestBagItemsWithChardata()
	av := uaaft.New(testPropertyName, items)

	uaav := av.(UnorderedAncestorArrayValue)

	// Test base Items() implementation.

	actualItems, err := uaav.UnorderedArrayValue.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		for _, ai := range actualItems {
			fmt.Printf("ACTUAL: %s\n", ai)
		}

		for _, ai := range expectedItems {
			fmt.Printf("EXPECTED: %s\n", ai)
		}

		t.Fatalf("Items() not correct.")
	}

	// Test type-specific Items() implementation.

	actualStrings, err := uaav.Items()
	log.PanicIf(err)

	expectedStrings := []string{
		"value2",
		"value1",
	}

	if reflect.DeepEqual(actualStrings, expectedStrings) != true {
		for _, s := range actualStrings {
			fmt.Printf("ACTUAL> %s\n", s)
		}

		for _, s := range expectedStrings {
			fmt.Printf("EXPECTED> %s\n", s)
		}

		t.Fatalf("Items not correct.")
	}
}

// Alternative array

func TestNewAlternativeArrayValue(t *testing.T) {
	// Not much that we can test here.
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, nil)
	newAlternativeArrayValue(bav)
}

func TestAlternativeArrayValue_String(t *testing.T) {
	xpn := xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "some_node"},
	}

	bav := newBaseArrayValue(xpn, []interface{}{"aa", "bb", "cc"})
	aav := newAlternativeArrayValue(bav)

	if aav.String() != "AlternativeArray<COUNT=(3)>" {
		t.Fatalf("String not correct: [%s]", aav.String())
	}
}

func TestAlternativeArrayValue_Items(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestAltBaseArrayValueWithChardata()
	aav := newAlternativeArrayValue(bav)

	actualItems, err := aav.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

func TestAlternativeArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	oaft := AlternativeArrayFieldType{}

	items := getTestAltItemsWithChardata()
	av := oaft.New(testPropertyName, items)

	aav := av.(AlternativeArrayValue)

	actualItems, err := aav.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		t.Fatalf("Items() not correct.")
	}
}

func TestLanguageAlternativeArrayValue_Items(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	bav := getTestAltBaseArrayValueWithChardata()
	aav := newAlternativeArrayValue(bav)

	laav := LanguageAlternativeArrayValue{
		AlternativeArrayValue: aav,
	}

	actual, err := laav.Items()
	log.PanicIf(err)

	expected := []string{
		"{[rdf]item1=[test_value_1] [rdf]item2=[test_value_2]} [value2]",
		"{[rdf]item1=[test_value_4] [rdf]item2=[test_value_3]} [value1]",
	}

	if reflect.DeepEqual(actual, expected) != true {
		fmt.Printf("Actual:\n")
		fmt.Printf("\n")

		for _, item := range actual {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		fmt.Printf("Expected:\n")
		fmt.Printf("\n")

		for _, item := range expected {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		t.Fatalf("Items not correct.")
	}
}

func TestLanguageAlternativeArrayFieldType_New(t *testing.T) {
	defer xmpregistry.Clear()
	registerTestNamespaces()

	laaft := LanguageAlternativeArrayFieldType{}

	items := getTestAltItemsWithChardata()
	av := laaft.New(testPropertyName, items)

	laav := av.(LanguageAlternativeArrayValue)

	// Test base Items() implementation.

	actualItems, err := laav.AlternativeArrayValue.Items()
	log.PanicIf(err)

	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	extractedAttributes1 := map[xml.Name]interface{}{
		attribute1Name: "test_value_1",
		attribute2Name: "test_value_2",
	}

	extractedAttributes2 := map[xml.Name]interface{}{
		attribute1Name: "test_value_4",
		attribute2Name: "test_value_3",
	}

	expectedItems := []ArrayItem{
		{
			Name:       itemName,
			Attributes: extractedAttributes1,
			CharData:   "value2",
		},
		{
			Name:       itemName,
			Attributes: extractedAttributes2,
			CharData:   "value1",
		},
	}

	if reflect.DeepEqual(actualItems, expectedItems) != true {
		for _, ai := range actualItems {
			fmt.Printf("ACTUAL: %s\n", ai)
		}

		for _, ai := range expectedItems {
			fmt.Printf("EXPECTED: %s\n", ai)
		}

		t.Fatalf("Items not correct.")
	}

	// Test type-specific Items() implementation.
	actual, err := laav.Items()
	log.PanicIf(err)

	expected := []string{
		"{[rdf]item1=[test_value_1] [rdf]item2=[test_value_2]} [value2]",
		"{[rdf]item1=[test_value_4] [rdf]item2=[test_value_3]} [value1]",
	}

	if reflect.DeepEqual(actual, expected) != true {
		fmt.Printf("Actual:\n")
		fmt.Printf("\n")

		for _, item := range actual {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		fmt.Printf("Expected:\n")
		fmt.Printf("\n")

		for _, item := range expected {
			fmt.Printf("%s\n", item)
		}

		fmt.Printf("\n")

		t.Fatalf("Items not correct.")
	}
}
