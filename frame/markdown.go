package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"
)

var (
	imgTag = regexp.MustCompile(`<img\s+[^>]*alt="([^"]*)"[^>]*src="([^"]+)"[^>]*>`)
)

func (f *Frame) FromMarkdown(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}
	md := string(content)

	// First, convert markdown to HTML
	var buf bytes.Buffer
	if err := (*f.Md).Convert([]byte(md), &buf); err != nil {
		return template.HTML("")
	}
	htmlStr := buf.String()

	// Then, process <img> tags for sizing
	htmlStr = imgTag.ReplaceAllStringFunc(htmlStr, func(match string) string {
		m := imgTag.FindStringSubmatch(match)
		if len(m) == 3 {
			alt, src := m[1], m[2]
			// Always use 50vw for width as per your requirement
			return string(f.Img(src, alt, "50vw"))
		}
		return match
	})

	wrapped := template.HTML(`<div class="text">` + htmlStr + `</div>`)
	allElements := make([]template.HTML, 0, len(elements)+1)
	allElements = append(allElements, wrapped)
	allElements = append(allElements, elements...)
	return f.CreateFrame(allElements...)
}
