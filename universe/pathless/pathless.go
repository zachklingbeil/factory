package pathless

import "html/template"

type Pathless struct {
	Favicon   string
	Title     string
	Font      string
	Primary   string
	Secondary string
	URL       string
	HTML      *template.HTML
}

func NewPathless(favicon, title, url string) *Pathless {
	pathless := &Pathless{
		Favicon:   favicon,
		Title:     title,
		Font:      "'Roboto', sans-serif",
		Primary:   "blue",
		Secondary: "red",
		URL:       url,
	}
	html := pathless.Zero()
	pathless.HTML = &html
	return pathless
}

func (p *Pathless) Zero() template.HTML {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <link rel="icon" href="` + p.Favicon + `" />
        <title>` + p.Title + `</title>
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
                font-family: ` + p.Font + `;
                scroll-behavior: smooth;
                box-sizing: border-box;
                border-radius: 0.3125em;
                display: flex;
                flex-direction: column;
            }
            body {
                border: medium solid ` + p.Primary + `;
            }
        </style>
        <script>
            async function loadFrame() {
                try {
                    const response = await fetch('` + p.URL + `');
                    const html = await response.text();
                    document.body.innerHTML = html;
                } catch (error) {
                    console.error('Failed to load frame:', error);
                }
            }            
            // Load initial frame on page load
            document.addEventListener('DOMContentLoaded', function() {
                loadFrame();
            });
        </script>
    </head>
    <body>
    </body>
</html>`
	return template.HTML(tmpl)
}
