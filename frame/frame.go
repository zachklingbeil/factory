package frame

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Frame struct {
	Md    *goldmark.Markdown
	Index []*template.HTML
	*mux.Router
}

func NewFrame(mux *mux.Router) *Frame {
	frame := &Frame{
		Md:     initGoldmark(),
		Index:  make([]*template.HTML, 0),
		Router: mux,
	}
	frame.HandleFunc("/frame", frame.FrameHandler())
	return frame
}

func (f *Frame) AddFrame(elements ...template.HTML) {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	value := builder.String()
	key, _ := intToWord(len(f.Index))
	value = `<div class="` + key + `">` + value + `</div>`
	frame := template.HTML(value)
	f.Index = append(f.Index, &frame)
}

// Serve frame by index from header "X", default to 0 if missing/invalid
func (f *Frame) FrameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		idx := 0
		count := len(f.Index)
		if idxStr := r.Header.Get("X"); idxStr != "" && count > 0 {
			if n, err := strconv.Atoi(idxStr); err == nil {
				idx = ((n % count) + count) % count
			}
		}
		if idx >= 0 && idx < count && f.Index[idx] != nil {
			_, _ = w.Write([]byte(string(*f.Index[idx])))
		} else {
			_, _ = w.Write([]byte("<div>404 Not Found</div>"))
		}
	}
}
