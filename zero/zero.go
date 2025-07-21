package zero

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type Zero struct {
	Frame
	Text
	Element
}

func NewZero() *Zero {
	return &Zero{
		Frame:   NewFrame(),
		Text:    NewText(),
		Element: NewElement(),
	}
}

type One template.HTML

func Tag(tag, text string) One {
	return One(template.HTML(fmt.Sprintf("<%s>%s</%s>", tag, html.EscapeString(text), tag)))
}

func ClosedTag(tag string, attrs map[string]string) One {
	var b strings.Builder
	b.WriteString("<")
	b.WriteString(tag)
	for k, v := range attrs {
		b.WriteString(fmt.Sprintf(` %s="%s"`, k, html.EscapeString(v)))
	}
	return One(template.HTML(b.String()))
}
