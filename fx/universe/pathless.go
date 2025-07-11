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

func NewPathless(favicon, title, font, primary, secondary string) *Pathless {
	p := &Pathless{
		Config: &Config{
			Favicon:   favicon,
			Title:     title,
			Font:      font,
			Primary:   primary,
			Secondary: secondary,
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
	return template.HTML(`
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="{{.Favicon}}" />
        <title>{{.Title}}</title>
        <style>
            :root {
                --font-family: {{.Font}};
                --primary: {{.Primary}};
                --secondary: {{.Secondary}};
            }
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
                font-family: var(--font-family);
                scroll-behavior: smooth;
                box-sizing: border-box;
                border-radius: 0.3125em;
                display: flex;
                flex-direction: column;
            }
            body {
                border: medium solid var(--primary);
            }
        </style>
    </head>
    <body>
        {{.Body}}
    </body>
</html>`)
}
