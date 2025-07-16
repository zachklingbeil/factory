package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"
)

var (
	img = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)(?:\s+"([^"]*)")?\)`)
)

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}
	md := string(content)

	md = img.ReplaceAllStringFunc(md, func(match string) string {
		m := img.FindStringSubmatch(match)
		if len(m) >= 3 {
			alt, src := m[1], m[2]
			return string(f.Img(src, alt, "50vw"))
		}
		return match
	})

	var buf bytes.Buffer
	if err := (*f.Md).Convert([]byte(md), &buf); err != nil {
		return template.HTML("")
	}
	wrapped := template.HTML(`<div class="text">` + buf.String() + `</div>`)
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, wrapped)
	allElements = append(allElements, elements...)
	return f.CreateFrame(allElements...)
}
