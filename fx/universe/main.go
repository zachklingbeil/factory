package universe

import (
	"html/template"
	"strings"
)

type Universe struct {
	Frame map[string]*template.HTML
}

func NewUniverse() *Universe {
	return &Universe{
		Frame: make(map[string]*template.HTML),
	}
}

func (u *Universe) CreateFrame(name string, elements ...template.HTML) {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	html := template.HTML(builder.String())
	u.Frame[name] = &html
}

func (u *Universe) AddCSS(pageName string, styles map[string]string) {
	if page, exists := u.Frame[pageName]; exists {
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
		u.Frame[pageName] = &html
	}
}

// Step 2: Add JS to existing page
func (u *Universe) AddJS(pageName string, js string) {
	if page, exists := u.Frame[pageName]; exists {
		var builder strings.Builder
		// Start with existing HTML
		builder.WriteString(string(*page))
		builder.WriteString("<script>")
		builder.WriteString(js)
		builder.WriteString("</script>")

		html := template.HTML(builder.String())
		u.Frame[pageName] = &html
	}
}

func (u *Universe) Render(pageName string) template.HTML {
	if page, exists := u.Frame[pageName]; exists {
		return *page
	}
	return template.HTML("")
}

func (u *Universe) ListPages() []string {
	pages := make([]string, 0, len(u.Frame))
	for name := range u.Frame {
		pages = append(pages, name)
	}
	return pages
}

func (u *Universe) DeletePage(pageName string) {
	delete(u.Frame, pageName)
}
