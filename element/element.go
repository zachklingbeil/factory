package element

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type Element struct {
	HTML template.HTML
}

type One interface {
	Render() template.HTML
}

func (e *Element) Render() template.HTML {
	return e.HTML
}
func NewElement(htmlStr string) *Element {
	return &Element{HTML: template.HTML(htmlStr)}
}

func (e *Element) BuildFrame(keys []One) template.HTML {
	body := simplify(keys)
	return template.HTML(body)
}

func (e *Element) WrapFrame(value string, keys []One) template.HTML {
	body := e.BuildFrame(keys)
	frame := fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(value), body)
	return template.HTML(frame)
}

func simplify(keys []One) string {
	var stylesBuilder, scriptsBuilder, htmlBuilder strings.Builder

	for _, item := range keys {
		s := string(item.Render())
		for {
			start, end := strings.Index(s, "<style>"), strings.Index(s, "</style>")
			if start < 0 || end <= start {
				break
			}
			stylesBuilder.WriteString(s[start+len("<style>") : end])
			s = s[:start] + s[end+len("</style>"):]
		}
		for {
			start, end := strings.Index(s, "<script>"), strings.Index(s, "</script>")
			if start < 0 || end <= start {
				break
			}
			scriptsBuilder.WriteString(s[start+len("<script>") : end])
			s = s[:start] + s[end+len("</script>"):]
		}
		htmlBuilder.WriteString(s)
	}

	var resultBuilder strings.Builder
	if stylesBuilder.Len() > 0 {
		resultBuilder.WriteString("<style>")
		resultBuilder.WriteString(stylesBuilder.String())
		resultBuilder.WriteString("</style>")
	}
	if scriptsBuilder.Len() > 0 {
		resultBuilder.WriteString("<script>")
		resultBuilder.WriteString(scriptsBuilder.String())
		resultBuilder.WriteString("</script>")
	}
	resultBuilder.WriteString(htmlBuilder.String())
	return resultBuilder.String()
}

// --- Generic tag builders ---
func Tag(tag, text string) string {
	escaped := html.EscapeString(text)
	return fmt.Sprintf("<%s>%s</%s>", tag, escaped, tag)
}
func ClosedTag(tag string, attrs map[string]string) string {
	attrStr := ""
	for k, v := range attrs {
		attrStr += fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v))
	}
	return fmt.Sprintf("<%s%s>", tag, attrStr)
}
func (e *Element) Div(value string) *Element {
	return NewElement(fmt.Sprintf(`<div class="%s"></div>`, html.EscapeString(value)))
}

func (e *Element) JS(js string) *Element {
	script := fmt.Sprintf("<script>%s</script>", js)
	return NewElement(string(e.HTML) + script)
}

func (e *Element) CSS(css string) *Element {
	style := fmt.Sprintf("<style>%s</style>", css)
	return NewElement(string(e.HTML) + style)
}

// Instance methods for Element
func (e *Element) H1(s string) *Element        { return NewElement(Tag("h1", s)) }
func (e *Element) H2(s string) *Element        { return NewElement(Tag("h2", s)) }
func (e *Element) H3(s string) *Element        { return NewElement(Tag("h3", s)) }
func (e *Element) H4(s string) *Element        { return NewElement(Tag("h4", s)) }
func (e *Element) H5(s string) *Element        { return NewElement(Tag("h5", s)) }
func (e *Element) H6(s string) *Element        { return NewElement(Tag("h6", s)) }
func (e *Element) Paragraph(s string) *Element { return NewElement(Tag("p", s)) }
func (e *Element) Span(s string) *Element      { return NewElement(Tag("span", s)) }
func (e *Element) Link(href, text string) *Element {
	return NewElement(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text)))
}
func (e *Element) List(keys []any, ordered bool) *Element {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	list := ""
	for _, item := range keys {
		list += Tag("li", fmt.Sprint(item))
	}
	return NewElement(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}
func (e *Element) Img(src, alt, reference string) *Element {
	img := ClosedTag("img", map[string]string{"src": src, "alt": alt})
	if reference != "" {
		return NewElement(`<div class="` + reference + `">` + img + `</div>`)
	}
	return NewElement(img)
}
func (e *Element) Video(src string) *Element {
	return NewElement(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func (e *Element) Audio(src string) *Element {
	return NewElement(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func (e *Element) Iframe(src string) *Element {
	return NewElement(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func (e *Element) Embed(src string) *Element {
	return NewElement(ClosedTag("embed", map[string]string{"src": src}))
}
func (e *Element) Source(src string) *Element {
	return NewElement(ClosedTag("source", map[string]string{"src": src}))
}
func (e *Element) Canvas(id string) *Element {
	return NewElement(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}
func (e *Element) Nav(attrs map[string]string) *Element {
	return NewElement(ClosedTag("nav", attrs))
}
func (e *Element) Button(label string) *Element {
	return NewElement(Tag("button", label))
}
func (e *Element) Code(code string) *Element {
	return NewElement(Tag("code", code))
}
func (e *Element) Table(cols uint8, rows uint64, data [][]string) *Element {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += Tag("td", cell)
		}
		table += "</tr>"
	}
	table += "</table>"
	return NewElement(table)
}
func (e *Element) Strong(s string) *Element  { return NewElement(Tag("strong", s)) }
func (e *Element) Em(s string) *Element      { return NewElement(Tag("em", s)) }
func (e *Element) Small(s string) *Element   { return NewElement(Tag("small", s)) }
func (e *Element) Mark(s string) *Element    { return NewElement(Tag("mark", s)) }
func (e *Element) Del(s string) *Element     { return NewElement(Tag("del", s)) }
func (e *Element) Ins(s string) *Element     { return NewElement(Tag("ins", s)) }
func (e *Element) Sub(s string) *Element     { return NewElement(Tag("sub", s)) }
func (e *Element) Sup(s string) *Element     { return NewElement(Tag("sup", s)) }
func (e *Element) Kbd(s string) *Element     { return NewElement(Tag("kbd", s)) }
func (e *Element) Samp(s string) *Element    { return NewElement(Tag("samp", s)) }
func (e *Element) VarElem(s string) *Element { return NewElement(Tag("var", s)) }
func (e *Element) Abbr(s string) *Element    { return NewElement(Tag("abbr", s)) }
func (e *Element) Time(s string) *Element    { return NewElement(Tag("time", s)) }
