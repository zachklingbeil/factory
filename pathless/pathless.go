package pathless

import "html/template"

type Pathless struct {
	Favicon   string
	Title     string
	Font      string
	Primary   string
	Secondary string
	HTML      template.HTML
	Body      template.HTML
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
    </head>
    <body>
        {{.Body}}
    </body>
</html>`
	return template.HTML(tmpl)
}
