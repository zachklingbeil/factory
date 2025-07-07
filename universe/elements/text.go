package elements

import (
	"fmt"
	"html"
	"html/template"
)

type text struct{}

func (t *text) H1(s string) template.HTML        { return Tag("h1", s) }
func (t *text) H2(s string) template.HTML        { return Tag("h2", s) }
func (t *text) H3(s string) template.HTML        { return Tag("h3", s) }
func (t *text) H4(s string) template.HTML        { return Tag("h4", s) }
func (t *text) H5(s string) template.HTML        { return Tag("h5", s) }
func (t *text) H6(s string) template.HTML        { return Tag("h6", s) }
func (t *text) Paragraph(s string) template.HTML { return Tag("p", s) }
func (t *text) Span(s string) template.HTML      { return Tag("span", s) }
func (t *text) Link(href, text string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text)))
}
func (t *text) List(items []any, ordered bool) template.HTML {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	list := ""
	for _, item := range items {
		list += string(Tag("li", fmt.Sprint(item)))
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}
