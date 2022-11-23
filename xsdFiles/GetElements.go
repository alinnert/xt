package xsdFiles

import (
	"fmt"
	"strings"

	"github.com/alinnert/xt/utils"
	"github.com/alinnert/xt/xsd"
	"github.com/antchfx/xmlquery"
	"github.com/dominikbraun/graph"
)

type ElementsGraph = graph.Graph[string, utils.XsdElement]

// GetElements is the entry point of the element processing logic.
func GetElements(filesGraph XsdFileGraph, verbose bool) (ElementsGraph, error) {
	elementsGraph := graph.New(func(element utils.XsdElement) string {
		return element.PathString()
	}, graph.Directed())

	// Add the main root element
	err := elementsGraph.AddVertex(utils.XsdElement{PathSegments: []string{}})
	if err != nil {
		return nil, err
	}

	filesMap, err := filesGraph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	for filePath := range filesMap {
		err = findRootElements(filePath, filesGraph, elementsGraph, verbose)
		if err != nil {
			return nil, err
		}
	}

	for filePath := range filesMap {
		err = findNestedElements(filePath, filesGraph, elementsGraph, verbose)
		if err != nil {
			return nil, err
		}
	}

	return elementsGraph, nil
}

// findRootElements finds and processes all root elements.
func findRootElements(filePath string, filesGraph XsdFileGraph, elementsGraph ElementsGraph, verbose bool) error {
	fileVertex, err := filesGraph.Vertex(filePath)
	if err != nil {
		return err
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(fileVertex.Content)))
	if err != nil {
		return err
	}

	rootElementsQuery := "/" + xsd.Element("schema") + "/" + xsd.Element("element") + "[@name]"

	rootElements, err := xmlquery.QueryAll(doc, rootElementsQuery)
	if err != nil {
		return err
	}

	for _, rootElement := range rootElements {
		rootElementName := rootElement.SelectAttr("name")

		if _, err := elementsGraph.Vertex(rootElementName); err == nil {
			if verbose {
				fmt.Println("  duplicate element:", rootElementName)
			}

			return nil
		}

		// Add root document elements
		if verbose {
			fmt.Println("add root element", rootElementName)
		}

		err = elementsGraph.AddVertex(utils.XsdElement{PathSegments: []string{rootElementName}})
		if err != nil {
			return err
		}

		err = elementsGraph.AddEdge("", rootElementName)
		if err != nil {
			return err
		}
	}

	return nil
}

// findNestedElements calls findLeafElements and findLeafElementRefs.
func findNestedElements(filePath string, filesGraph XsdFileGraph, elementsGraph ElementsGraph, verbose bool) error {
	fileVertex, err := filesGraph.Vertex(filePath)
	if err != nil {
		return err
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(fileVertex.Content)))
	if err != nil {
		return err
	}

	rootElementsQuery := "/" + xsd.Element("schema") + "/" + xsd.Element("element") + "[@name]"

	rootElements, err := xmlquery.QueryAll(doc, rootElementsQuery)
	if err != nil {
		return err
	}

	for _, rootElement := range rootElements {
		if verbose {
			fmt.Println("[", rootElement.SelectAttr("name"), "]")
		}

		err := findLeafElements(rootElement, elementsGraph, verbose)
		if err != nil {
			return err
		}

		err = findLeafElementRefs(rootElement, elementsGraph, verbose)
		if err != nil {
			return err
		}

		if verbose {
			fmt.Println()
		}
	}

	return nil
}

// findLeafElements finds and processes all leaf elements. Those are the most deeply nested ones.
func findLeafElements(rootElement *xmlquery.Node, elementsGraph ElementsGraph, verbose bool) error {
	leafElementsQuery := "/descendant::" + xsd.Element("element") + "[not(descendant::" + xsd.Element("element") + "[@name])][@name]"
	leafElements, err := xmlquery.QueryAll(rootElement, leafElementsQuery)
	if err != nil {
		return err
	}

	if verbose {
		if lenLeafElements := len(leafElements); lenLeafElements == 1 {
			fmt.Println("1 leaf element")
		} else {
			fmt.Println(lenLeafElements, "leaf elements")
		}
	}

	for _, leafElement := range leafElements {
		err = findLeafAncestorElements(leafElement, elementsGraph, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// findLeafElementRefs finds and processes all leaf elements that are references to root elements.
func findLeafElementRefs(rootElement *xmlquery.Node, elementsGraph ElementsGraph, verbose bool) error {
	leafElementRefsQuery := "/descendant::" + xsd.Element("element") + "[@ref]"
	leafElementRefs, err := xmlquery.QueryAll(rootElement, leafElementRefsQuery)
	if err != nil {
		return err
	}

	if verbose {
		if lenLeafElementRefs := len(leafElementRefs); lenLeafElementRefs == 1 {
			fmt.Println("1 leaf element ref")
		} else {
			fmt.Println(lenLeafElementRefs, "leaf element refs")
		}
	}

	for _, elementRef := range leafElementRefs {
		referencedElementName := elementRef.SelectAttr("ref")
		if strings.Contains(referencedElementName, ":") {
			continue
		}

		closestAncestorQuery := "/ancestor::" + xsd.Element("element") + "[@name][1]"
		closestAncestor, err := xmlquery.Query(elementRef, closestAncestorQuery)
		if err != nil {
			return err
		}

		closestAncestorName := closestAncestor.SelectAttr("name")

		if verbose {
			fmt.Println("add edge for reference", closestAncestorName, "->", referencedElementName)
		}

		err = elementsGraph.AddEdge(closestAncestorName, referencedElementName)
		if err != nil {
			continue
		}
	}

	return nil
}

// findLeafAncestorElements finds and processes all ancestors of a leaf element.
func findLeafAncestorElements(leafElement *xmlquery.Node, elementsGraph ElementsGraph, verbose bool) error {
	leafElementName := leafElement.SelectAttr("name")

	ancestorsQuery := "/ancestor::" + xsd.Element("element") + "[@name]"
	ancestors, err := xmlquery.QueryAll(leafElement, ancestorsQuery)
	if err != nil {
		return err
	}

	ancestors = utils.Reverse(ancestors)

	ancestorElementNames := utils.Map(ancestors, func(ancestor *xmlquery.Node) string {
		return ancestor.SelectAttr("name")
	})

	leafElementPathSegments := append(ancestorElementNames, leafElementName)
	ancestorPaths := xsd.GetAncestorPathSegments(leafElementPathSegments)

	for _, currentItemPathSegments := range ancestorPaths {
		currentItemPath := strings.Join(currentItemPathSegments, "/")
		parentPathSegments := xsd.GetParent(currentItemPathSegments)
		parentPath := strings.Join(parentPathSegments, "/")

		if _, err := elementsGraph.Vertex(currentItemPath); err == nil {
			continue
		}

		// Add descendant elements
		if verbose {
			fmt.Println("add nested element", currentItemPath)
		}

		err = elementsGraph.AddVertex(utils.XsdElement{PathSegments: currentItemPathSegments})
		if err != nil {
			return err
		}

		if verbose {
			fmt.Println("add edge for nested element", parentPath, "->", currentItemPath)
		}

		err = elementsGraph.AddEdge(parentPath, currentItemPath)
		if err != nil {
			return err
		}
	}

	return nil
}
