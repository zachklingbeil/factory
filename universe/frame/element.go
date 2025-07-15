package frame

import (
	"fmt"
	"html"
	"html/template"
)

// --- Generic tag builders ---
func (f *Frame) Tag(tag, text string) template.HTML {
	escaped := html.EscapeString(text)
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, escaped, tag))
}
func (f *Frame) ClosedTag(tag string, attrs map[string]string) template.HTML {
	attrStr := ""
	for k, v := range attrs {
		attrStr += fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v))
	}
	return template.HTML(fmt.Sprintf("<%s%s>", tag, attrStr))
}

func (f *Frame) H1(s string) template.HTML        { return f.Tag("h1", s) }
func (f *Frame) H2(s string) template.HTML        { return f.Tag("h2", s) }
func (f *Frame) H3(s string) template.HTML        { return f.Tag("h3", s) }
func (f *Frame) H4(s string) template.HTML        { return f.Tag("h4", s) }
func (f *Frame) H5(s string) template.HTML        { return f.Tag("h5", s) }
func (f *Frame) H6(s string) template.HTML        { return f.Tag("h6", s) }
func (f *Frame) Paragraph(s string) template.HTML { return f.Tag("p", s) }
func (f *Frame) Span(s string) template.HTML      { return f.Tag("span", s) }
func (f *Frame) Link(href, text string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text)))
}
func (f *Frame) List(items []any, ordered bool) template.HTML {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	list := ""
	for _, item := range items {
		list += string(f.Tag("li", fmt.Sprint(item)))
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}

func (f *Frame) Img(src, alt, width, height string) template.HTML {
	return f.ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": width, "height": height})
}
func (f *Frame) Video(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func (f *Frame) Audio(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func (f *Frame) Iframe(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func (f *Frame) Embed(src string) template.HTML {
	return f.ClosedTag("embed", map[string]string{"src": src})
}
func (f *Frame) Source(src string) template.HTML {
	return f.ClosedTag("source", map[string]string{"src": src})
}
func (f *Frame) Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

func (f *Frame) Nav(attrs map[string]string) template.HTML {
	return f.ClosedTag("nav", attrs)
}

func (f *Frame) Button(label string) template.HTML {
	return f.Tag("button", label)
}
func (f *Frame) Code(code string) template.HTML {
	return f.Tag("code", code)
}
func (f *Frame) Table(cols uint8, rows uint64, data [][]string) template.HTML {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += string(f.Tag("td", cell))
		}
		table += "</tr>"
	}
	table += "</table>"
	return template.HTML(table)
}

func (f *Frame) Strong(s string) template.HTML { return f.Tag("strong", s) }
func (f *Frame) Em(s string) template.HTML     { return f.Tag("em", s) }
func (f *Frame) Small(s string) template.HTML  { return f.Tag("small", s) }
func (f *Frame) Mark(s string) template.HTML   { return f.Tag("mark", s) }
func (f *Frame) Del(s string) template.HTML    { return f.Tag("del", s) }
func (f *Frame) Ins(s string) template.HTML    { return f.Tag("ins", s) }
func (f *Frame) Sub(s string) template.HTML    { return f.Tag("sub", s) }
func (f *Frame) Sup(s string) template.HTML    { return f.Tag("sup", s) }
func (f *Frame) Kbd(s string) template.HTML    { return f.Tag("kbd", s) }
func (f *Frame) Samp(s string) template.HTML   { return f.Tag("samp", s) }
func (f *Frame) Var(s string) template.HTML    { return f.Tag("var", s) }
func (f *Frame) Abbr(s string) template.HTML   { return f.Tag("abbr", s) }
func (f *Frame) Time(s string) template.HTML   { return f.Tag("time", s) }
