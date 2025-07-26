package zero

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"strings"
)

type Element interface {
	Div(class string) One
	Link(href, text string) One
	List(items []any, ordered bool) One
	Img(src, alt, reference string) One
	Video(src string) One
	Audio(src string) One
	Iframe(src string) One
	Embed(src string) One
	Source(src string) One
	Nav(attrs map[string]string) One
	Canvas(id string) One
	Table(cols uint8, rows uint64, data [][]string) One
	FileToString(path string) string
	CoordinatePlane() One // <-- Add this line

}

// --- element Implementation ---
type element struct{}

func NewElement() Element {
	return &element{}
}

func (e *element) FileToString(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(file)
}

func (e *element) Div(class string) One {
	return One(template.HTML(fmt.Sprintf(`<div class="%s"></div>`, html.EscapeString(class))))
}

func (e *element) Link(href, text string) One {
	return One(template.HTML(fmt.Sprintf(`<a href="%s">%s</a>`, html.EscapeString(href), html.EscapeString(text))))
}

func (e *element) List(items []any, ordered bool) One {
	tag := "ul"
	if ordered {
		tag = "ol"
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("<%s>", tag))
	for _, item := range items {
		b.WriteString(fmt.Sprintf("<li>%v</li>", html.EscapeString(fmt.Sprintf("%v", item))))
	}
	b.WriteString(fmt.Sprintf("</%s>", tag))
	return One(template.HTML(b.String()))
}

func (e *element) Img(src, alt, reference string) One {
	return One(template.HTML(fmt.Sprintf(`<img src="%s" alt="%s" ref="%s"/>`, html.EscapeString(src), html.EscapeString(alt), html.EscapeString(reference))))
}

func (e *element) Video(src string) One {
	return One(template.HTML(fmt.Sprintf(`<video src="%s"></video>`, html.EscapeString(src))))
}

func (e *element) Audio(src string) One {
	return One(template.HTML(fmt.Sprintf(`<audio src="%s"></audio>`, html.EscapeString(src))))
}

func (e *element) Iframe(src string) One {
	return One(template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src))))
}

func (e *element) Embed(src string) One {
	return One(template.HTML(fmt.Sprintf(`<embed src="%s"/>`, html.EscapeString(src))))
}

func (e *element) Source(src string) One {
	return One(template.HTML(fmt.Sprintf(`<source src="%s"/>`, html.EscapeString(src))))
}

func (e *element) Nav(attrs map[string]string) One {
	return ClosedTag("nav", attrs)
}

func (e *element) Canvas(id string) One {
	return One(template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id))))
}

func (e *element) Table(cols uint8, rows uint64, data [][]string) One {
	var b strings.Builder
	b.WriteString("<table>")
	for _, row := range data {
		b.WriteString("<tr>")
		for _, cell := range row {
			b.WriteString(fmt.Sprintf("<td>%s</td>", html.EscapeString(cell)))
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")
	return One(template.HTML(b.String()))
}

func (e *element) CoordinatePlane() One {
	return One(template.HTML(`<div class="coordinate-plane" id="coordinate-plane"></div>`))
}
