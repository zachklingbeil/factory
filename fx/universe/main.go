package universe

import (
	"html/template"
	"strings"
)

type Universe struct {
	Element
	Components map[string]*Component
}

func NewUniverse() *Universe {
	return &Universe{
		Components: make(map[string]*Component),
	}
}

// Add adds a component to the universe with a key
func (u *Universe) Add(key string, component *Component) *Universe {
	u.Components[key] = component
	return u
}

// Render renders the entire universe as HTML
func (u *Universe) Render() template.HTML {
	var content strings.Builder

	for _, component := range u.Components {
		content.WriteString(string(component.Render()))
	}

	return template.HTML(content.String())
}

// GetCSS returns all CSS from all components
func (u *Universe) GetCSS() template.CSS {
	var builder strings.Builder

	for _, component := range u.Components {
		if component.CSS != "" {
			builder.WriteString(string(component.CSS))
			builder.WriteString("\n")
		}
	}
	return template.CSS(builder.String())
}

// GetJS returns all JavaScript from all components
func (u *Universe) GetJS() template.JS {
	var builder strings.Builder

	for _, component := range u.Components {
		if component.JS != "" {
			builder.WriteString(string(component.JS))
			builder.WriteString("\n")
		}
	}

	return template.JS(builder.String())
}
