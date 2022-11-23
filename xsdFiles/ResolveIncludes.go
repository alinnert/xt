package xsdFiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alinnert/xt/xsd"
	"github.com/dominikbraun/graph"

	"github.com/antchfx/xmlquery"
)

type XsdFile struct {
	Path    string
	Content []byte
}

type XsdFileGraph = graph.Graph[string, XsdFile]

// ResolveIncludes reads the contents of all included files in an XSD file.
func ResolveIncludes(mainFilePath string, verbose bool) (XsdFileGraph, error) {
	fileGraph := graph.New(func(xsdFile XsdFile) string {
		return xsdFile.Path
	})

	err := processFile(fileGraph, mainFilePath, verbose)
	if err != nil {
		return nil, err
	}

	return fileGraph, nil
}

func processFile(fileGraph XsdFileGraph, filePath string, verbose bool) error {
	if verbose {
		fmt.Println("Add file:", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = fileGraph.AddVertex(XsdFile{Path: filePath, Content: content})
	if err != nil {
		return err
	}

	document, err := xmlquery.Parse(strings.NewReader(string(content)))
	if err != nil {
		return err
	}

	includes, err := xmlquery.QueryAll(document, "//"+xsd.Element("include")+"[@schemaLocation]")
	if err != nil {
		return err
	}

	currentDir := filepath.Dir(filePath)

	for _, include := range includes {
		nextFilePath := filepath.Join(currentDir, include.SelectAttr("schemaLocation"))

		if _, err := fileGraph.Vertex(nextFilePath); err == nil {
			fmt.Println("duplicate include:", nextFilePath)
			continue
		}

		err := processFile(fileGraph, nextFilePath, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}
