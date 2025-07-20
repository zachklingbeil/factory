package frame

import (
	"fmt"
	"html/template"
	"strings"
)

func (f *Frame) AddCSS(styles map[string]string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	html := template.HTML(builder.String())
	return &html
}

func (f *Frame) AddJS(js string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	html := template.HTML(builder.String())
	return &html
}

func intToWord(n int) (string, error) {
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}
	if n >= 0 && n < len(words) {
		return words[n], nil
	}
	return "", fmt.Errorf("intToWord: more words needed for index %d", n)
}
