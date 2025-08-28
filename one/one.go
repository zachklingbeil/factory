package one

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zachklingbeil/factory/fx"
	"github.com/zachklingbeil/factory/zero"
)

type One struct {
	zero.Zero
	*fx.Fx
	Api map[string]zero.Coordinate
	*mux.Router
}

func NewOne() *One {
	o := &One{
		Zero:   zero.NewZero(),
		Fx:     fx.InitFx(),
		Router: mux.NewRouter(),
		Api:    make(map[string]zero.Coordinate),
	}
	o.Circuit()
	return o
}

func (o *One) NewRouter() *mux.Router {
	router := mux.NewRouter()
	go func() {
		http.ListenAndServe(":1001", router)
	}()
	return router
}
