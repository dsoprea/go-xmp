package xmptype

import (
	"errors"
	"fmt"
	"reflect"

	"encoding/xml"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-xmp/registry"
)

var (
	typeLogger = log.NewLogger("xmp.type")
)

var (
	// ErrArrayItemsNotOverridden indicates that a particular type is not
	// correctly overridden.
	ErrArrayItemsNotOverridden = errors.New("array type method must be overridden")
)

const (
	// RdfUri is the URI for the "rdf" namespace. We can't use the same value
	// from xmpnamespace because xmptype can't import from it.
	RdfUri = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
)

var (
	rdfSeqTag = xml.Name{
		Space: RdfUri,
		Local: "Seq",
	}

	rdfBagTag = xml.Name{
		Space: RdfUri,
		Local: "Bag",
	}

	rdfAltTag = xml.Name{
		Space: RdfUri,
		Local: "Alt",
	}

	rdfLiTag = xml.Name{
		Space: RdfUri,
		Local: "li",
	}
)

// ArrayItem is the item type of the extracted array items.
type ArrayItem struct {
	// Name is the name of the item node.
	Name xml.Name

	// Attributes are the attributes, if any, of the array item.
	Attributes map[xml.Name]interface{}

	// CharData is the trimmed char-data found in the array item.
	CharData string
}

// String returns a string representation of the item.
func (ai ArrayItem) String() string {
	return fmt.Sprintf(
		"ArrayItem<NAME={%s} ATTR={%s} CHAR-DATA=[%s]>",
		xmpregistry.XmlName(ai.Name),
		ai.InlineAttributes(),
		ai.CharData)
}

// InlineAttributes returns an inline string representation of all attributes.
func (ai ArrayItem) InlineAttributes() string {
	return xmpregistry.InlineAttributes(ai.Attributes)
}

// ArrayValue is satisfied by all array value types.
type ArrayValue interface {
	// FullName returns the name of the array container.
	FullName() xmpregistry.XmpPropertyName

	// Count returns the number of items.
	Count() int
}

// ArrayStringValueLister is any array type that has an Items() method that
// returns a string slice.
type ArrayStringValueLister interface {
	// Items returns a slice of strings.
	Items() (items []string, err error)
}

// ArrayFieldType is satisfied by all array field-types..
type ArrayFieldType interface {
	// New returns a new value struct encapsulating the given arguments.
	New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue
}

// elementTagName returns the xml.Name for the ith element. isTag indicates
// whether that element is actually a tag.
func elementTagName(elements []interface{}, i int) (name xml.Name, isTag bool, isOpenTag bool) {
	item := elements[i]
	if se, ok := item.(xml.StartElement); ok == true {
		return se.Name, true, true
	} else if se, ok := item.(xml.EndElement); ok == true {
		return se.Name, true, false
	}

	return name, false, false
}

