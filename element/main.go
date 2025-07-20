package element

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
)

type Zero struct {
	X []Frame
}

func NewZero() *Zero {
	return &Zero{
		X: make([]Frame, 0),
	}
}

type Element struct {
	template.HTML
	data any
}

type Frame interface {
	Render() template.HTML
	Data() any
	Update(data any)
}

func (e *Element) Render() template.HTML {
	return e.HTML
}

func (e *Element) Data() any {
	return e.data
}

func (e *Element) Update(data any) {
	e.data = data
}

func NewElement(htmlStr string) *Element {
	return &Element{HTML: template.HTML(htmlStr)}
}

// BuildFrame wraps multiple elements into a single frame and appends it to z.X.
func (z *Zero) BuildFrame(class string, elements []*Element) {
	var frames []Frame
	for _, e := range elements {
		frames = append(frames, e)
	}
	body := simplify(frames)
	wrapped := fmt.Sprintf(`<div class="%s">%s</div>`, html.EscapeString(class), body)
	z.X = append(z.X, NewElement(wrapped))
}

func (z *Zero) SetPathless(css, script string) {
	var body template.HTML
	if len(z.X) > 1 {
		body = z.X[1].Render()
	}
	templateStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>hello universe</title>
    <style>{{.CSS}}</style>
    <script>{{.Script}}</script>
</head>
<body>
    <div id="frame">{{.Body}}</div>
</body>
</html>`
	tmpl, err := template.New("pathless").Parse(templateStr)
	if err != nil {
		return
	}
	data := struct {
		CSS    string
		Script string
		Body   template.HTML
	}{
		CSS:    css,
		Script: script,
		Body:   body,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return
	}
	elem := NewElement(buf.String())
	if len(z.X) == 0 {
		z.X = append(z.X, elem)
	} else {
		z.X[0] = elem
	}
}
