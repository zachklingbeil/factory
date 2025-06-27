package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Pathless struct {
	Dir    string
	Router *mux.Router
	CTX    context.Context
}

func ServePathless(dir, port string, router *mux.Router) *Pathless {
	root := &Pathless{
		Dir:    dir,
		Router: router,
	}
	root.Router.PathPrefix("/").Handler(http.FileServer(http.Dir(root.Dir)))
	go func() {
		log.Fatal(http.ListenAndServe(":"+port, root.Router))
	}()
	return root
}
