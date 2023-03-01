package xsd

import (
	"fmt"
	"strings"

	"github.com/alinnert/xt/log"
	"github.com/alinnert/xt/utils"
	"github.com/antchfx/xmlquery"
	"github.com/dominikbraun/graph"
)

// ElementsGraph is a graph.Graph that includes all found element definitions and their relations to each other.
type ElementsGraph = graph.Graph[string, XsdElement]

// GetElements is the entry point of the element processing logic.
func GetElements(filesGraph FilesGraph, verbose bool) (ElementsGraph, error) {
	elementsGraph := graph.New(func(element XsdElement) string {
		return element.PathString()
	}, graph.Directed())

	// Add the main root element
	err := elementsGraph.AddVertex(XsdElement{PathSegments: []string{}})
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
func findRootElements(filePath string, filesGraph FilesGraph, elementsGraph ElementsGraph, verbose bool) error {
	fileVertex, err := filesGraph.Vertex(filePath)
	if err != nil {
		return err
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(fileVertex.Content)))
	if err != nil {
		return err
	}

	rootElementsQuery := "/" + Element("schema") + "/" + Element("element") + "[@name]"

	rootElements, err := xmlquery.QueryAll(doc, rootElementsQuery)
	if err != nil {
		return err
	}

	for _, rootElement := range rootElements {
		rootElementName := rootElement.SelectAttr("name")

		if _, err := elementsGraph.Vertex(rootElementName); err == nil {
			if verbose {
				log.DuplicateElement(rootElementName)
			}

			return nil
		}

		// Add root document elements
		if verbose {
			log.AddElement(rootElementName)
		}

		err = elementsGraph.AddVertex(XsdElement{PathSegments: []string{rootElementName}})
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
func findNestedElements(filePath string, filesGraph FilesGraph, elementsGraph ElementsGraph, verbose bool) error {
	fileVertex, err := filesGraph.Vertex(filePath)
	if err != nil {
		return err
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(fileVertex.Content)))
	if err != nil {
		return err
	}

	rootElementsQuery := "/" + Element("schema") + "/" + Element("element") + "[@name]"

	rootElements, err := xmlquery.QueryAll(doc, rootElementsQuery)
	if err != nil {
		return err
	}

	for _, rootElement := range rootElements {
		if verbose {
			log.ElementHeadline(rootElement.SelectAttr("name"))
		}

		err := findLeafElements(rootElement, elementsGraph, verbose)
		if err != nil {
			return err
		}

		err = findLeafElementRefs(rootElement, elementsGraph, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

// findLeafElements finds and processes all leaf elements. Those are the most deeply nested ones.
func findLeafElements(rootElement *xmlquery.Node, elementsGraph ElementsGraph, verbose bool) error {
	leafElementsQuery := "/descendant::" + Element("element") + "[not(descendant::" + Element("element") + "[@name])][@name]"
	leafElements, err := xmlquery.QueryAll(rootElement, leafElementsQuery)
	if err != nil {
		return err
	}

	if verbose {
		log.LeafElementsCount(len(leafElements))
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
	leafElementRefsQuery := "/descendant::" + Element("element") + "[@ref]"
	leafElementRefs, err := xmlquery.QueryAll(rootElement, leafElementRefsQuery)
	if err != nil {
		return err
	}

	if verbose {
		log.ElementRefsCount(len(leafElementRefs))
	}

	for _, elementRef := range leafElementRefs {
		referencedElementName := elementRef.SelectAttr("ref")
		if strings.Contains(referencedElementName, ":") {
			continue
		}

		closestAncestorQuery := "/ancestor::" + Element("element") + "[@name][1]"
		closestAncestor, err := xmlquery.Query(elementRef, closestAncestorQuery)
		if err != nil {
			return err
		}

		closestAncestorName := closestAncestor.SelectAttr("name")

		if verbose {
			log.AddReference(closestAncestorName, referencedElementName)
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

	ancestorsQuery := "/ancestor::" + Element("element") + "[@name]"
	ancestors, err := xmlquery.QueryAll(leafElement, ancestorsQuery)
	if err != nil {
		return err
	}

	ancestors = utils.Reverse(ancestors)

	ancestorElementNames := utils.Map(&ancestors, func(ctx *utils.MapContext[*xmlquery.Node]) string {
		return (*ctx.Item).SelectAttr("name")
	})

	leafElementPathSegments := append(ancestorElementNames, leafElementName)
	ancestorPaths := GetAncestorPathSegments(leafElementPathSegments)

	for _, currentItemPathSegments := range ancestorPaths {
		currentItemPath := strings.Join(currentItemPathSegments, "/")
		parentPathSegments := GetParent(currentItemPathSegments)
		parentPath := strings.Join(parentPathSegments, "/")

		if _, err := elementsGraph.Vertex(currentItemPath); err == nil {
			continue
		}

		// Add descendant elements
		if verbose {
			log.AddElement(currentItemPath)
		}

		err = elementsGraph.AddVertex(XsdElement{PathSegments: currentItemPathSegments})
		if err != nil {
			return err
		}

		if verbose {
			log.AddReference(parentPath, currentItemPath)
			fmt.Println()
		}

		err = elementsGraph.AddEdge(parentPath, currentItemPath)
		if err != nil {
			return err
		}
	}

	return nil
}
