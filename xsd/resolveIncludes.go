package xsd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alinnert/xt/log"
	"github.com/dominikbraun/graph"

	"github.com/antchfx/xmlquery"
)

// XsdFile represents an XML Schema file. It contains the filename as well as the file content.
type XsdFile struct {
	Path    string
	Content []byte
}

// FilesGraph is the graph.Graph that contains all loaded XML Schema files.
type FilesGraph = graph.Graph[string, XsdFile]

// ResolveIncludes reads the contents of all included files in an XSD file.
func ResolveIncludes(mainFilePath string, verbose bool) (FilesGraph, error) {
	filesGraph := graph.New(func(xsdFile XsdFile) string {
		return xsdFile.Path
	})

	err := processFile(filesGraph, mainFilePath, verbose)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Println()
	}

	return filesGraph, nil
}

func processFile(fileGraph FilesGraph, filePath string, verbose bool) error {
	if verbose {
		log.AddFile(filePath)
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

	includes, err := xmlquery.QueryAll(document, "//"+Element("include")+"[@schemaLocation]")
	if err != nil {
		return err
	}

	currentDir := filepath.Dir(filePath)

	for _, include := range includes {
		nextFilePath := filepath.Join(currentDir, include.SelectAttr("schemaLocation"))

		if _, err := fileGraph.Vertex(nextFilePath); err == nil {
			log.DuplicateFile(nextFilePath)
			continue
		}

		err := processFile(fileGraph, nextFilePath, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}
