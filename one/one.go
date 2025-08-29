package one

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"

	"github.com/zachklingbeil/factory/fx"
)

type Factory struct {
	*fx.Fx
}

func NewFactory() *Factory {
	one := &Factory{
		Fx: fx.Init(),
	}
	one.Path("/").HandlerFunc(one.Pathless)
	return one
}

func (o *Factory) Pathless(w http.ResponseWriter, r *http.Request) {
	valChan := o.Observe("count")
	val := <-valChan
	count := val.(*atomic.Value).Load().(int)

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

	frame := o.Frames[current]
	fmt.Fprint(w, *frame)
}
