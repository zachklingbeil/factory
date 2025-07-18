package frame

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Frame struct {
	Md     *goldmark.Markdown
	Frames []*template.HTML
	Map    map[string]*template.HTML
	*mux.Router
}

func NewFrame(mux *mux.Router) *Frame {
	frame := &Frame{
		Md:     initGoldmark(),
		Frames: []*template.HTML{},
		Router: mux,
		Map:    make(map[string]*template.HTML),
	}

	// prev := frame.AddNavKeybind("q", "/frame/prev")
	// next := frame.AddNavKeybind("e", "/frame/next")
	return frame
}

func (f *Frame) AddFrame(reference string, elements ...template.HTML) *template.HTML {
	var builder strings.Builder
	for _, element := range elements {
		builder.WriteString(string(element))
	}
	result := builder.String()
	if reference != "" {
		result = `<div class="` + reference + `">` + result + `</div>`
	}
	frame := template.HTML(result)
	f.Frames = append(f.Frames, &frame)
	f.Map["/frame/"+reference] = &frame // Use "/reference" as path
	return &frame
}

// Register endpoint template and route
func (f *Frame) RegisterFrameRoutes() {
	for path, tmpl := range f.Map {
		f.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			f.WriteResponse(w, tmpl)
		})
	}
}

// Write HTML response
func (f *Frame) WriteResponse(w http.ResponseWriter, tmpl *template.HTML) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if tmpl != nil {
		w.Write([]byte(string(*tmpl)))
	} else {
		w.Write([]byte("<div>404 Not Found</div>"))
	}
}

func (f *Frame) AddCSS(styles map[string]string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<style>")
	for selector, rules := range styles {
		builder.WriteString(selector)
		builder.WriteString(" { ")
		builder.WriteString(rules)
		builder.WriteString(" }\n")
	}
	builder.WriteString("</style>")
	html := template.HTML(builder.String())
	return &html
}

func (f *Frame) AddJS(js string) *template.HTML {
	var builder strings.Builder
	builder.WriteString("<script>")
	builder.WriteString(js)
	builder.WriteString("</script>")
	html := template.HTML(builder.String())
	return &html
}

func (f *Frame) AddNavKeybind(key, path string) *template.HTML {
	return f.AddJS(
		`document.addEventListener('keydown', function(event) {
            if (event.key === '` + key + `') {
                console.log('Keybind pressed: ` + key + ` -> ` + path + `');
                fetch('` + path + `')
                    .then(r => r.text())
                    .then(html => {
                        const c = document.getElementById('frame');
                        if (c) c.innerHTML = html;
                    });
            }
        });`,
	)
}

func (f *Frame) AddScrollKeybinds() *template.HTML {
	return f.AddJS(
		`document.addEventListener('keydown', function(event) {
            const c = document.getElementById('frame');
            if (!c) return;
            if (event.key === 'w') {
                c.scrollBy({ top: -100, behavior: 'smooth' });
            }
            if (event.key === 's') {
                c.scrollBy({ top: 100, behavior: 'smooth' });
            }
        });`,
	)
}

func (f *Frame) RegisterNavRoutes() {
	f.Router.HandleFunc("/frame/prev", func(w http.ResponseWriter, r *http.Request) {
		// Serve previous frame content
		w.Write([]byte("<div>Previous Frame</div>"))
	})
	f.Router.HandleFunc("/frame/next", func(w http.ResponseWriter, r *http.Request) {
		// Serve next frame content
		w.Write([]byte("<div>Next Frame</div>"))
	})
}
