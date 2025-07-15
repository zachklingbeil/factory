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
	Md    *goldmark.Markdown
}

func InitPathless(color string, body template.HTML) *Pathless {
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
	pathless := &Pathless{
		Font:  "'Roboto', sans-serif",
		Color: color,
		Md:    &md,
	}
	pathless.Zero(body)
	return pathless
}

func (p *Pathless) Zero(body template.HTML) {
	templateStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>hello universe</title>
    <style>
        *, *::before, *::after {
            box-sizing: border-box;
            margin: 0;
            scrollbar-width: none;
            -ms-overflow-style: none;
            user-select: none;
        }
        *::-webkit-scrollbar { display: none; }
        html, body {
            color: white;
            background-color: black;
            height: 100vh;
            width: 100vw;
            font-family: ` + p.Font + `;
            scroll-behavior: smooth;
            overflow: hidden;
        }
        body {
            border: medium solid ` + p.Color + `; 
            border-radius: 0.3125em;
            display: flex;
        }
        main {
            flex: 1;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            overflow-y: auto;
        }
    </style>
</head>
<body><main>{{.Body}}</main></body>
</html>`

	tmpl := template.Must(template.New("page").Parse(templateStr))
	var buf bytes.Buffer

	data := struct{ Body template.HTML }{Body: body}
	if err := tmpl.Execute(&buf, data); err != nil {
		result := template.HTML(strings.ReplaceAll(templateStr, "{{.Body}}", string(body)))
		p.HTML = &result
		return
	}
	result := template.HTML(buf.String())
	p.HTML = &result
}

func (p *Pathless) CreateFrame(elements ...template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}

	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}

	return template.HTML(builder.String())
}

func (p *Pathless) AddCSS(frame template.HTML, styles map[string]string) template.HTML {
	if len(styles) == 0 {
		return frame
	}

	var builder strings.Builder
	builder.WriteString(string(frame))
	builder.WriteString("<style>")

	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}

	builder.WriteString("</style>")
	return template.HTML(builder.String())
}

func (p *Pathless) AddJS(frame template.HTML, js string) template.HTML {
	if js == "" {
		return frame
	}

	var builder strings.Builder
	builder.WriteString(string(frame))
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")

	return template.HTML(builder.String())
}

func (p *Pathless) FromMarkdown(filePath string, elements ...template.HTML) (template.HTML, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return template.HTML(""), err
	}

	var buf bytes.Buffer
	if err := (*p.Md).Convert(content, &buf); err != nil {
		return template.HTML(""), err
	}
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, template.HTML(buf.String()))
	allElements = append(allElements, elements...)
	return p.CreateFrame(allElements...), nil
}
