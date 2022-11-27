package commands

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alinnert/xt/log"
	"github.com/alinnert/xt/utils"
	"github.com/alinnert/xt/xsd"
	"github.com/dominikbraun/graph"
	"github.com/spf13/cobra"
)

type adjacencyMap[K comparable] map[K]map[K]graph.Edge[K]

// MainCommand finds the shortest path from all source elements to the given target element.
func MainCommand() *cobra.Command {
	var limit int
	var exact bool
	var verbose bool

	mainCommand := &cobra.Command{
		Use:     "xt [file] [element]",
		Short:   "XML Schema Tools",
		Long:    "XML Schema Tools displays information about an XML Schema.",
		Version: "0.2",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			file, element := parseArgs(args)

			filesGraph, err := xsd.ResolveIncludes(file, verbose)
			if err != nil {
				panic(err)
			}

			elementsGraph, err := xsd.GetElements(filesGraph, verbose)
			if err != nil {
				panic(err)
			}

			adjMap, _ := elementsGraph.AdjacencyMap()
			sourceElements := getSourceElements(adjMap)
			targetElements := getTargetElements(element, adjMap, exact)
			results := getResults(elementsGraph, sourceElements, targetElements)
			printResults(results, element, limit, verbose)
		},
	}

	mainCommand.Flags().IntVarP(&limit, "limit", "l", 5, "Limit the number of results. Only the shortest results are shown. Use \"--limit 0\" to show all results.")
	mainCommand.Flags().BoolVarP(&exact, "exact", "e", false, "If flag is set and search term is \"elem\" only \"elem\" is found. Otherwise \"parent/elem\" is also found.")
	mainCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Output additional information about the parsed XML Schema.")

	return mainCommand
}

func parseArgs(args []string) (string, string) {
	filepathArg, err := filepath.Abs(args[0])
	if err != nil {
		panic(err)
	}
	elementArg := args[1]

	return filepathArg, elementArg
}

func getSourceElements(adjMap adjacencyMap[string]) *[]string {
	var sourceElements []string

	for rootElementName := range adjMap[""] {
		sourceElements = append(sourceElements, rootElementName)
	}

	return &sourceElements
}

func getTargetElements(element string, adjMap adjacencyMap[string], exact bool) *[]string {
	targetElements := []string{element}

	if !exact {
		for element := range adjMap {
			if !strings.HasSuffix(element, "/"+element) {
				continue
			}
			targetElements = append(targetElements, element)
		}
	}

	return &targetElements
}

func getResults(elementsGraph xsd.ElementsGraph, sourceElements, targetElements *[]string) *[][]string {
	var results [][]string

	for _, targetElement := range *targetElements {
		for _, element := range *sourceElements {
			path, err := graph.ShortestPath(elementsGraph, element, targetElement)
			if err != nil {
				continue
			}

			results = append(results, path)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return len(results[i]) < len(results[j])
	})

	return &results
}

func printResults(results *[][]string, element string, limit int, verbose bool) {
	if verbose {
		fmt.Println()
	}

	if len(*results) == 0 {
		fmt.Printf("No paths found for element \"" + element + "\"\n")
		return
	}

	resultsCount := utils.Min(limit, len(*results))
	if limit == 0 {
		resultsCount = len(*results)
	}

	countLabel := fmt.Sprintf("(showing all %d)", resultsCount)

	if resultsCount != len(*results) {
		countLabel = fmt.Sprintf("(showing %d of %d)", resultsCount, len(*results))
	}

	fmt.Printf("Possible paths for element \"%s\" %s\n", element, countLabel)

	for _, result := range (*results)[:resultsCount] {
		log.PathsResult(result)
	}
}
