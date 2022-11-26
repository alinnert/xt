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

	colorizedSegments := utils.Map(segments, func(ctx utils.MapContext[string]) string {
		if ctx.LastItem {
			return colorizeSegment(ctx.Item, true)
		}
		return colorizeSegment(ctx.Item, false)
	})

	return strings.Join(colorizedSegments, color.HiBlackString("/"))
}

func colorizePath(path []string) string {
	colorizedPath := utils.Map(path, func(ctx utils.MapContext[string]) string {
		segments := strings.Split(ctx.Item, "/")
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

func AddFile(filePath string) {
	fmt.Println(colorizeAction("Add file"), color.GreenString(filePath))
}

func DuplicateFile(filePath string) {
	fmt.Println(colorizeAction("Duplicate file"), color.YellowString(filePath))
}

func ElementHeadline(elementName string) {
	fmt.Println()
	fmt.Println(addBrackets(color.HiYellowString(elementName)))
}

func AddElement(elementName string) {
	fmt.Println(colorizeAction("Add element"), colorizeElement(elementName))
}

func DuplicateElement(elementName string) {
	fmt.Println(colorizeAction("Duplicate element"), colorizeElement(elementName))
}

func LeafElementsCount(count int) {
	c := color.New(color.FgGreen).SprintFunc()
	if count == 1 {
		fmt.Println("1", c("leaf element"))
	} else {
		fmt.Println(color.HiWhiteString("%d", count), c("leaf elements"))
	}
}

func ElementRefsCount(count int) {
	c := color.New(color.FgGreen).SprintFunc()
	if count == 1 {
		fmt.Println("1", c("element ref"))
	} else {
		fmt.Println(color.HiWhiteString("%d", count), c("element refs"))
	}
}

func AddReference(from, to string) {
	fmt.Println(colorizeAction("Add reference"), colorizeElement(from), referenceDivider, colorizeElement(to))
}

func PathsResult(segments []string) {
	fmt.Println("-", colorizePath(segments))
}
