package main

import (
	"net/http"

	"github.com/zachklingbeil/factory/one"
)

func main() {
	factory := one.NewFactory()
	go func() {
		http.ListenAndServe(":1001", factory.Router)
	}()
	select {}
}
