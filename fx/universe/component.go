package universe

import (
	"fmt"
	"html/template"
	"strings"
)

type Component struct {
	Name     string
	elements []Element
	CSS      template.CSS
	JS       template.JS
}

// Element interface for composable HTML elements
type Element interface {
	Render() template.HTML
}

// NewComponent creates a new component with tag as name
func NewComponent(name string) *Component {
	return &Component{
		Name:     name,
		elements: make([]Element, 0),
	}
}

// SetCSS converts a map of CSS styles to template.CSS
func (c *Component) SetCSS(styles map[string]string) {
	var builder strings.Builder
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	c.CSS = template.CSS(builder.String())
}

// SetJS converts a string to template.JS
func (c *Component) SetJS(js string) {
	c.JS = template.JS(js)
}

// Add appends an element to the component
func (c *Component) Add(element Element) *Component {
	c.elements = append(c.elements, element)
	return c
}

// AddMultiple appends multiple elements to the component
func (c *Component) AddMultiple(elements ...Element) *Component {
	c.elements = append(c.elements, elements...)
	return c
}

// Render renders all elements in the component
func (c *Component) Render() template.HTML {
	var content strings.Builder

	for _, element := range c.elements {
		content.WriteString(string(element.Render()))
	}

	if c.Name != "" {
		return template.HTML(fmt.Sprintf("<%s>%s</%s>", c.Name, content.String(), c.Name))
	}

	return template.HTML(content.String())
}
