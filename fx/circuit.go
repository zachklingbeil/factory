package fx

import (
	"context"
	"maps"
	"sync"
	"time"
)

type Circuit struct {
	One  map[any]any
	Zero map[Zero]any
	Ctx  context.Context
	Mu   *sync.Mutex
}

type Zero struct {
	Block       int64  `json:"block"`
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
		One: make(map[any]any),
		Ctx: ctx,
		Mu:  mu,
	}
	return circuit
}

// Add safely adds one to Circuit.One
func (c *Circuit) Add(one map[any]any) {
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

func (c *Circuit) Coordinates(blockNumber, timestamp int64, ones []any) (Zero, any, error) {
	for i := range ones {
		if tx, ok := ones[i].(map[string]any); ok {
			tx["index"] = i + 1
		}
	}

	t := time.UnixMilli(timestamp)
	coord := Zero{
		Block:       blockNumber,
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
		Index:       0,
	}
	return coord, ones, nil
}
