package universe

import (
	"html/template"
	"strings"
)

type Universe struct {
	*Element
	Components []Component
}

func NewUniverse() *Universe {
	return &Universe{
		Components: make([]Component, 0),
		Element:    NewElement(),
	}
}

type Component struct {
	HTML template.HTML
	CSS  template.CSS
	JS   template.JS
}

// ComponentOption is a function that configures a Component
type ComponentOption func(*Component)

// NewComponent creates a Component from HTML elements and optional CSS/JS, and adds it to the Universe.
func (u *Universe) NewComponent(htmls []template.HTML, opts ...ComponentOption) Component {
	component := Component{
		HTML: combineHTML(htmls...),
	}

	for _, opt := range opts {
		opt(&component)
	}

	u.Components = append(u.Components, component)
	return component
}

// combineHTML efficiently aggregates multiple template.HTML elements.
func combineHTML(elements ...template.HTML) template.HTML {
	if len(elements) == 0 {
		return ""
	}

	var builder strings.Builder
	for _, el := range elements {
		builder.WriteString(string(el))
	}
	return template.HTML(builder.String())
}

// buildCSS efficiently builds CSS from a selector-to-rules map.
func buildCSS(styles map[string]string) template.CSS {
	if len(styles) == 0 {
		return ""
	}

	var builder strings.Builder
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	return template.CSS(builder.String())
}

// CSS is an option to add CSS to a Component.
func CSS(styles map[string]string) ComponentOption {
	return func(c *Component) {
		c.CSS = buildCSS(styles)
	}
}

// JS is an option to add JS to a Component.
func JS(js template.JS) ComponentOption {
	return func(c *Component) {
		c.JS = js
	}
}
