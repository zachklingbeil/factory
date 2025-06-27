package json

import (
	"context"
	"net/http"
)

type JSON struct {
	Ctx  context.Context
	HTTP *http.Client
}

func NewJSON(ctx context.Context) *JSON {
	return &JSON{
		Ctx:  ctx,
		HTTP: &http.Client{},
	}
}
