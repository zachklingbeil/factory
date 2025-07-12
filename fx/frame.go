package fx

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

func (u *Universe) CreateFrame(name string, elements ...template.HTML) {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	u.Frame[name] = &html
}

func (u *Universe) AddCSS(frame string, styles map[string]string) {
	if page, exists := u.Frame[frame]; exists {
		var builder strings.Builder

		// Start with existing HTML
		builder.WriteString(string(*page))
		builder.WriteString("<style>")
		for selector, rules := range styles {
			builder.WriteString(selector)
			builder.WriteString(" { ")
			builder.WriteString(rules)
			builder.WriteString(" }\n")
		}
		builder.WriteString("</style>")

		html := template.HTML(builder.String())
		u.Frame[frame] = &html
	}
}

// Step 2: Add JS to existing page
func (u *Universe) AddJS(frame string, js string) {
	if page, exists := u.Frame[frame]; exists {
		var builder strings.Builder
		// Start with existing HTML
		builder.WriteString(string(*page))
		builder.WriteString("<script>")
		builder.WriteString(js)
		builder.WriteString("</script>")

		html := template.HTML(builder.String())
		u.Frame[frame] = &html
	}
}

func (u *Universe) Render(frame string) template.HTML {
	if page, exists := u.Frame[frame]; exists {
		return *page
	}
	return template.HTML("")
}

func (u *Universe) ListFrames() []string {
	frames := make([]string, 0, len(u.Frame))
	for name := range u.Frame {
		frames = append(frames, name)
	}
	return frames
}

func (u *Universe) DeleteFrame(frame string) {
	delete(u.Frame, frame)
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

func H1(s string) template.HTML        { return Tag("h1", s) }
func H2(s string) template.HTML        { return Tag("h2", s) }
func H3(s string) template.HTML        { return Tag("h3", s) }
func H4(s string) template.HTML        { return Tag("h4", s) }
func H5(s string) template.HTML        { return Tag("h5", s) }
func H6(s string) template.HTML        { return Tag("h6", s) }
func Paragraph(s string) template.HTML { return Tag("p", s) }
func Span(s string) template.HTML      { return Tag("span", s) }
func Link(href, text string) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text)))
}
func List(items []any, ordered bool) template.HTML {
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

func Img(src, alt, width, height string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": width, "height": height})
}
func Video(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func Audio(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func Iframe(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func Embed(src string) template.HTML {
	return ClosedTag("embed", map[string]string{"src": src})
}
func Source(src string) template.HTML {
	return ClosedTag("source", map[string]string{"src": src})
}
func Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

func Nav(attrs map[string]string) template.HTML {
	return ClosedTag("nav", attrs)
}

func Button(label string) template.HTML {
	return Tag("button", label)
}
func Code(code string) template.HTML {
	return Tag("code", code)
}
func Table(cols uint8, rows uint64, data [][]string) template.HTML {
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

func Strong(s string) template.HTML { return Tag("strong", s) }
func Em(s string) template.HTML     { return Tag("em", s) }
func Small(s string) template.HTML  { return Tag("small", s) }
func Mark(s string) template.HTML   { return Tag("mark", s) }
func Del(s string) template.HTML    { return Tag("del", s) }
func Ins(s string) template.HTML    { return Tag("ins", s) }
func Sub(s string) template.HTML    { return Tag("sub", s) }
func Sup(s string) template.HTML    { return Tag("sup", s) }
func Kbd(s string) template.HTML    { return Tag("kbd", s) }
func Samp(s string) template.HTML   { return Tag("samp", s) }
func Var(s string) template.HTML    { return Tag("var", s) }
func Abbr(s string) template.HTML   { return Tag("abbr", s) }
func Time(s string) template.HTML   { return Tag("time", s) }
