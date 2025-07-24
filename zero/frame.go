package zero

import (
	"fmt"
	"html"
	"html/template"
	"os"
	"strconv"
	"strings"
)

type Frame interface {
	Build(class string, elements []One)
	Merge(elements ...One) One // Merge multiple elements into one
	Wrap(class string, elements ...One) One
	Pathless(css, js string, body One)
	JS(js string) One
	CSS(css string) One
	AddKeybind(containerId string, keyHandlers map[string]string) One
	AddScrollKeybinds() One
	FileToString(path string) string
	AddFrame(frame *One)
	GetFrame(index int) (*One, bool)
	FrameCount() string
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

func (f *frame) Merge(elements ...One) One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}
	return One(template.HTML(b.String()))
}

func (f *frame) Wrap(class string, elements ...One) One {
	var b strings.Builder
	for _, el := range elements {
		b.WriteString(string(el))
	}
	consolidatedContent := template.HTML(b.String())
	result := fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), string(consolidatedContent))
	return One(template.HTML(result))
}

func (f *frame) Pathless(css, js string, body One) {
	c := f.FileToString(css)
	j := f.FileToString(js)

	result := f.Merge(
		One(`<!DOCTYPE html>`),
		One(`<html lang="en">`),
		One(`<head>`),
		One(`<meta charset="UTF-8" />`),
		One(`<meta name="viewport" content="width=device-width, initial-scale=1.0" />`),
		One(`<title>hello universe</title>`),
		f.CSS(c),
		f.JS(j),
		One(`</head>`),
		One(fmt.Sprintf(`<body><div id="one">%s</div></body></html>`, string(body))),
	)
	f.AddFrame(&result)
}

func (f *frame) FrameCount() string {
	return strconv.Itoa(int(f.count))
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
