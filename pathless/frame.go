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

func initGoldmark() *goldmark.Markdown {
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
	return &md
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
