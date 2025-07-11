package universe

import (
	"html/template"
	"strings"
)

func (u *Universe) CreateFrame(name string, elements ...template.HTML) {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	u.Frame[name] = &html
}

func (u *Universe) AddCSS(frame string, styles map[string]string) {
	if page, exists := u.Frame[frame]; exists {
		var builder strings.Builder

		// Start with existing HTML
		builder.WriteString(string(*page))
		builder.WriteString("<style>")
		for selector, rules := range styles {
			builder.WriteString(selector)
			builder.WriteString(" { ")
			builder.WriteString(rules)
			builder.WriteString(" }\n")
		}
		builder.WriteString("</style>")

		html := template.HTML(builder.String())
		u.Frame[frame] = &html
	}
}

// Step 2: Add JS to existing page
func (u *Universe) AddJS(frame string, js string) {
	if page, exists := u.Frame[frame]; exists {
		var builder strings.Builder
		// Start with existing HTML
		builder.WriteString(string(*page))
		builder.WriteString("<script>")
		builder.WriteString(js)
		builder.WriteString("</script>")

		html := template.HTML(builder.String())
		u.Frame[frame] = &html
	}
}

func (u *Universe) Render(frame string) template.HTML {
	if page, exists := u.Frame[frame]; exists {
		return *page
	}
	return template.HTML("")
}

func (u *Universe) ListFrames() []string {
	frames := make([]string, 0, len(u.Frame))
	for name := range u.Frame {
		frames = append(frames, name)
	}
	return frames
}

func (u *Universe) DeleteFrame(frame string) {
	delete(u.Frame, frame)
}
