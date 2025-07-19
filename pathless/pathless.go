package pathless

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Pathless struct {
	HTML *template.HTML
}

func NewPathless() *Pathless {
	return &Pathless{}
}

func (p *Pathless) Zero(body template.HTML, cssPath string) {
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
</head>
<body>
    <div id="frame">{{.Body}}</div>
</body>
</html>`

	tmpl := template.Must(template.New("page").Parse(templateStr))
	var buf bytes.Buffer

	data := struct {
		Body template.HTML
		CSS  template.CSS
		Nav  template.HTML
	}{Body: body, CSS: cssContent, Nav: p.Nav()}

	if err := tmpl.Execute(&buf, data); err != nil {
		result := template.HTML(strings.ReplaceAll(templateStr, "{{.Body}}", string(body)))
		p.HTML = &result
		return
	}
	result := template.HTML(buf.String())
	p.HTML = &result
}

func (p *Pathless) Nav() template.HTML {
	return template.HTML(`
<script>
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
</script>
`)
}
func (p *Pathless) One(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(*p.HTML))
}
