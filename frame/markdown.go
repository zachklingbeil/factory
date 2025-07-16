package frame

import (
	"bytes"
	"html/template"
	"os"
	"regexp"
)

var (
	img = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]+)\)`)
)

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
	return f.CreateFrame(allElements...)
}

func (f *Frame) FromMarkdown2(file string, elements ...template.HTML) template.HTML {
	content, err := os.ReadFile(file)
	if err != nil {
		return template.HTML("")
	}
	md := string(content)

	parts := img.FindAllStringIndex(md, -1)
	result := make([]template.HTML, 0, len(parts)+len(elements)+1)
	last := 0

	for _, idx := range parts {
		// Convert text before the image
		if idx[0] > last {
			section := md[last:idx[0]]
			var buf bytes.Buffer
			if err := (*f.Md).Convert([]byte(section), &buf); err == nil {
				result = append(result, template.HTML(buf.String()))
			}
		}
		// Handle the image itself
		imgMatch := md[idx[0]:idx[1]]
		m := img.FindStringSubmatch(imgMatch)
		if len(m) == 3 {
			alt, src := m[1], m[2]
			result = append(result, f.Img(src, alt, "50vw"))
		}
		last = idx[1]
	}

	// Convert any remaining text after the last image
	if last < len(md) {
		section := md[last:]
		var buf bytes.Buffer
		if err := (*f.Md).Convert([]byte(section), &buf); err == nil {
			result = append(result, template.HTML(buf.String()))
		}
	}

	// Add any extra elements
	result = append(result, elements...)

	// Join all sections and wrap in .text
	final := `<div class="text">`
	for _, r := range result {
		final += string(r)
	}
	final += `</div>`

	return f.CreateFrame(template.HTML(final))
}
