package log

import (
	"fmt"
	"strings"

	"github.com/alinnert/xt/utils"
	"github.com/fatih/color"
)

var childDivider = color.YellowString("/")
var referenceDivider = color.MagentaString("=>")

func addBrackets(elementName string) string { return "<" + elementName + ">" }

func colorizeSegment(elementName string, highlight bool) string {
	if highlight {
		return addBrackets(color.CyanString(elementName))
	}
	return addBrackets(color.WhiteString(elementName))
}

func colorizeElement(elementName string) string {
	segments := strings.Split(elementName, "/")

	if len(segments) == 1 {
		return colorizeSegment(elementName, true)
	}

	colorizedSegments := utils.Map(&segments, func(ctx *utils.MapContext[string]) string {
		if ctx.LastItem {
			return colorizeSegment(*ctx.Item, true)
		}
		return colorizeSegment(*ctx.Item, false)
	})

	return strings.Join(colorizedSegments, color.HiBlackString("/"))
}

func colorizePath(path []string) string {
	colorizedPath := utils.Map(&path, func(ctx *utils.MapContext[string]) string {
		segments := strings.Split(*ctx.Item, "/")
		lastSegment := segments[len(segments)-1]
		colorizedLastSegment := colorizeElement(lastSegment)
		if ctx.FirstItem {
			return colorizedLastSegment
		}
		if len(segments) > 1 {
			return fmt.Sprintf("%s %s", childDivider, colorizedLastSegment)
		}
		return fmt.Sprintf("%s %s", referenceDivider, colorizedLastSegment)
	})

	return strings.Join(colorizedPath, " ")
}

func colorizeAction(action string) string {
	return color.HiWhiteString(action + ":")
}

// AddFile logs a file when it gets loaded while following all `include` elements from an XML Schema file.
func AddFile(filePath string) {
	fmt.Println(colorizeAction("Add file"), color.GreenString(filePath))
}

// DuplicateFile logs a file that gets loaded at least for a second time.
func DuplicateFile(filePath string) {
	fmt.Println(colorizeAction("Duplicate file"), color.YellowString(filePath))
}

// ElementHeadline prints a headline for an element that gets processed.
func ElementHeadline(elementName string) {
	fmt.Println()
	fmt.Println(addBrackets(color.HiYellowString(elementName)))
}

// AddElement logs an element that was found while parsing an XML Schema file.
func AddElement(elementName string) {
	fmt.Println(colorizeAction("Add element"), colorizeElement(elementName))
}

// DuplicateElement logs an element that was found at least for a second time.
func DuplicateElement(elementName string) {
	fmt.Println(colorizeAction("Duplicate element"), colorizeElement(elementName))
}

// LeafElementsCount prints the number of leaf elements that were found inside a given root level element.
func LeafElementsCount(count int) {
	c := color.New(color.FgGreen).SprintFunc()
	if count == 1 {
		fmt.Println("1", c("leaf element"))
	} else {
		fmt.Println(color.HiWhiteString("%d", count), c("leaf elements"))
	}
}

// ElementRefsCount prints the number of element references that were found inside a given root level element.
func ElementRefsCount(count int) {
	c := color.New(color.FgGreen).SprintFunc()
	if count == 1 {
		fmt.Println("1", c("element ref"))
	} else {
		fmt.Println(color.HiWhiteString("%d", count), c("element refs"))
	}
}

// AddReference logs a reference that is being created.
func AddReference(from, to string) {
	fmt.Println(colorizeAction("Add reference"), colorizeElement(from), referenceDivider, colorizeElement(to))
}

// PathsResult prints a result line for the main command.
func PathsResult(segments []string) {
	fmt.Println("-", colorizePath(segments))
}
