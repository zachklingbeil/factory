package json

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type JSON struct {
	HTTP *http.Client
	CTX  context.Context
}

// Json initializes and returns a new JSON utility instance, using http.Client and context.Context created in NewFactory.
func Json(http http.Client, ctx context.Context) *JSON {
	return &JSON{
		HTTP: &http,
		CTX:  context.Background(),
	}
}

// Print value as indented JSON to the standard output or logs error when value cannot be marshaled.
func (j *JSON) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}
