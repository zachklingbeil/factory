package zero

import (
	"fmt"
	"html/template"
)

type Zero struct {
	Frame
	Pathless *One
	Frames   []*One
}

func NewZero(css, js string) *Zero {
	z := &Zero{
		Frame:  NewFrame(),
		Frames: make([]*One, 0),
	}
	pathless := z.BuildPathless(css, js)
	z.Pathless = &pathless
	z.AddFrame(z.Pathless)
	return z
}

type One template.HTML

func (z *Zero) AddFrame(frame *One) {
	z.Frames = append(z.Frames, frame)
}

func (z *Zero) FrameCount() uint {
	return uint(len(z.Frames))
}

func (z *Zero) BuildPathless(css, js string) One {
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
<body><div id="one"></div></body>
</html>`, *head))

	return html
}
