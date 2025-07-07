package element

import (
	"fmt"
	"html"
	"html/template"
)

type Element struct {
	text
	media
	other
	format
}

type Component struct {
	HTML template.HTML `json:"html"`
	CSS  template.CSS  `json:"css"`
}

func NewElements() *Element {
	return &Element{
		text:   text{},
		media:  media{},
		other:  other{},
		format: format{},
	}
}

func (e *Element) NewComponent(html template.HTML, css template.CSS) *Component {
	return &Component{
		HTML: html,
		CSS:  css,
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

type media struct{}

func (m *media) Img(src, alt string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": "100%", "height": "auto"})
}
func (m *media) Video(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func (m *media) Audio(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func (m *media) Iframe(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func (m *media) Embed(src string) template.HTML {
	return ClosedTag("embed", map[string]string{"src": src})
}
func (m *media) Source(src string) template.HTML {
	return ClosedTag("source", map[string]string{"src": src})
}
func (m *media) Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

type other struct{}

func (o *other) Nav(attrs map[string]string) template.HTML {
	return ClosedTag("nav", attrs)
}

func (o *other) Button(label string) template.HTML {
	return Tag("button", label)
}
func (o *other) Code(code string) template.HTML {
	return Tag("code", code)
}
func (o *other) Table(cols uint8, rows uint64, data [][]string) template.HTML {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += string(Tag("td", cell))
		}
		table += "</tr>"
	}
	table += "</table>"
	return template.HTML(table)
}

type format struct{}

func (f *format) Strong(s string) template.HTML { return Tag("strong", s) }
func (f *format) Em(s string) template.HTML     { return Tag("em", s) }
func (f *format) Small(s string) template.HTML  { return Tag("small", s) }
func (f *format) Mark(s string) template.HTML   { return Tag("mark", s) }
func (f *format) Del(s string) template.HTML    { return Tag("del", s) }
func (f *format) Ins(s string) template.HTML    { return Tag("ins", s) }
func (f *format) Sub(s string) template.HTML    { return Tag("sub", s) }
func (f *format) Sup(s string) template.HTML    { return Tag("sup", s) }
func (f *format) Kbd(s string) template.HTML    { return Tag("kbd", s) }
func (f *format) Samp(s string) template.HTML   { return Tag("samp", s) }
func (f *format) Var(s string) template.HTML    { return Tag("var", s) }
func (f *format) Abbr(s string) template.HTML   { return Tag("abbr", s) }
func (f *format) Time(s string) template.HTML   { return Tag("time", s) }
