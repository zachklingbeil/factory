package zero

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//go:embed pathless.html
var pathless string

var count int // package-level variable for frame count

type One template.HTML

type Zero struct {
	Build
	*mux.Router
	context.Context
	Frames []*One
}

func NewZero() *Zero {
	return &Zero{
		Build:   NewBuild(),
		Frames:  make([]*One, 0),
		Router:  mux.NewRouter(),
		Context: context.Background(),
	}
}

func (z *Zero) ZeroZero() {
	root := One(template.HTML(pathless))
	z.AddFrame(&root)
	z.Path("/").HandlerFunc(z.Pathless)
	go func() {
		http.ListenAndServe(":1001", z.Router)
	}()
}

func (z *Zero) AddFrame(frame *One) {
	z.Frames = append(z.Frames, frame)
	count = len(z.Frames)
}

func (z *Zero) Pathless(w http.ResponseWriter, r *http.Request) {
	current, err := strconv.Atoi(r.Header.Get("Y"))
	if err != nil || current < 0 || current >= count {
		current = 0
	}

	prev := (current - 1 + count) % count
	next := (current + 1) % count

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X", strconv.Itoa(prev))
	w.Header().Set("Y", strconv.Itoa(current))
	w.Header().Set("Z", strconv.Itoa(next))

	frame := z.Frames[current]
	fmt.Fprint(w, *frame)
}
