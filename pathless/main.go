package pathless

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
)

type Pathless struct {
	Driver *mux.Router
	HTML   *template.HTML
	Font   string
	Color  string
	Md     *goldmark.Markdown
}

func InitPathless(color string, body template.HTML) *Pathless {
	p := &Pathless{
		Driver: mux.NewRouter().StrictSlash(true),
		Font:   "'Roboto', sans-serif",
		Color:  color,
		Md:     initGoldmark(),
	}
	p.Zero(body)

	p.Driver.HandleFunc("/", p.one).Methods("GET")
	go func() {
		log.Println("Starting pathless on :10101")
		http.ListenAndServe(":10101", p.Driver)
	}()
	return p
}

func (p *Pathless) one(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(*p.HTML))
}
