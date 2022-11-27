package xsd

import "fmt"

var xsdNamespace = "http://www.w3.org/2001/XMLSchema"

// Element creates a part of an xpath selector that selects an element with a given name in the xml schema namespace.
func Element(name string) string {
	return fmt.Sprintf("*[local-name()='%s' and namespace-uri()='%s']", name, xsdNamespace)
}
