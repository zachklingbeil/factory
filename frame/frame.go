package frame

import (
	"bytes"
	"html/template"
	"os"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Frame struct {
	Md *goldmark.Markdown
}

func NewFrame(mux *mux.Router) *Frame {
	frame := &Frame{
		Md: initGoldmark(),
	}
	return frame
}

func (f *Frame) Index(cssPath string) {
	// Render the first frame in Frames (index 0)
	var body template.HTML

	file, err := os.ReadFile(cssPath)
	cssContent := template.CSS("")
	if err == nil {
		cssContent = template.CSS(file)
	}
	templateStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>hello universe</title>
    <style>{{.CSS}}</style>
    <script>{{.Nav}}</script>
</head>
<body>
    <div id="frame">{{.Body}}</div>
</body>
</html>`

	tmpl := template.Must(template.New("page").Parse(templateStr))
	var buf bytes.Buffer

}
func (f *Frame) NavJS() template.JS {
	return template.JS(`
document.addEventListener('DOMContentLoaded', function() {
    let frameIdx = 0;
    document.addEventListener('keydown', function(event) {
        if (event.key === 'q') frameIdx--;
        if (event.key === 'e') frameIdx++;
        if (event.key === 'q' || event.key === 'e') {
            fetch('/frame', {
                headers: { 'X': frameIdx }
            })
            .then(r => r.text())
            .then(html => {
                const c = document.getElementById('frame');
                if (c) c.innerHTML = html;
            });
        }
    });
});
`)
}
