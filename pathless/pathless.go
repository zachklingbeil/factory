package pathless

import (
	"bytes"
	"html/template"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Pathless struct {
	Font  string
	Color string
	HTML  *template.HTML
}

func InitPathless(color string, body template.HTML) *Pathless {
	pathless := &Pathless{
		Font:  "'Roboto', sans-serif",
		Color: color,
	}
	html := pathless.Zero(body)
	pathless.HTML = &html
	return pathless
}

func (p *Pathless) Zero(body template.HTML) template.HTML {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>hello universe</title>
        <style>
            *,
            *::before,
            *::after {
                box-sizing: border-box;
                margin: 0;
                scrollbar-width: none;
                -ms-overflow-style: none;
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -ms-user-select: none;
            }
            *::-webkit-scrollbar {
                display: none;
            }
            html,
            body {
                color: white;
                background-color: black;
                overflow-y: auto;
                height: 100vh;
                width: 100vw;
                font-family: ` + p.Font + `;
                scroll-behavior: smooth;
                box-sizing: border-box;
                border-radius: 0.3125em;
                display: flex;
                flex-direction: column;
            }
            body {
                border: medium solid ` + p.Color + `;
            }
        </style>
    </head>
    <body>{{.Body}}</body>
</html>`

	t := template.Must(template.New("page").Parse(tmpl))
	var buf bytes.Buffer

	data := struct {
		Body template.HTML
	}{
		Body: body,
	}

	t.Execute(&buf, data)
	return template.HTML(buf.String())
}

func (p *Pathless) CreateFrame(elements ...template.HTML) *template.HTML {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	return &html
}

func (p *Pathless) AddCSS(frame *template.HTML, styles map[string]string) *template.HTML {
	var builder strings.Builder
	builder.WriteString(string(*frame))
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	html := template.HTML(builder.String())
	return &html
}

func (p *Pathless) AddJS(frame *template.HTML, js string) *template.HTML {
	var builder strings.Builder
	builder.WriteString(string(*frame))
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	html := template.HTML(builder.String())
	return &html
}

func (p *Pathless) FromMarkdown(filePath string, elements ...template.HTML) (*template.HTML, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, err
	}

	markdownHTML := template.HTML(buf.String())
	allElements := append([]template.HTML{markdownHTML}, elements...)

	return p.CreateFrame(allElements...), nil
}
