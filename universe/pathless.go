package universe

import (
	"html/template"
	"net/http"
)

type Pathless struct {
	Config *Config
	HTML   template.HTML
	Body   template.HTML
}

type Config struct {
	Favicon   string
	Title     string
	Font      string
	Primary   string
	Secondary string
}

func NewPathless(favicon, title string) *Pathless {
	p := &Pathless{
		Config: &Config{
			Favicon:   favicon,
			Title:     title,
			Font:      "'Roboto', sans-serif",
			Primary:   "blue",
			Secondary: "red",
		},
	}
	p.HTML = p.baseTemplate()
	p.Body = template.HTML("")
	return p
}

func (p *Pathless) Serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(p.HTML))
}

func (p *Pathless) Update(w http.ResponseWriter, r *http.Request, content string) {
	p.Body = template.HTML(content)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(p.Body))
}

func (p *Pathless) baseTemplate() template.HTML {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="` + p.Config.Favicon + `" />
        <title>` + p.Config.Title + `</title>
        <style>
            *,
            *::before,
            *::after {
                box-sizing: border-box;
                margin: 0;
                scrollbar-width: none;
                -ms-overflow-style: none;
                user-select: none;
                -webkit-user-select: none;
                -moz-user-select: none;
                -ms-user-select: none;
            }
            *::-webkit-scrollbar {
                display: none;
            }
            html,
            body {
                color: white;
                background-color: black;
                overflow: hidden;
                height: 100vh;
                width: 100vw;
                font-family: ` + p.Config.Font + `;
                scroll-behavior: smooth;
                box-sizing: border-box;
                border-radius: 0.3125em;
                display: flex;
                flex-direction: column;
            }
            body {
                border: medium solid ` + p.Config.Primary + `;
            }
        </style>
    </head>
    <body>
        {{.Body}}
    </body>
</html>`
	return template.HTML(tmpl)
}
