package elements

import (
	"fmt"
	"html"
	"html/template"
)

type media struct{}

func (m *media) Img(src, alt string) template.HTML {
	return ClosedTag("img", map[string]string{"src": src, "alt": alt, "width": "100%", "height": "auto"})
}
func (m *media) Video(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<video controls src="%s"></video>`, html.EscapeString(src)))
}
func (m *media) Audio(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<audio controls src="%s"></audio>`, html.EscapeString(src)))
}
func (m *media) Iframe(src string) template.HTML {
	return template.HTML(fmt.Sprintf(`<iframe src="%s"></iframe>`, html.EscapeString(src)))
}
func (m *media) Embed(src string) template.HTML {
	return ClosedTag("embed", map[string]string{"src": src})
}
func (m *media) Source(src string) template.HTML {
	return ClosedTag("source", map[string]string{"src": src})
}
func (m *media) Canvas(id string) template.HTML {
	return template.HTML(fmt.Sprintf(`<canvas id="%s"></canvas>`, html.EscapeString(id)))
}
