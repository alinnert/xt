package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alinnert/xt/utils"
	"github.com/alinnert/xt/xsdFiles"
	"github.com/dominikbraun/graph"
	"github.com/spf13/cobra"
)

func main() {
	var limit int
	var exact bool
	var verbose bool

	// xt some-file.xsd --limit 5 --verbose
	rootCmd := &cobra.Command{
		Use:   "xt",
		Short: "XML Schema Tools",
		Long:  "XML Schema Tools displays information about an XML Schema.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			filepathArg, err := filepath.Abs(args[0])
			if err != nil {
				panic(err)
			}
			elementArg := args[1]

			filesGraph, err := xsdFiles.ResolveIncludes(filepathArg, verbose)
			if err != nil {
				panic(err)
			}

			elementsGraph, err := xsdFiles.GetElements(filesGraph, verbose)
			if err != nil {
				panic(err)
			}

			var sourceElements []string
			targetElements := []string{elementArg}

			adjMap, _ := elementsGraph.AdjacencyMap()

			// Get all possible source vertex hashes.
			for rootElementName := range adjMap[""] {
				sourceElements = append(sourceElements, rootElementName)
			}

			// Get all possible target vertex hashes.
			if !exact {
				for element := range adjMap {
					if !strings.HasSuffix(element, "/"+elementArg) {
						continue
					}
					targetElements = append(targetElements, element)
				}
			}

			// Fetch all results.
			var results [][]string

			for _, targetElement := range targetElements {
				for _, element := range sourceElements {
					path, err := graph.ShortestPath(elementsGraph, element, targetElement)
					if err != nil {
						continue
					}

					results = append(results, path)
				}
			}

			// Sort results
			sort.Slice(results, func(i, j int) bool {
				return len(results[i]) < len(results[j])
			})

			// Print all results.
			if len(results) == 0 {
				fmt.Printf("No paths found for element \"" + elementArg + "\"\n")
				return
			}

			resultsCount := utils.MinInt(limit, len(results))
			if limit == 0 {
				resultsCount = len(results)
			}

			countLabel := fmt.Sprintf("(showing all %d)", resultsCount)

			if resultsCount != len(results) {
				countLabel = fmt.Sprintf("(showing %d of %d)", resultsCount, len(results))
			}

			fmt.Printf("Possible paths for element \"%s\" %s\n", elementArg, countLabel)

			for _, result := range results[:resultsCount] {
				fmt.Println("-", strings.Join(result, " > "))
			}
		},
	}

	rootCmd.Flags().IntVarP(
		&limit,
		"limit",
		"l",
		5,
		"Limit the number of results. Only the shortest results are shown.",
	)

	rootCmd.Flags().BoolVarP(
		&exact,
		"exact",
		"e",
		false,
		"If flag is set and search term is \"elem\" only \"elem\" is found. Otherwise \"parent/elem\" is also found.",
	)

	rootCmd.Flags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"Output additional information about the parsed XML Schema.",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
