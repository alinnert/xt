package xsd

import "strings"

// XsdElement represents an XSD element.
type XsdElement struct {
	PathSegments []string
}

// Name returns the effective element name excluding the path.
func (e XsdElement) Name() string {
	return e.PathSegments[len(e.PathSegments)-1]
}

// PathString returns the element's path as a string, divided by a "/".
func (e XsdElement) PathString() string {
	return strings.Join(e.PathSegments, "/")
}
