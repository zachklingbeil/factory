package frame

import (
	"html/template"
	"strings"

	"github.com/yuin/goldmark"
)

type Frame struct {
	Md        *goldmark.Markdown
	Keybinds  map[string]string
	Templates []*template.HTML
}

func NewFrame() *Frame {
	return &Frame{
		Md:       initGoldmark(),
		Keybinds: make(map[string]string),
	}
}

// Add or update a keybind
func (f *Frame) SetKeybind(key, action string) {
	f.Keybinds[key] = action
}

// Retrieve an action for a keybind
func (f *Frame) GetKeybind(key string) (string, bool) {
	action, exists := f.Keybinds[key]
	return action, exists
}

// Remove a keybind
func (f *Frame) RemoveKeybind(key string) {
	delete(f.Keybinds, key)
}

func (f *Frame) CreateFrame(reference string, elements ...template.HTML) template.HTML {
	if len(elements) == 0 {
		return template.HTML("")
	}
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	result := builder.String()
	if reference != "" {
		result = `<div class="` + reference + `">` + result + `</div>`
	}
	return template.HTML(result)
}

func (f *Frame) AddCSS(styles map[string]string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	return template.HTML(builder.String())
}

func (f *Frame) AddJS(js string) template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	return template.HTML(builder.String())
}

// AddKeybindJS generates and injects JavaScript for all keybinds in the Frame's Keybinds map.
func (f *Frame) AddKeybindJS() template.HTML {
	var builder strings.Builder
	builder.WriteString("document.addEventListener('keydown',function(event){")
	builder.WriteString("if(event.target.tagName==='INPUT'||event.target.tagName==='TEXTAREA'){return;}")
	for key, action := range f.Keybinds {
		builder.WriteString("if(event.key==='")
		builder.WriteString(key)
		builder.WriteString("'){event.preventDefault();")
		builder.WriteString(action)
		builder.WriteString("}")
	}
	builder.WriteString("});")
	return f.AddJS(builder.String())
}
