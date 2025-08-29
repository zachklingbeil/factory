package main

import (
	"fmt"
	"net/http"

	"github.com/zachklingbeil/factory/one"
)

func main() {
	factory := one.NewFactory()
	go func() {
		http.ListenAndServe(":1001", factory.Router)
		fmt.Println("Server started at http://localhost:1001")
	}()
	select {}
}
