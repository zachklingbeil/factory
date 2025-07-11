package pathless

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Pathless struct {
	Config *Config
	HTML   template.HTML
	Body   template.HTML
	router *mux.Router
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
		router: mux.NewRouter().StrictSlash(true),
		Config: &Config{
			Favicon:   favicon,
			Title:     title,
			Font:      font,
			Primary:   primary,
			Secondary: secondary,
		},
	}

	// Finalize the template with {{.Body}} placeholder
	p.HTML = p.baseTemplate()
	p.Body = template.HTML("")

	p.router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-API-KEY", "Content-Type", "Peer", "Cache-Control", "Connection"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	))
	p.router.HandleFunc("/", p.serve)
	go func() {
		log.Fatal(http.ListenAndServe(":10001", p.router))
	}()
	return p
}

func (p *Pathless) serve(w http.ResponseWriter, r *http.Request) {
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
