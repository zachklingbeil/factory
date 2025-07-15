package pathless

import (
	"fmt"
	"html"
	"html/template"
)

// --- Generic tag builders ---
func (p *Pathless) Tag(tag, text string) template.HTML {
	escaped := html.EscapeString(text)
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, escaped, tag))
}
func (p *Pathless) ClosedTag(tag string, attrs map[string]string) template.HTML {
	attrStr := ""
	for k, v := range attrs {
		attrStr += fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v))
	}
	return template.HTML(fmt.Sprintf("<%s%s>", tag, attrStr))
}

func (p *Pathless) H1(s string) template.HTML        { return p.Tag("h1", s) }
func (p *Pathless) H2(s string) template.HTML        { return p.Tag("h2", s) }
func (p *Pathless) H3(s string) template.HTML        { return p.Tag("h3", s) }
func (p *Pathless) H4(s string) template.HTML        { return p.Tag("h4", s) }
func (p *Pathless) H5(s string) template.HTML        { return p.Tag("h5", s) }
func (p *Pathless) H6(s string) template.HTML        { return p.Tag("h6", s) }
func (p *Pathless) Paragraph(s string) template.HTML { return p.Tag("p", s) }
func (p *Pathless) Span(s string) template.HTML      { return p.Tag("span", s) }
func (p *Pathless) Link(href, text string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text)))
}
func (p *Pathless) List(items []any, ordered bool) template.HTML {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	list := ""
	for _, item := range items {
		list += string(p.Tag("li", fmt.Sprint(item)))
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}

func (p *Pathless) Img(src, alt, width, height string) template.HTML {
	return p.ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": width, "height": height})
}
func (p *Pathless) Video(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func (p *Pathless) Audio(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func (p *Pathless) Iframe(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func (p *Pathless) Embed(src string) template.HTML {
	return p.ClosedTag("embed", map[string]string{"src": src})
}
func (p *Pathless) Source(src string) template.HTML {
	return p.ClosedTag("source", map[string]string{"src": src})
}
func (p *Pathless) Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

func (p *Pathless) Nav(attrs map[string]string) template.HTML {
	return p.ClosedTag("nav", attrs)
}

func (p *Pathless) Button(label string) template.HTML {
	return p.Tag("button", label)
}
func (p *Pathless) Code(code string) template.HTML {
	return p.Tag("code", code)
}
func (p *Pathless) Table(cols uint8, rows uint64, data [][]string) template.HTML {
	table := "<table>"
	for _, row := range data {
		table += "<tr>"
		for _, cell := range row {
			table += string(p.Tag("td", cell))
		}
		table += "</tr>"
	}
	table += "</table>"
	return template.HTML(table)
}

func (p *Pathless) Strong(s string) template.HTML { return p.Tag("strong", s) }
func (p *Pathless) Em(s string) template.HTML     { return p.Tag("em", s) }
func (p *Pathless) Small(s string) template.HTML  { return p.Tag("small", s) }
func (p *Pathless) Mark(s string) template.HTML   { return p.Tag("mark", s) }
func (p *Pathless) Del(s string) template.HTML    { return p.Tag("del", s) }
func (p *Pathless) Ins(s string) template.HTML    { return p.Tag("ins", s) }
func (p *Pathless) Sub(s string) template.HTML    { return p.Tag("sub", s) }
func (p *Pathless) Sup(s string) template.HTML    { return p.Tag("sup", s) }
func (p *Pathless) Kbd(s string) template.HTML    { return p.Tag("kbd", s) }
func (p *Pathless) Samp(s string) template.HTML   { return p.Tag("samp", s) }
func (p *Pathless) Var(s string) template.HTML    { return p.Tag("var", s) }
func (p *Pathless) Abbr(s string) template.HTML   { return p.Tag("abbr", s) }
func (p *Pathless) Time(s string) template.HTML   { return p.Tag("time", s) }
