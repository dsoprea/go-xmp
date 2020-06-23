// Package xmpnamespace describe hierarchies found within XMP documents that are
// associated with namespaces.
//
// Nodes in an unregistered namespace will stringify with a prefix of "?".
// Values in an unregistered namespace will either be parsed as a string or
// skipped.
//
// The standard Go XML parser is used. Even if a field is not defined in the
// field-mapping for a particular namespace, it will still be descended into by
// the Go parser.
package xmpnamespace
