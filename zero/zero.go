package zero

import (
	"context"

	"html/template"
	"sync"
	"sync/atomic"

	"github.com/gorilla/mux"
)

type One template.HTML

type Zero struct {
	Build
	Element Element
	context.Context
	*mux.Router
	*sync.RWMutex
	*sync.Cond
	Frames   []*One
	Map      map[string]*atomic.Value
	watchers map[string][]chan any
}

func NewZero() *Zero {
	rw := &sync.RWMutex{}
	zero := &Zero{
		RWMutex:  rw,
		Cond:     sync.NewCond(rw),
		Context:  context.Background(),
		Build:    NewBuild(),
		Frames:   make([]*One, 0),
		Router:   mux.NewRouter().StrictSlash(false),
		watchers: make(map[string][]chan any),
		Map:      make(map[string]*atomic.Value),
		Element:  NewElement(),
	}
	zero.AddFrame(zero.Pathless())
	return zero
}

func (z *Zero) AddFrame(frame *One) {
	z.Frames = append(z.Frames, frame)
	z.Add("count", len(z.Frames))
}
