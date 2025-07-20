package element

import (
	"fmt"
	"html"
)

// --- Generic tag builders ---
func (z *Zero) Tag(tag, text string) *Element {
	escaped := html.EscapeString(text)
	return NewElement(fmt.Sprintf("<%s>%s</%s>", tag, escaped, tag))
}

func (z *Zero) ClosedTag(tag string, attrs map[string]string) *Element {
	attrStr := ""
	for k, v := range attrs {
		attrStr += fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v))
	}
	return NewElement(fmt.Sprintf("<%s%s>", tag, attrStr))
}

func (z *Zero) Div(value string) *Element {
	return NewElement(fmt.Sprintf(`<div class="%s"></div>`, html.EscapeString(value)))
}

func (z *Zero) JS(js string) *Element {
	script := fmt.Sprintf("<script>%s</script>", js)
	return NewElement(script)
}

func (z *Zero) CSS(css string) *Element {
	style := fmt.Sprintf("<style>%s</style>", css)
	return NewElement(style)
}

// Instance methods for Element
func (z *Zero) H1(s string) *Element        { return z.Tag("h1", s) }
func (z *Zero) H2(s string) *Element        { return z.Tag("h2", s) }
func (z *Zero) H3(s string) *Element        { return z.Tag("h3", s) }
func (z *Zero) H4(s string) *Element        { return z.Tag("h4", s) }
func (z *Zero) H5(s string) *Element        { return z.Tag("h5", s) }
func (z *Zero) H6(s string) *Element        { return z.Tag("h6", s) }
func (z *Zero) Paragraph(s string) *Element { return z.Tag("p", s) }
func (z *Zero) Span(s string) *Element      { return z.Tag("span", s) }
func (z *Zero) Link(href, text string) *Element {
	attrStr := fmt.Sprintf(` href="%s"`, html.EscapeString(href))
	return NewElement(fmt.Sprintf("<a%s>%s</a>", attrStr, html.EscapeString(text)))
}

func (z *Zero) List(keys []any, ordered bool) *Element {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	list := ""
	for _, item := range keys {
		list += string(z.Tag("li", fmt.Sprint(item)).HTML)
	}
	return NewElement(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}

func (z *Zero) Img(src, alt, reference string) *Element {
	img := string(z.ClosedTag("img", map[string]string{"src": src, "alt": alt}).HTML)
	if reference != "" {
		return NewElement(fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(reference), img))
	}
	return NewElement(img)
}

func (z *Zero) Video(src string) *Element {
	return NewElement(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}

func (z *Zero) Audio(src string) *Element {
	return NewElement(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}

func (z *Zero) Iframe(src string) *Element {
	return NewElement(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}

func (z *Zero) Embed(src string) *Element {
	return z.ClosedTag("embed", map[string]string{"src": src})
}

func (z *Zero) Source(src string) *Element {
	return z.ClosedTag("source", map[string]string{"src": src})
}

func (z *Zero) Canvas(id string) *Element {
	return NewElement(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

func (z *Zero) Nav(attrs map[string]string) *Element {
	return z.ClosedTag("nav", attrs)
}

func (z *Zero) Button(label string) *Element {
	return z.Tag("button", label)
}

func (z *Zero) Code(code string) *Element {
	return z.Tag("code", code)
}

func (z *Zero) Table(cols uint8, rows uint64, data [][]string) *Element {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += string(z.Tag("td", cell).HTML)
		}
		table += "</tr>"
	}
	table += "</table>"
	return NewElement(table)
}

func (z *Zero) Strong(s string) *Element  { return z.Tag("strong", s) }
func (z *Zero) Em(s string) *Element      { return z.Tag("em", s) }
func (z *Zero) Small(s string) *Element   { return z.Tag("small", s) }
func (z *Zero) Mark(s string) *Element    { return z.Tag("mark", s) }
func (z *Zero) Del(s string) *Element     { return z.Tag("del", s) }
func (z *Zero) Ins(s string) *Element     { return z.Tag("ins", s) }
func (z *Zero) Sub(s string) *Element     { return z.Tag("sub", s) }
func (z *Zero) Sup(s string) *Element     { return z.Tag("sup", s) }
func (z *Zero) Kbd(s string) *Element     { return z.Tag("kbd", s) }
func (z *Zero) Samp(s string) *Element    { return z.Tag("samp", s) }
func (z *Zero) VarElem(s string) *Element { return z.Tag("var", s) }
func (z *Zero) Abbr(s string) *Element    { return z.Tag("abbr", s) }
func (z *Zero) Time(s string) *Element    { return z.Tag("time", s) }
