package xmptype

import (
	"reflect"
	"testing"
	"time"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

const (
	RdfUri = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	XmpUri = "http://ns.adobe.com/xap/1.0/"
)

func TestParseValue_Good(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestNumber": IntegerFieldType{},
		},
	}

	parsedValue, err := ParseValue(namespace, "TestNumber", "123")
	log.PanicIf(err)

	if parsedValue.(int64) != int64(123) {
		t.Fatalf("Parse failed.")
	}
}

func TestParseValue_Bad(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestField": IntegerFieldType{},
		},
	}

	_, err := ParseValue(namespace, "TestField", "abc")
	if err != ErrValueNotValid {
		log.Panic(err)
	}
}

func TestParseValue_InvalidChild(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestField": IntegerFieldType{},
		},
	}

	_, err := ParseValue(namespace, "InvalidField", "abc")
	if err == nil {
		t.Fatalf("Expected error for invalid child.")
	} else if err != ErrChildFieldNotFound {
		log.Panic(err)
	}
}

func TestIsArrayType_Hit(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestField": OrderedArrayFieldType{},
		},
	}

	flag, err := IsArrayType(namespace, "TestField")
	log.PanicIf(err)

	if flag != true {
		t.Fatalf("Expected array-type.")
	}
}

func TestIsArrayType_Miss(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestField": IntegerFieldType{},
		},
	}

	flag, err := IsArrayType(namespace, "TestField")
	log.PanicIf(err)

	if flag != false {
		t.Fatalf("Expected non-array type.")
	}
}

func TestIsArrayType_InvalidChild(t *testing.T) {
	namespace := xmpregistry.Namespace{
		Uri: "some/uri",
		Fields: map[string]interface{}{
			"TestField": IntegerFieldType{},
		},
	}

	_, err := IsArrayType(namespace, "InvalidField")
	if err == nil {
		t.Fatalf("Expected error for invalid child.")
	} else if err != ErrChildFieldNotFound {
		log.Panic(err)
	}
}

func TestParseAttributes_Ok(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpNamespace := xmpregistry.Namespace{
		Uri:             XmpUri,
		PreferredPrefix: "xmp",
		Fields: map[string]interface{}{
			"BaseURL":      UrlFieldType{},
			"CreateDate":   DateFieldType{},
			"CreatorTool":  AgentNameFieldType{},
			"Identifier":   UnorderedTextArrayFieldType{},
			"Label":        TextFieldType{},
			"MetadataDate": DateFieldType{},
			"ModifyDate":   DateFieldType{},
			"Nickname":     TextFieldType{},
			"Rating":       RealFieldType{},
		},
	}

	xmpregistry.Register(xmpNamespace)

	labelName := xml.Name{Space: XmpUri, Local: "Label"}
	modifyDateName := xml.Name{Space: XmpUri, Local: "ModifyDate"}

	rawAttributes := []xml.Attr{
		{
			Name:  labelName,
			Value: "test_label_value",
		},
		{
			Name:  modifyDateName,
			Value: "2020-06-26",
		},
	}

	se := xml.StartElement{
		Attr: rawAttributes,
	}

	actual, err := ParseAttributes(se)
	log.PanicIf(err)

	expected := map[xml.Name]interface{}{
		labelName:      "test_label_value",
		modifyDateName: time.Date(2020, 6, 26, 0, 0, 0, 0, time.UTC),
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Attributes not parsed correctly.")
	}
}

func TestParseAttributes_SkipUnknownNamespaces(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	labelName := xml.Name{Space: XmpUri, Local: "Label"}
	modifyDateName := xml.Name{Space: XmpUri, Local: "ModifyDate"}

	rawAttributes := []xml.Attr{
		{
			Name:  labelName,
			Value: "test_label_value",
		},
		{
			Name:  modifyDateName,
			Value: "2020-06-26",
		},
	}

	se := xml.StartElement{
		Attr: rawAttributes,
	}

	actual, err := ParseAttributes(se)
	log.PanicIf(err)

	expected := map[xml.Name]interface{}{}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Attributes not parsed correctly.")
	}
}

func TestParseAttributes_SkipInvalidFields(t *testing.T) {
	xmpregistry.Clear()
	defer xmpregistry.Clear()

	xmpNamespace := xmpregistry.Namespace{
		Uri:             XmpUri,
		PreferredPrefix: "xmp",
		Fields: map[string]interface{}{
			"BaseURL":      UrlFieldType{},
			"CreateDate":   DateFieldType{},
			"CreatorTool":  AgentNameFieldType{},
			"Identifier":   UnorderedTextArrayFieldType{},
			"Label":        TextFieldType{},
			"MetadataDate": DateFieldType{},
			"Nickname":     TextFieldType{},
			"Rating":       RealFieldType{},
		},
	}

	xmpregistry.Register(xmpNamespace)

	labelName := xml.Name{Space: XmpUri, Local: "Label"}
	modifyDateName := xml.Name{Space: XmpUri, Local: "ModifyDate"}

	rawAttributes := []xml.Attr{
		{
			Name:  labelName,
			Value: "test_label_value",
		},
		{
			Name:  modifyDateName,
			Value: "2020-06-26",
		},
	}

	se := xml.StartElement{
		Attr: rawAttributes,
	}

	actual, err := ParseAttributes(se)
	log.PanicIf(err)

	expected := map[xml.Name]interface{}{
		labelName: "test_label_value",
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Attributes not parsed correctly.")
	}
}
