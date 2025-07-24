package zero

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"strings"
)

type Frame interface {
	Build(class string, elements []One)
	Pathless(css, js string, body One)
	JS(js string) One
	CSS(css string) One
	AddKeybind(containerId string, keyHandlers map[string]string) One
	AddScrollKeybinds() One
	FileToString(path string) string
	AddFrame(frame *One)
	GetFrame(index int) (*One, bool)
	Text
	Element
}

// --- frame Implementation ---
type frame struct {
	Text
	Element
	frames []*One
	count  uint
}

func NewFrame() Frame {
	return &frame{
		Text:    NewText(),
		Element: NewElement(),
		frames:  make([]*One, 0),
		count:   0,
	}
}

func (f *frame) Pathless(css, js string, body One) {
	c := f.FileToString(css)
	j := f.FileToString(js)
	f.Build("pathless", []One{
		One(`<!DOCTYPE html>`),
		One(`<html lang="en">`),
		One(`<head>`),
		One(`<meta charset="UTF-8" />`),
		One(`<meta name="viewport" content="width=device-width, initial-scale=1.0" />`),
		One(`<title>hello universe</title>`),
		One(f.CSS(string(c))),
		One(f.JS(string(j))),
		One(`</head>`),
		One(fmt.Sprintf(`<body><div id="one">%s</div></body></html>`, string(body))),
	})
}

// Add a finalized frame to the collection
func (f *frame) AddFrame(frame *One) {
	f.frames = append(f.frames, frame)
	f.count++
}

// Retrieve a frame by index
func (f *frame) GetFrame(index int) (*One, bool) {
	if index < 0 || index >= int(f.count) {
		return nil, false
	}
	return f.frames[index], true
}

func (f *frame) FileToString(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(file)
}

func (f *frame) Build(class string, elements []One) {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}

	consolidatedContent := One(template.HTML(b.String()))
	result := One(template.HTML(fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), string(consolidatedContent))))
	f.AddFrame(&result)
}

func (f *frame) JS(js string) One {
	return One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
}

func (f *frame) CSS(css string) One {
	return One(template.HTML(fmt.Sprintf(`<style>%s</style>`, css)))
}

func (f *frame) AddKeybind(containerId string, keyHandlers map[string]string) One {
	var handlers strings.Builder
	for key, handlerCode := range keyHandlers {
		handlers.WriteString(fmt.Sprintf(`
         if (event.key === %q) {
            %s
         }
        `, key, handlerCode))
	}
	js := fmt.Sprintf(`
document.addEventListener('DOMContentLoaded', () => {
   const container = document.getElementById(%q);
   if (!container) return;
   container.tabIndex = 0;
   container.addEventListener('keydown', (event) => {
      %s
   });
});
`, containerId, handlers.String())
	return One(template.HTML(fmt.Sprintf(`<script>%s</script>`, js)))
}

func (f *frame) AddScrollKeybinds() One {
	return f.JS(
		`document.addEventListener('keydown', function(event) {
            const c = document.getElementById('one');
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
