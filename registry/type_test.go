package xmpregistry

import (
	"reflect"
	"testing"

	"encoding/xml"
)

func TestClearCachedPrefixes(t *testing.T) {
	originalCachedPrefixes := cachedPrefixes
	cachedPrefixes = make(map[string]string)

	defer func() {
		cachedPrefixes = originalCachedPrefixes
	}()

	cachedPrefixes["aa"] = "bb"

	ClearCachedPrefixes()

	if len(cachedPrefixes) != 0 {
		t.Fatalf("ClearCachedPrefixes did not clear prefix cache.")
	}
}

func TestXmlName_Prefix_Known(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	namespaceUri := "http://ns.adobe.com/pdf/1.3/"

	namespace := Namespace{
		Uri:             namespaceUri,
		PreferredPrefix: "pdf",
	}

	Register(namespace)

	// Test.

	name := XmlName{
		Space: namespaceUri,
		Local: "bb",
	}

	if name.Prefix() != "pdf" {
		t.Fatalf("Prefix not correct: [%s]", name.Prefix())
	}

	if len(cachedPrefixes) != 1 {
		t.Fatalf("Expected prefix to be cached.")
	}

	cachedPrefix := cachedPrefixes[namespaceUri]

	if cachedPrefix != "pdf" {
		t.Fatalf("Cached prefix is not correct.")
	}
}

func TestXmlName_String_Known(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	namespaceUri := "http://ns.adobe.com/pdf/1.3/"

	namespace := Namespace{
		Uri:             namespaceUri,
		PreferredPrefix: "pdf",
	}

	Register(namespace)

	// Test.

	name := XmlName{
		Space: namespaceUri,
		Local: "bb",
	}

	if name.String() != "[pdf]bb" {
		t.Fatalf("String not correct: [%s]", name.String())
	}

	if len(cachedPrefixes) != 1 {
		t.Fatalf("Expected prefix to be cached.")
	}

	cachedPrefix := cachedPrefixes[namespaceUri]

	if cachedPrefix != "pdf" {
		t.Fatalf("Cached prefix is not correct.")
	}
}

func TestXmlName_Prefix_NotKnown(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	// Test.

	namespaceUri := "http://ns.adobe.com/pdf/1.3/"

	name := XmlName{
		Space: namespaceUri,
		Local: "bb",
	}

	if name.Prefix() != "?" {
		t.Fatalf("Prefix not correct: [%s]", name.Prefix())
	}

	if len(cachedPrefixes) != 1 {
		t.Fatalf("Expected prefix to be cached.")
	}

	cachedPrefix := cachedPrefixes[namespaceUri]

	if cachedPrefix != "?" {
		t.Fatalf("Cached prefix is not correct.")
	}
}

func TestXmlName_String_NotKnown(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	// Test.

	name := XmlName{
		Space: "aa",
		Local: "bb",
	}

	if name.String() != "[?]bb" {
		t.Fatalf("String not correct: [%s]", name.String())
	}
}

func TestXmpPropertyName_Parts(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	namespaceUri1 := "adobe:ns:meta/"

	namespace1 := Namespace{
		Uri:             namespaceUri1,
		PreferredPrefix: "x",
	}

	Register(namespace1)

	namespaceUri2 := "http://ns.adobe.com/xap/1.0/sType/Version#"

	namespace2 := Namespace{
		Uri:             namespaceUri2,
		PreferredPrefix: "stVer",
	}

	Register(namespace2)

	// Test.

	name := XmpPropertyName{
		{namespaceUri1, "aa"},
		{namespaceUri2, "bb"},
	}

	actual := name.Parts()

	expected := []string{
		"[x]aa",
		"[stVer]bb",
	}

	if reflect.DeepEqual(actual, expected) != true {
		t.Fatalf("Parts not correct: %v != %v\n", actual, expected)
	}
}

func TestXmpPropertyName_String(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	namespaceUri1 := "adobe:ns:meta/"

	namespace1 := Namespace{
		Uri:             namespaceUri1,
		PreferredPrefix: "x",
	}

	Register(namespace1)

	namespaceUri2 := "http://ns.adobe.com/xap/1.0/sType/Version#"

	namespace2 := Namespace{
		Uri:             namespaceUri2,
		PreferredPrefix: "stVer",
	}

	Register(namespace2)

	// Test.

	name := XmpPropertyName{
		{namespaceUri1, "aa"},
		{namespaceUri2, "bb"},
	}

	namePhrase := name.String()

	if namePhrase != "[x]aa.[stVer]bb" {
		t.Fatalf("String not correct: %s\n", namePhrase)
	}
}

func TestInlineAttributes(t *testing.T) {
	// Stage.

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	ClearCachedPrefixes()

	namespaceUri1 := "adobe:ns:meta/"

	namespace1 := Namespace{
		Uri:             namespaceUri1,
		PreferredPrefix: "x",
	}

	Register(namespace1)

	namespaceUri2 := "http://ns.adobe.com/xap/1.0/sType/Version#"

	namespace2 := Namespace{
		Uri:             namespaceUri2,
		PreferredPrefix: "stVer",
	}

	Register(namespace2)

	attributes := map[xml.Name]interface{}{
		xml.Name{Space: namespaceUri1, Local: "aa"}: "value1",
		xml.Name{Space: namespaceUri2, Local: "bb"}: "value2",
	}

	phrase := InlineAttributes(attributes)

	if phrase != "[stVer]bb=[value2] [x]aa=[value1]" {
		t.Fatalf("Inlined attributes are not correct: [%s]", phrase)
	}
}
