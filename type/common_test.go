package xmptype

import (
	"encoding/xml"

	"github.com/dsoprea/go-xmp/registry"
)

const (
	xmpUri = "http://ns.adobe.com/xap/1.0/"
)

var (
	testPropertyName = xmpregistry.XmpPropertyName{
		xmpregistry.XmlName{Space: RdfUri, Local: "aa"},
		xmpregistry.XmlName{Space: xmpUri, Local: "bb"},
	}
)

func registerTestNamespaces() {
	xmpregistry.Clear()

	namespace := xmpregistry.Namespace{
		Uri:             RdfUri,
		PreferredPrefix: "rdf",
		Fields: map[string]interface{}{
			"item1": TextFieldType{},
			"item2": TextFieldType{},
		},
	}

	xmpregistry.Register(namespace)

	namespace = xmpregistry.Namespace{
		Uri:             xmpUri,
		PreferredPrefix: "xmp",
	}

	xmpregistry.Register(namespace)
}

func getTestSequenceItemsWithChardata() []interface{} {
	arrayName := xml.Name{Space: RdfUri, Local: "Seq"}
	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	attributes1 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_1"},
		{Name: attribute2Name, Value: "test_value_2"},
	}

	// We're deliberately misordering these.
	attributes2 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_4"},
		{Name: attribute2Name, Value: "test_value_3"},
	}

	// We're deliberately misordering these.
	items := []interface{}{
		xml.StartElement{Name: arrayName},
		xml.StartElement{Name: itemName, Attr: attributes1},
		"value2",
		xml.EndElement{Name: itemName},
		xml.StartElement{Name: itemName, Attr: attributes2},
		"value1",
		xml.EndElement{Name: itemName},
		xml.EndElement{Name: arrayName},
	}

	return items
}

func getTestSequenceBaseArrayValueWithChardata() baseArrayValue {
	items := getTestSequenceItemsWithChardata()

	bav := newBaseArrayValue(testPropertyName, items)
	return bav
}

func getTestBagItemsWithChardata() []interface{} {
	arrayName := xml.Name{Space: RdfUri, Local: "Bag"}
	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	attributes1 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_1"},
		{Name: attribute2Name, Value: "test_value_2"},
	}

	// We're deliberately misordering these.
	attributes2 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_4"},
		{Name: attribute2Name, Value: "test_value_3"},
	}

	// We're deliberately misordering these.
	items := []interface{}{
		xml.StartElement{Name: arrayName},
		xml.StartElement{Name: itemName, Attr: attributes1},
		"value2",
		xml.EndElement{Name: itemName},
		xml.StartElement{Name: itemName, Attr: attributes2},
		"value1",
		xml.EndElement{Name: itemName},
		xml.EndElement{Name: arrayName},
	}

	return items
}

func getTestBagBaseArrayValueWithChardata() baseArrayValue {
	items := getTestBagItemsWithChardata()

	bav := newBaseArrayValue(testPropertyName, items)
	return bav
}

func getTestAltItemsWithChardata() []interface{} {
	arrayName := xml.Name{Space: RdfUri, Local: "Alt"}
	itemName := xml.Name{Space: RdfUri, Local: "li"}

	attribute1Name := xml.Name{Space: RdfUri, Local: "item1"}
	attribute2Name := xml.Name{Space: RdfUri, Local: "item2"}

	attributes1 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_1"},
		{Name: attribute2Name, Value: "test_value_2"},
	}

	// We're deliberately misordering these.
	attributes2 := []xml.Attr{
		{Name: attribute1Name, Value: "test_value_4"},
		{Name: attribute2Name, Value: "test_value_3"},
	}

	// We're deliberately misordering these.
	items := []interface{}{
		xml.StartElement{Name: arrayName},
		xml.StartElement{Name: itemName, Attr: attributes1},
		"value2",
		xml.EndElement{Name: itemName},
		xml.StartElement{Name: itemName, Attr: attributes2},
		"value1",
		xml.EndElement{Name: itemName},
		xml.EndElement{Name: arrayName},
	}

	return items
}

func getTestAltBaseArrayValueWithChardata() baseArrayValue {
	items := getTestAltItemsWithChardata()

	bav := newBaseArrayValue(testPropertyName, items)
	return bav
}
