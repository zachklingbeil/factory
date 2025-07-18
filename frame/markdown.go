package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"

	math "github.com/litao91/goldmark-mathjax"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	h "github.com/yuin/goldmark/renderer/html"
)

var (
	imageRe = regexp.MustCompile(`<img\s+[^>]*alt="(img\+?|img-)"[^>]*>`)
	altRe   = regexp.MustCompile(`alt="([^"]*)"`)
	styleRe = regexp.MustCompile(`style="[^"]*"`)
)

func initGoldmark() *goldmark.Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, math.MathJax),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
			parser.WithBlockParsers(),
			parser.WithInlineParsers(),
		),
		goldmark.WithRendererOptions(
			h.WithHardWraps(),
			h.WithXHTML(),
		),
	)
	return &md
}

func (f *Frame) AddText(file string, elements ...template.HTML) {
	content, err := os.ReadFile(file)
	if err != nil {
		return
	}

	var buf bytes.Buffer
	if err := (*f.Md).Convert(content, &buf); err != nil {
		return
	}

	markdownHTML := template.HTML(buf.String())
	all := append([]template.HTML{markdownHTML}, elements...)

	processed := imageRe.ReplaceAllStringFunc(string(markdownHTML), func(imgTag string) string {
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

	all[0] = template.HTML(processed)
	all = append(all, *f.AddScrollKeybinds()) // Add scroll keybinds at the end

	f.AddFrame(all...) // Only call AddFrame, no return
}

func (f *Frame) AddScrollKeybinds() *template.HTML {
	return f.AddJS(
		`document.addEventListener('keydown', function(event) {
            const c = document.getElementById('frame');
            if (!c) return;
            if (event.key === 'w') {
                c.scrollBy({ top: -100, behavior: 'smooth' });
            }
            if (event.key === 's') {
                c.scrollBy({ top: 100, behavior: 'smooth' });
            }
        });`,
	)
}
