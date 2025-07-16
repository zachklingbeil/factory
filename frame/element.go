package frame

import (
	"fmt"
	"html"
	"html/template"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	yhtml "github.com/yuin/goldmark/renderer/html"
)

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

func (f *Frame) H1(s string) template.HTML        { return Tag("h1", s) }
func (f *Frame) H2(s string) template.HTML        { return Tag("h2", s) }
func (f *Frame) H3(s string) template.HTML        { return Tag("h3", s) }
func (f *Frame) H4(s string) template.HTML        { return Tag("h4", s) }
func (f *Frame) H5(s string) template.HTML        { return Tag("h5", s) }
func (f *Frame) H6(s string) template.HTML        { return Tag("h6", s) }
func (f *Frame) Paragraph(s string) template.HTML { return Tag("p", s) }
func (f *Frame) Span(s string) template.HTML      { return Tag("span", s) }
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
		list += string(Tag("li", fmt.Sprint(item)))
	}
	return template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, list, tag))
}

func (f *Frame) Img(src, alt, width, height string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": width, "height": height})
}
func (f *Frame) Img2(src, alt string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": "75vw", "height": "auto"})
}
func (f *Frame) Img3(src, alt string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": "50vw", "height": "auto"})
}
func (f *Frame) Img4(src, alt string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": "25vw", "height": "auto"})
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
	return ClosedTag("embed", map[string]string{"src": src})
}
func (f *Frame) Source(src string) template.HTML {
	return ClosedTag("source", map[string]string{"src": src})
}
func (f *Frame) Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}

func (f *Frame) Nav(attrs map[string]string) template.HTML {
	return ClosedTag("nav", attrs)
}

func (f *Frame) Button(label string) template.HTML {
	return Tag("button", label)
}
func (f *Frame) Code(code string) template.HTML {
	return Tag("code", code)
}
func (f *Frame) Table(cols uint8, rows uint64, data [][]string) template.HTML {
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

func (f *Frame) Strong(s string) template.HTML { return Tag("strong", s) }
func (f *Frame) Em(s string) template.HTML     { return Tag("em", s) }
func (f *Frame) Small(s string) template.HTML  { return Tag("small", s) }
func (f *Frame) Mark(s string) template.HTML   { return Tag("mark", s) }
func (f *Frame) Del(s string) template.HTML    { return Tag("del", s) }
func (f *Frame) Ins(s string) template.HTML    { return Tag("ins", s) }
func (f *Frame) Sub(s string) template.HTML    { return Tag("sub", s) }
func (f *Frame) Sup(s string) template.HTML    { return Tag("sup", s) }
func (f *Frame) Kbd(s string) template.HTML    { return Tag("kbd", s) }
func (f *Frame) Samp(s string) template.HTML   { return Tag("samp", s) }
func (f *Frame) Var(s string) template.HTML    { return Tag("var", s) }
func (f *Frame) Abbr(s string) template.HTML   { return Tag("abbr", s) }
func (f *Frame) Time(s string) template.HTML   { return Tag("time", s) }

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, mathjax.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			yhtml.WithHardWraps(),
			// yhtml.WithXHTML(),
		),
	)
	return &md
}
