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

func (p *Pathless) One(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(*p.HTML))
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
<body>{{.Body}}</body>
</html>`

	tmpl := template.Must(template.New("page").Parse(templateStr))
	var buf bytes.Buffer

	data := struct {
		Body template.HTML
		CSS  template.CSS
	}{Body: body, CSS: cssContent}

	if err := tmpl.Execute(&buf, data); err != nil {
		result := template.HTML(strings.ReplaceAll(templateStr, "{{.Body}}", string(body)))
		p.HTML = &result
		return
	}
	result := template.HTML(buf.String())
	p.HTML = &result
}
