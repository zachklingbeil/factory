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
	"github.com/yuin/goldmark/renderer/html"
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
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	return &md
}

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return template.HTML("")
	}
	wrapped := template.HTML(`<div class="text">` + buf.String() + `</div>`)
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, wrapped)
	allElements = append(allElements, elements...)
	frameHTML := f.CreateFrame(allElements...)

	// Post-process <img> tags in frameHTML, set width directly
	imgRe := regexp.MustCompile(`<img\s+([^>]*alt="(img\+?|img-)"[^>]*)>`)
	processed := imgRe.ReplaceAllStringFunc(string(frameHTML), func(imgTag string) string {
		altRe := regexp.MustCompile(`alt="([^"]*)"`)
		alt := "img"
		if m := altRe.FindStringSubmatch(imgTag); m != nil {
			alt = m[1]
		}
		width := "50vw"
		switch alt {
		case "img+":
			width = "75vw"
		case "img-":
			width = "25vw"
		}
		// Replace or add width attribute
		widthRe := regexp.MustCompile(`width="[^"]*"`)
		if widthRe.MatchString(imgTag) {
			imgTag = widthRe.ReplaceAllString(imgTag, `width="`+width+`"`)
		} else {
			imgTag = imgTag[:len(imgTag)-1] + ` width="` + width + `" height="auto">`
		}
		return imgTag
	})

	return template.HTML(processed)
}

// ...existing code...
func (f *Frame) FromMarkdown2(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return template.HTML("")
	}
	wrapped := template.HTML(`<div class="text">` + buf.String() + `</div>`)
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, wrapped)
	allElements = append(allElements, elements...)
	return f.CreateFrame(allElements...)
}
