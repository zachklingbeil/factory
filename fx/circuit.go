package fx

import (
	"context"
	"maps"
	"sync"
)

type Circuit struct {
	One map[Zero]any
	Ctx context.Context
	Mu  *sync.Mutex
}

type Zero struct {
	Block       uint64 `json:"block"`
	Year        uint8  `json:"year"`
	Month       uint8  `json:"month"`
	Day         uint8  `json:"day"`
	Hour        uint8  `json:"hour"`
	Minute      uint8  `json:"minute"`
	Second      uint8  `json:"second"`
	Millisecond uint16 `json:"millisecond"`
	Index       uint16 `json:"index"`
}

func NewCircuit(ctx context.Context, mu *sync.Mutex) *Circuit {
	circuit := &Circuit{
		One: make(map[Zero]any),
		Ctx: ctx,
		Mu:  mu,
	}
	return circuit
}

// Add safely adds one to Circuit.One
func (c *Circuit) Add(one map[Zero]any) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	maps.Copy(c.One, one)
}

func (c *Circuit) Read(zero Zero) any {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	value := c.One[zero]
	return value
}
