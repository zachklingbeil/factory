package pathless

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
)

type Pathless struct {
	HTML  *template.HTML
	Font  string
	Color string
}

func NewPathless(color string) *Pathless {
	pathless := &Pathless{
		Font:  "'Roboto', sans-serif",
		Color: color,
	}
	return pathless
}

func (p *Pathless) One(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(*p.HTML))
}

func (p *Pathless) Zero(body template.HTML) {
	templateStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>hello universe</title>
    <style>
        *, *::before, *::after {
            box-sizing: border-box;
            margin: 0;
            scrollbar-width: none;
            -ms-overflow-style: none;
            user-select: none;
        }
        *::-webkit-scrollbar { display: none; }
        html, body {
            color: white;
            background-color: black;
            height: 100vh;
            width: 100vw;
            font-family: ` + p.Font + `;
            scroll-behavior: smooth;
            overflow: hidden;
        }
        body {
            border: medium solid ` + p.Color + `; 
            border-radius: 0.3125em;
            display: flex;
        }
        main {
            flex: 1;
            display: flex;
            flex-direction: column;
            overflow-y: auto;
        }
    </style>
</head>
<body><main>{{.Body}}</main></body>
</html>`

	tmpl := template.Must(template.New("page").Parse(templateStr))
	var buf bytes.Buffer

	data := struct{ Body template.HTML }{Body: body}
	if err := tmpl.Execute(&buf, data); err != nil {
		result := template.HTML(strings.ReplaceAll(templateStr, "{{.Body}}", string(body)))
		p.HTML = &result
		return
	}
	result := template.HTML(buf.String())
	p.HTML = &result
}
