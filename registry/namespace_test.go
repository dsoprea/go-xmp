package xmpregistry

import (
	"reflect"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestNamespace_String(t *testing.T) {
	namespace := Namespace{
		Uri:             "http://some/uri",
		PreferredPrefix: "someprefix",
	}

	if namespace.String() != "Namespace<URI=[http://some/uri] PREFIX=[someprefix]>" {
		t.Fatalf("String not expected: [%s]", namespace.String())
	}
}

func TestRegister_Hit(t *testing.T) {
	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	uri := "http://some/uri/TestRegister"

	namespace := Namespace{
		Uri:             uri,
		PreferredPrefix: "TestRegister",
	}

	Register(namespace)

	if len(namespaces) != 1 {
		t.Fatalf("Registrations count not correct: (%d)", len(namespaces))
	}

	recalled, err := Get(uri)
	log.PanicIf(err)

	if reflect.DeepEqual(recalled, namespace) != true {
		t.Fatalf("Recalled namespace not correct: %s", recalled)
	}
}

func TestClear(t *testing.T) {
	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	uri := "http://some/uri/TestRegister"

	namespace := Namespace{
		Uri:             uri,
		PreferredPrefix: "TestRegister",
	}

	Register(namespace)

	if len(namespaces) != 1 {
		t.Fatalf("Registrations count not correct: (%d)", len(namespaces))
	}

	Clear()

	if len(namespaces) != 0 {
		t.Fatalf("Registrations not cleared.")
	}
}

func TestRegister_Get_Hit(t *testing.T) {
	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	uri := "http://some/uri/TestRegister"

	namespace := Namespace{
		Uri:             uri,
		PreferredPrefix: "TestRegister",
	}

	Register(namespace)

	if len(namespaces) != 1 {
		t.Fatalf("Registrations count not correct: (%d)", len(namespaces))
	}

	recalled, err := Get(uri)
	log.PanicIf(err)

	if reflect.DeepEqual(recalled, namespace) != true {
		t.Fatalf("Recalled namespace not correct: %s", recalled)
	}
}

func TestRegister_Get_Miss(t *testing.T) {
	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	_, err := Get("unknown/uri")

	if err != ErrNamespaceNotFound {
		t.Fatalf("Expected namespace miss for unknown namespace-URI: [%v]", err)
	}
}

func TestRegister_MustGet_Hit(t *testing.T) {
	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	uri := "http://some/uri/TestRegister"

	namespace := Namespace{
		Uri:             uri,
		PreferredPrefix: "TestRegister",
	}

	Register(namespace)

	if len(namespaces) != 1 {
		t.Fatalf("Registrations count not correct: (%d)", len(namespaces))
	}

	recalled := MustGet(uri)

	if reflect.DeepEqual(recalled, namespace) != true {
		t.Fatalf("Recalled namespace not correct: %s", recalled)
	}
}

func TestRegister_MustGet_Miss(t *testing.T) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err := errRaw.(error)

			if err != ErrNamespaceNotFound {
				log.Panic(err)
			}
		} else {
			t.Fatalf("Expected failure for MustGet miss.")
		}
	}()

	originalNamespaces := namespaces
	namespaces = make(map[string]Namespace)

	defer func() {
		namespaces = originalNamespaces
	}()

	MustGet("unknown/uri")
}
