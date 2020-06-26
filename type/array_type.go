package xmptype

import (
	"fmt"
	"reflect"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

var (
	typeLogger = log.NewLogger("xmp.type")
)

const (
	// rdfUri is the URI for the "rdf" namespace. We can't use the same value
	// from xmpnamespace because xmptype can't import from it.
	rdfUri = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

var (
	rdfSeqTag = xml.Name{
		Space: rdfUri,
		Local: "Seq",
	}

	rdfBagTag = xml.Name{
		Space: rdfUri,
		Local: "Bag",
	}

	rdfAltTag = xml.Name{
		Space: rdfUri,
		Local: "Alt",
	}

	rdfLiTag = xml.Name{
		Space: rdfUri,
		Local: "li",
	}
)

type ArrayItem struct {
	Name       xml.Name
	Attributes map[xml.Name]interface{}
	CharData   string
}

func (ai ArrayItem) InlineAttributes() string {
	return xmpregistry.InlineAttributes(ai.Attributes)
}

type ArrayValue interface {
	FullName() xmpregistry.XmpPropertyName
	Count() int
}

// ArrayStringValueLister is any array type that has an Items() method that
// returns a string slice.
type ArrayStringValueLister interface {
	Items() (items []string, err error)
}

type ArrayFieldType interface {
	New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue
}

type baseArrayValue struct {
	fullName  xmpregistry.XmpPropertyName
	collected []interface{}
}

func newBaseArrayValue(fullName xmpregistry.XmpPropertyName, collected []interface{}) baseArrayValue {
	return baseArrayValue{
		fullName:  fullName,
		collected: collected,
	}
}

func (bav baseArrayValue) elementTagName(elements []interface{}, i int) (name xml.Name, isTag bool) {
	item := elements[i]
	if se, ok := item.(xml.StartElement); ok == true {
		return se.Name, true
	} else if se, ok := item.(xml.EndElement); ok == true {
		return se.Name, true
	}

	return name, false
}

// validateAnchorElements asserts that the list of elements starts and ends with
// the given tag.
func (bav baseArrayValue) validateAnchorElements(elements []interface{}, name xml.Name) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	elementCount := len(elements)

	if elementCount < 2 {
		log.Panicf(
			"expected at least two items for anchor-tag check: [%s]",
			bav.FullName())
	}

	firstElementTagName, firstElementIsTag := bav.elementTagName(elements, 0)

	if firstElementIsTag == false {
		log.Panicf(
			"expected first element in [%s] array to be a tag",
			bav.FullName(),
		)
	}

	if firstElementTagName != name {
		log.Panicf(
			"expected first element in [%s] array to be a [%s] tag: [%s]",
			bav.FullName(), xmpregistry.XmlName(name), xmpregistry.XmlName(firstElementTagName))
	}

	lastElementTagName, lastElementIsTag := bav.elementTagName(elements, elementCount-1)

	if lastElementIsTag == false {
		log.Panicf(
			"expected last element in [%s] array to be a tag",
			bav.FullName(),
		)
	}

	if lastElementTagName != name {
		log.Panicf(
			"expected last element in [%s] array to be a [%s] tag: [%s]",
			bav.FullName(), xmpregistry.XmlName(name), xmpregistry.XmlName(lastElementTagName))
	}

	return nil
}

func (bav baseArrayValue) FullName() xmpregistry.XmpPropertyName {
	return bav.fullName
}

func (bav baseArrayValue) Count() int {
	return len(bav.collected)
}

// innerItems will extract the attributes and char-data from all elements except
// the first and last. It expects these to all be "li" tags with an optional
// char-data value between them.
func (bav baseArrayValue) innerItems(hasSandwichedCharData bool) (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	elementCount := bav.Count()
	items = make([]ArrayItem, 0, elementCount-2)

	if hasSandwichedCharData == true {
		for i := 1; i < elementCount-1; {
			// Check for the first and the third elements to be "li" tags.

			subslice := bav.collected[i : i+3]

			err = bav.validateAnchorElements(subslice, rdfLiTag)
			log.PanicIf(err)

			se := subslice[0].(xml.StartElement)

			attributes, err := ParseAttributes(se)
			log.PanicIf(err)

			// Extract the character-data between the tags.

			charDataRaw := subslice[1]

			s, ok := charDataRaw.(string)
			if ok == false {
				log.Panicf(
					"expected element between 'li' tags in unordered-array to be char-data: [%s] [%s]",
					bav.FullName(), reflect.TypeOf(s))
			}

			ci := ArrayItem{
				Name:       se.Name,
				Attributes: attributes,
				CharData:   s,
			}

			items = append(items, ci)

			i += len(subslice)
		}
	} else {
		for i := 1; i < elementCount-1; {
			subslice := bav.collected[i : i+2]

			err = bav.validateAnchorElements(subslice, rdfLiTag)
			log.PanicIf(err)

			se := subslice[0].(xml.StartElement)

			attributes, err := ParseAttributes(se)
			log.PanicIf(err)

			ci := ArrayItem{
				Name:       se.Name,
				Attributes: attributes,
			}

			items = append(items, ci)

			i += len(subslice)
		}
	}

	return items, nil
}

// Ordered array semantics

// TODO(dustin): Ordered array yet-to-implement: CuePointParam, Marker, ResourceEvent, Version, Colorant, Marker, Layer, "point" (?)

type OrderedArrayValue struct {
	baseArrayValue
}

func newOrderedArrayValue(bav baseArrayValue) OrderedArrayValue {
	return OrderedArrayValue{
		baseArrayValue: bav,
	}
}

func (oav OrderedArrayValue) String() string {
	return fmt.Sprintf("OrderedArray<COUNT=(%d)>", oav.Count())
}

func (oav OrderedArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = oav.validateAnchorElements(oav.baseArrayValue.collected, rdfSeqTag)
	log.PanicIf(err)

	items, err = oav.innerItems(false)
	log.PanicIf(err)

	return items, nil
}

type OrderedArrayFieldType struct {
}

func (oat OrderedArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return OrderedArrayValue{
		baseArrayValue: bav,
	}
}

type OrderedTextArrayFieldType struct {
	OrderedArrayFieldType
}

type OrderedUriArrayFieldType struct {
	OrderedArrayFieldType
}

type OrderedResourceEventArrayValue struct {
	OrderedArrayValue
}

// Items this is a wrapper that returns a simple list of strings from inner
// underlying array-items, thereby satisfying the ArrayStringValueLister
// interface. In the case of these, we return a stringification of the
// attributes.
func (oreav OrderedResourceEventArrayValue) Items() (items []string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	innerItems, err := oreav.OrderedArrayValue.Items()
	log.PanicIf(err)

	items = make([]string, len(innerItems))
	for i, ai := range innerItems {
		items[i] = ai.InlineAttributes()
	}

	return items, nil
}

type OrderedResourceEventArrayFieldType struct {
}

func (oreat OrderedResourceEventArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	oav := newOrderedArrayValue(bav)

	return OrderedResourceEventArrayValue{
		OrderedArrayValue: oav,
	}
}

// Unordered array semantics

// TODO(dustin): Unordered array yet-to-implement: XPath, ResourceRef, "struct" (?), Job, Font, Media, Track

type UnorderedArrayValue struct {
	baseArrayValue
}

func newUnorderedArrayValue(bav baseArrayValue) UnorderedArrayValue {
	return UnorderedArrayValue{
		baseArrayValue: bav,
	}
}

func (uav UnorderedArrayValue) String() string {
	return fmt.Sprintf("UnorderedArray<COUNT=(%d)>", uav.Count())
}

func (uav UnorderedArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = uav.validateAnchorElements(uav.baseArrayValue.collected, rdfBagTag)
	log.PanicIf(err)

	items, err = uav.innerItems(true)
	log.PanicIf(err)

	return items, nil
}

type UnorderedArrayFieldType struct {
}

func (uat UnorderedArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return newUnorderedArrayValue(bav)
}

type UnorderedTextArrayFieldType struct {
	UnorderedArrayFieldType
}

type UnorderedAncestorArrayValue struct {
	UnorderedArrayValue
}

// Items this is a wrapper that returns a simple list of strings from inner
// underlying array-items, thereby satisfying the ArrayStringValueLister
// interface.
func (uaav UnorderedAncestorArrayValue) Items() (items []string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	innerItems, err := uaav.UnorderedArrayValue.Items()
	log.PanicIf(err)

	items = make([]string, len(innerItems))
	for i, ai := range innerItems {
		items[i] = ai.CharData
	}

	return items, nil
}

type UnorderedAncestorArrayFieldType struct {
	UnorderedArrayFieldType
}

func (uaat UnorderedAncestorArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return UnorderedAncestorArrayValue{
		UnorderedArrayValue: newUnorderedArrayValue(bav),
	}
}

// Alternatives array semantics

type AlternativeArrayValue struct {
	baseArrayValue
}

func newAlternativeArrayValue(bav baseArrayValue) AlternativeArrayValue {
	return AlternativeArrayValue{
		baseArrayValue: bav,
	}
}

func (aav AlternativeArrayValue) String() string {
	return fmt.Sprintf("AlternativeArray<COUNT=(%d)>", aav.Count())
}

func (aav AlternativeArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = aav.validateAnchorElements(aav.baseArrayValue.collected, rdfAltTag)
	log.PanicIf(err)

	items, err = aav.innerItems(true)
	log.PanicIf(err)

	return items, nil
}

type AlternativeArrayFieldType struct {
}

func (aat AlternativeArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return AlternativeArrayValue{
		baseArrayValue: bav,
	}
}

type LanguageAlternativeArrayValue struct {
	AlternativeArrayValue
}

func (laav LanguageAlternativeArrayValue) Items() (items []string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	innerItems, err := laav.AlternativeArrayValue.Items()
	log.PanicIf(err)

	items = make([]string, len(innerItems))
	for i, ai := range innerItems {
		items[i] = fmt.Sprintf("{%s} [%s]", ai.InlineAttributes(), ai.CharData)
	}

	return items, nil
}

type LanguageAlternativeArrayFieldType struct {
}

func (laat LanguageAlternativeArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	aav := newAlternativeArrayValue(bav)

	return LanguageAlternativeArrayValue{
		AlternativeArrayValue: aav,
	}
}
