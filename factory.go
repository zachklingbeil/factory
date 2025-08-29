package main

import (
	"net/http"

	"github.com/zachklingbeil/factory/one"
)

type Factory struct {
	*one.One
}

func NewFactory() *Factory {
	return &Factory{
		One: one.NewOne(),
	}
}

func main() {
	factory := one.NewOne()
	go func() {
		http.ListenAndServe(":1001", factory.Router)
	}()
	select {}
}
