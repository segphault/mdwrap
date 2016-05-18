package main

import (
	"fmt"
	"github.com/golang-commonmark/markdown"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

var safeLineBreakPattern *regexp.Regexp = regexp.MustCompile("(\\[[^\\]]+\\]\\[[^\\]]*\\][^\\s]?)|(\\[[^\\]]+\\]\\([^\\]]*\\)[^\\s]?)|(`[^`]+`[^\\s]?)|([^\\s]+)")

func findLinesToWrap(data []byte) []int {
	output := []int{}

	md := markdown.New(markdown.HTML(true), markdown.Tables(true))

	for _, token := range md.Parse(data) {
		if p, ok := token.(*markdown.ParagraphOpen); ok && p.Lvl < 1 {
			output = append(output, p.Map[0])
		}
	}

	return output
}

func wrapMarkdownLine(text string, width int) string {
	output := []string{}
	current := ""

	for _, word := range safeLineBreakPattern.FindAllString(text, -1) {
		if len(current+word) < width {
			if current != "" {
				current += " "
			}
			current += word
		} else {
			output = append(output, current)
			current = word
		}
	}

	return strings.Join(append(output, current), "\n")
}

func main() {
	data, _ := ioutil.ReadAll(os.Stdin)
	output := findLinesToWrap(data)
	next := 0

	for line, text := range strings.Split(string(data), "\n") {
		if len(output) > next && line == output[next] {
			next++
			fmt.Println(wrapMarkdownLine(text, 80))
		} else {
			fmt.Println(text)
		}
	}
}
