package xsd

var xsdNamespace = "http://www.w3.org/2001/XMLSchema"

// Element creates a part of an xpath selector that selects an element with a given name in the xml schema namespace.
func Element(name string) string {
	return "*[local-name()='" + name + "' and namespace-uri()='" + xsdNamespace + "']"
}
