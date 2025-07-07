package elements

import (
	"fmt"
	"html"
	"html/template"
)

type Elements struct {
	Head template.HTML
	text
	media
	other
	format
}

func NewElements() *Elements {
	return &Elements{
		Head:   Head(),
		text:   text{},
		media:  media{},
		other:  other{},
		format: format{},
	}
}

// --- Generic tag builders ---
func Tag(tag, text string) template.HTML {
	escaped := html.EscapeString(text)
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, escaped, tag))
}
func ClosedTag(tag string, attrs map[string]string) template.HTML {
	attrStr := ""
	for k, v := range attrs {
		attrStr += fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v))
	}
	return template.HTML(fmt.Sprintf("<%s%s>", tag, attrStr))
}
