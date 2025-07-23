package zero

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Zero struct {
	Frame
	Pathless One
	Body     One
	Frames   []One
}

func NewZero() *Zero {
	return &Zero{
		Frame:    NewFrame(),
		Pathless: "",
		Body:     "",
		Frames:   make([]One, 0),
	}
}

type One template.HTML

func (z *Zero) AddFrame(frame One, router *mux.Router) {
	z.Frames = append(z.Frames, frame)
	index := len(z.Frames) - 1
	path := "/"
	if index > 0 {
		path = "/frame/" + fmt.Sprint(index)
	}
	router.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Frames", fmt.Sprint(len(z.Frames)))
		w.Write([]byte(z.Frames[index]))
	}).Methods("GET")
}

func (z *Zero) FrameCount() uint {
	return uint(len(z.Frames))
}
func (z *Zero) Swap(id uint) {
	z.Body = z.Frames[id]
}

func (z *Zero) BuildPathless(css, js string) {
	c := z.FileToString(css)
	j := z.FileToString(js)
	head := z.Build([]One{
		One(`<meta charset="UTF-8" />`),
		One(`<meta name="viewport" content="width=device-width, initial-scale=1.0" />`),
		One(`<title>hello universe</title>`),
		z.CSS(c),
		z.JS(j),
	})
	html := One(fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>%s</head>
<body>%s</body>
</html>`, head, z.Body))
	z.Pathless = html
}
