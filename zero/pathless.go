package zero

import (
	"fmt"
	"os"
)

func (z *Zero) BuildPathless(cssPath, domain string, frames int) One {
	css := z.ReadFileAsString(cssPath)
	head := z.Build([]One{
		One(`<meta charset="UTF-8" />`),
		One(`<meta name="viewport" content="width=device-width, initial-scale=1.0" />`),
		One(`<title>hello universe</title>`),
		z.CSS(css),
		z.JS(z.PathlessJS(frames, domain)),
	})

	mainBody := One(`<div id="one"></div>`)

	html := One(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>%s</head>
<body>%s</body>
</html>`, head, mainBody))

	return html
}

func (z *Zero) ReadFileAsString(path string) string {
	file, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(file)
}

func (z *Zero) PathlessJS(numFrames int, domain string) string {
	return fmt.Sprintf(`
let currentFrame = 0;
const totalFrames = %d;
const domain = "%s";

function updateBody(endpoint) {
    fetch(endpoint)
        .then((res) => res.text())
        .then((html) => {
            document.getElementById('one').innerHTML = html;
        });
}

function showFrame(idx) {
    currentFrame = ((idx %% totalFrames) + totalFrames) %% totalFrames; // wrap around
    updateBody(domain + "/frame/" + currentFrame);
}

document.addEventListener('DOMContentLoaded', function () {
    showFrame(currentFrame); // load initial frame

    document.addEventListener('keydown', function (e) {
        if (e.key === "q") {
            showFrame(currentFrame - 1);
        }
        if (e.key === "e") {
            showFrame(currentFrame + 1);
        }
    });
});
`, numFrames, domain)
}
