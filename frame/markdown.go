package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gold "github.com/yuin/goldmark/renderer/html"
)

var (
	imageRe = regexp.MustCompile(`<img\s+[^>]*alt="(img\+?|img-)"[^>]*>`)
	altRe   = regexp.MustCompile(`alt="([^"]*)"`)
	styleRe = regexp.MustCompile(`style="[^"]*"`)
)

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, mathjax.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
		),
		goldmark.WithRendererOptions(
			gold.WithHardWraps(),
			gold.WithXHTML(),
		),
	)
	return &md
}

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return ""
	}

	markdownHTML := template.HTML(buf.String())
	all := append([]template.HTML{markdownHTML}, elements...)
	frameHTML := f.CreateFrame("text", all...)

	processed := imageRe.ReplaceAllStringFunc(string(frameHTML), func(imgTag string) string {
		alt := "img"
		if m := altRe.FindStringSubmatch(imgTag); m != nil {
			alt = m[1]
		}
		style := "width:50vw;display:block;margin:0 auto;"
		switch alt {
		case "img+":
			style = "width:75vw;display:block;margin:0 auto;"
		case "img-":
			style = "width:25vw;display:block;margin:0 auto;"
		}
		if styleRe.MatchString(imgTag) {
			return styleRe.ReplaceAllString(imgTag, `style="`+style+`"`)
		}
		return imgTag[:len(imgTag)-1] + ` style="` + style + `">`
	})

	return template.HTML(processed)
}