// validateAnchorElements asserts that the list of elements starts and ends with
// the given tag. Note that any failures here are likely due to mistyping a
// field in such a way that we expect or don't expect char-data between tags
// when we shouldn't.
func validateAnchorElements(elements []interface{}, name xml.Name) (err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	elementCount := len(elements)

	if elementCount < 2 {
		log.Panicf("expected at least two items for anchor-tag check")
	}

	firstElementTagName, firstElementIsTag, firstElementIsTagOpen := elementTagName(elements, 0)

	if firstElementIsTag == false {
		log.Panicf("expected first element in array to be a tag")
	} else if firstElementIsTagOpen != true {
		log.Panicf("expected first tag to be an open-tag")
	} else if firstElementTagName != name {
		log.Panicf(
			"expected first element in array to be a [%s] tag: [%s]",
			xmpregistry.XmlName(name), xmpregistry.XmlName(firstElementTagName))
	}

	lastElementTagName, lastElementIsTag, firstElementIsTagOpen := elementTagName(elements, elementCount-1)

	printElements := func() {
		fmt.Printf("Element dump:\n")
		fmt.Printf("\n")

		for i, x := range elements {
			fmt.Printf("%d: [%v] [%v]\n", i, reflect.TypeOf(x), x)
		}
	}

	if lastElementIsTag == false {
		fmt.Printf("\n")

		printElements()
		fmt.Printf("\n")

		e := elements[len(elements)-1]
		log.Panicf("expected last element in array to be a tag: [%v] [%v]", reflect.TypeOf(e), e)
	} else if firstElementIsTagOpen != false {
		fmt.Printf("\n")

		printElements()
		fmt.Printf("\n")

		log.Panicf("expected last tag to be a close-tag")
	} else if lastElementTagName != name {
		fmt.Printf("\n")

		printElements()
		fmt.Printf("\n")

		log.Panicf(
			"expected last element in array to be a [%s] tag: [%s]",
			xmpregistry.XmlName(name), xmpregistry.XmlName(lastElementTagName))
	}

	return nil
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

// FullName returns the fully-qualified name of the node encapsulating the
// array.
func (bav baseArrayValue) FullName() xmpregistry.XmpPropertyName {
	return bav.fullName
}

// Count returns the number of items found/extracted from the array.
func (bav baseArrayValue) Count() int {
	return len(bav.collected)
}

func (bav baseArrayValue) constructArrayItem(subslice []interface{}) (ai ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	subsliceLen := len(subslice)
	if subsliceLen != 3 {
		log.Panicf("sublice length is not valid: (%d)", subsliceLen)
	}

	err = validateAnchorElements(subslice, rdfLiTag)
	log.PanicIf(err)

	se := subslice[0].(xml.StartElement)

	attributes, err := ParseAttributes(se)
	log.PanicIf(err)

	var charData string

	if len(subslice) == 3 {
		// There is character-data between the tags. Extract it.

		charDataRaw := subslice[1]

		var ok bool

		charData, ok = charDataRaw.(string)
		if ok == false {
			log.Panicf(
				"expected element between 'li' tags in unordered-array to be char-data: [%s] [%s]",
				bav.FullName(), reflect.TypeOf(charData))
		}
	}

	ai = ArrayItem{
		Name:       se.Name,
		Attributes: attributes,
		CharData:   charData,
	}

	return ai, nil
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

	firstTag := bav.collected[0].(xml.StartElement)
	firstTagName := firstTag.Name
	lastTag := bav.collected[len(bav.collected)-1].(xml.EndElement)
	lastTagName := lastTag.Name

	if firstTagName != lastTagName {
		log.Panicf(
			"open and close anchor tags do not have the same name: [%s] [%s] != [%s] [%s]",
			firstTagName.Space, firstTagName.Local, lastTagName.Space, lastTagName.Local)
	}

	for i := 1; i < elementCount-1; {
		var subsliceLen int

		if hasSandwichedCharData == true {
			// Two tags with char-data in the middle (three elements total).
			subsliceLen = 3
		} else {
			// Two tags with no char-data in the middle (two elements total).
			subsliceLen = 2
		}

		subslice := bav.collected[i : i+subsliceLen]
		if len(subslice) != subsliceLen {
			log.Panicf("number of elements does not make sense (with char-data)")
		}

		// constructArrayItem will assume that the second item is char-data if
		// there are three items.

		ai, err := bav.constructArrayItem(subslice)
		log.PanicIf(err)

		items = append(items, ai)

		i += len(subslice)
	}

	return items, nil
}

// Ordered array semantics

// TODO(dustin): Ordered array yet-to-implement: CuePointParam, Marker, ResourceEvent, Version, Colorant, Marker, Layer, "point" (?)

// OrderedArrayValue represents the items of an ordered-array.
type OrderedArrayValue struct {
	baseArrayValue
}

func newOrderedArrayValue(bav baseArrayValue) OrderedArrayValue {
	return OrderedArrayValue{
		baseArrayValue: bav,
	}
}

// String returns a string representation of the array state.
func (oav OrderedArrayValue) String() string {
	return fmt.Sprintf("OrderedArray<COUNT=(%d)>", oav.Count())
}

// Items returns a slice of all raw array items.
func (oav OrderedArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = validateAnchorElements(oav.baseArrayValue.collected, rdfSeqTag)
	log.PanicIf(err)

	items, err = oav.innerItems(true)
	log.PanicIf(err)

	return items, nil
}

// OrderedArrayFieldType is a field-type that acts as a factory for the ordered-
// array value type.
type OrderedArrayFieldType struct {
}

// New returns a value-type for the given arguments.
func (oat OrderedArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return newOrderedArrayValue(bav)
}

// OrderedTextArrayValue identifies the array as having resource-event
// items.
type OrderedTextArrayValue struct {
	OrderedArrayValue
}

// Items this is a wrapper that returns a simple list of strings from inner
// underlying array-items, thereby satisfying the ArrayStringValueLister
// interface. In the case of these, we return a stringification of the
// attributes.
func (otav OrderedTextArrayValue) Items() (items []string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	innerItems, err := otav.OrderedArrayValue.Items()
	log.PanicIf(err)

	items = make([]string, len(innerItems))
	for i, ai := range innerItems {
		items[i] = ai.CharData
	}

	return items, nil
}

// OrderedTextArrayFieldType identifies the array as having text items.
type OrderedTextArrayFieldType struct {
	OrderedArrayFieldType
}

// New returns a value-type for the given arguments.
func (oat OrderedTextArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	oav := newOrderedArrayValue(bav)

	return OrderedTextArrayValue{
		OrderedArrayValue: oav,
	}
}

// OrderedUriArrayFieldType identifies the array as having URI items.
type OrderedUriArrayFieldType struct {
	OrderedArrayFieldType
}

// OrderedResourceEventArrayValue identifies the array as having resource-event
// items.
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

// OrderedResourceEventArrayFieldType identifies the array as having resource-
// event items.
type OrderedResourceEventArrayFieldType struct {
}

// New returns a value-type for the given arguments.
func (oreat OrderedResourceEventArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	oav := newOrderedArrayValue(bav)

	return OrderedResourceEventArrayValue{
		OrderedArrayValue: oav,
	}
}

// Unordered array semantics

// TODO(dustin): Unordered array yet-to-implement: XPath, ResourceRef, "struct" (?), Job, Font, Media, Track

// UnorderedArrayValue represents the items of an unordered-array.
type UnorderedArrayValue struct {
	baseArrayValue
}

func newUnorderedArrayValue(bav baseArrayValue) UnorderedArrayValue {
	return UnorderedArrayValue{
		baseArrayValue: bav,
	}
}

// String returns a string representation of the array state.
func (uav UnorderedArrayValue) String() string {
	return fmt.Sprintf("UnorderedArray<COUNT=(%d)>", uav.Count())
}

// Items returns a slice of all raw array items.
func (uav UnorderedArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = validateAnchorElements(uav.baseArrayValue.collected, rdfBagTag)
	log.PanicIf(err)

	items, err = uav.innerItems(true)
	log.PanicIf(err)

	return items, nil
}

// UnorderedArrayFieldType is a field-type that acts as a factory for the
// unordered-array value type.
type UnorderedArrayFieldType struct {
}

// New returns a value-type for the given arguments.
func (uat UnorderedArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return newUnorderedArrayValue(bav)
}

// UnorderedTextArrayValue represents the items of an unordered-array with
// Ancestor items.
type UnorderedTextArrayValue struct {
	UnorderedArrayValue
}

// Items this is a wrapper that returns a simple list of strings from inner
// underlying array-items, thereby satisfying the ArrayStringValueLister
// interface.
func (utav UnorderedTextArrayValue) Items() (items []string, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	innerItems, err := utav.UnorderedArrayValue.Items()
	log.PanicIf(err)

	items = make([]string, len(innerItems))
	for i, ai := range innerItems {
		items[i] = ai.CharData
	}

	return items, nil
}

// UnorderedTextArrayFieldType identifies the array as having text items.
type UnorderedTextArrayFieldType struct {
	UnorderedArrayFieldType
}

// New returns a value-type for the given arguments.
func (uat UnorderedTextArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	uav := newUnorderedArrayValue(bav)

	return UnorderedTextArrayValue{
		UnorderedArrayValue: uav,
	}
}

// UnorderedAncestorArrayValue represents the items of an unordered-array with
// Ancestor items.
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

// UnorderedAncestorArrayFieldType identifies the array as having Ancestor
// items.
type UnorderedAncestorArrayFieldType struct {
	UnorderedArrayFieldType
}

// New returns a value-type for the given arguments.
func (uaat UnorderedAncestorArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return UnorderedAncestorArrayValue{
		UnorderedArrayValue: newUnorderedArrayValue(bav),
	}
}

// Alternatives array semantics

// AlternativeArrayValue represents the items of an alternatives-array
type AlternativeArrayValue struct {
	baseArrayValue
}

func newAlternativeArrayValue(bav baseArrayValue) AlternativeArrayValue {
	return AlternativeArrayValue{
		baseArrayValue: bav,
	}
}

// String returns a string representation of the array state.
func (aav AlternativeArrayValue) String() string {
	return fmt.Sprintf("AlternativeArray<COUNT=(%d)>", aav.Count())
}

// Items returns a slice of all raw array items.
func (aav AlternativeArrayValue) Items() (items []ArrayItem, err error) {
	defer func() {
		if errRaw := recover(); errRaw != nil {
			err = log.Wrap(errRaw.(error))
		}
	}()

	err = validateAnchorElements(aav.baseArrayValue.collected, rdfAltTag)
	log.PanicIf(err)

	items, err = aav.innerItems(true)
	log.PanicIf(err)

	return items, nil
}

// AlternativeArrayFieldType is a field-type that acts as a factory for the
// actual value type.
type AlternativeArrayFieldType struct {
}

// New returns a value-type for the given arguments.
func (aat AlternativeArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)

	return AlternativeArrayValue{
		baseArrayValue: bav,
	}
}

// LanguageAlternativeArrayValue represents the items of an alternatives-array.
type LanguageAlternativeArrayValue struct {
	AlternativeArrayValue
}

// Items this is a wrapper that returns a simple list of strings from inner
// underlying array-items, thereby satisfying the ArrayStringValueLister
// interface.
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

// LanguageAlternativeArrayFieldType is a field-type that acts as a factory for
// the alternatives-array value type.
type LanguageAlternativeArrayFieldType struct {
}

// New returns a value-type for the given arguments.
func (laat LanguageAlternativeArrayFieldType) New(fullName xmpregistry.XmpPropertyName, collected []interface{}) ArrayValue {
	bav := newBaseArrayValue(fullName, collected)
	aav := newAlternativeArrayValue(bav)

	return LanguageAlternativeArrayValue{
		AlternativeArrayValue: aav,
	}
}
