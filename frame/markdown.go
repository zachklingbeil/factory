package frame

import (
	"bytes"
	"html"
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
	image = regexp.MustCompile(`<img\s+[^>]*alt="(img\+?|img-)"[^>]*>`)
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

	// Post-process <img> tags in frameHTML, assign class based on alt
	imgRe := regexp.MustCompile(`<img\s+[^>]*alt="(img\+?|img-)"[^>]*>`)
	processed := imgRe.ReplaceAllStringFunc(string(frameHTML), func(imgTag string) string {
		altRe := regexp.MustCompile(`alt="([^"]*)"`)
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
		// Replace or add style attribute
		styleRe := regexp.MustCompile(`style="[^"]*"`)
		if styleRe.MatchString(imgTag) {
			imgTag = styleRe.ReplaceAllString(imgTag, `style="`+style+`"`)
		} else {
			imgTag = imgTag[:len(imgTag)-1] + ` style="` + style + `">`
		}
		return imgTag
	})

	return template.HTML(processed)
}

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
	frameHTML := f.CreateFrame(allElements...)

	// Post-process <img> tags in frameHTML, replace with f.Img
	processed := image.ReplaceAllStringFunc(string(frameHTML), func(imgTag string) string {
		// Extract src and alt
		srcRe := regexp.MustCompile(`src="([^"]*)"`)
		altRe := regexp.MustCompile(`alt="([^"]*)"`)
		src := ""
		alt := ""
		if m := srcRe.FindStringSubmatch(imgTag); m != nil {
			src = html.UnescapeString(m[1])
		}
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
		return string(f.Img(src, alt, width))
	})

	return template.HTML(processed)
}
