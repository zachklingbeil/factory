package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Circuit struct {
	One map[Zero]any
	Db  *Database
	Ctx context.Context
}

func NewCircuit(db *Database, ctx context.Context, distance time.Duration) *Circuit {
	circuit := &Circuit{
		One: make(map[Zero]any),
		Db:  db,
		Ctx: ctx,
	}
	if err := circuit.Loop(distance); err != nil {
		fmt.Printf("Error initializing Circuit: %v\n", err)
	}
	return circuit
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

func (c *Circuit) Coordinate(timestamp int64) Zero {
	t := time.UnixMilli(timestamp)
	return Zero{
		Year:        uint8(t.Year() - 2015),
		Month:       uint8(t.Month()),
		Day:         uint8(t.Day()),
		Hour:        uint8(t.Hour()),
		Minute:      uint8(t.Minute()),
		Second:      uint8(t.Second()),
		Millisecond: uint16(t.Nanosecond() / 1e6),
	}
}

func (c *Circuit) Index(data any) error {
	// Assert that the input is a Block
	block, ok := data.(Block)
	if !ok {
		return fmt.Errorf("invalid data type, expected Block")
	}

	// Generate the Zero coordinate using the block's timestamp
	coord := c.Coordinate(block.Timestamp)

	// Iterate over the transactions in the block
	for i, transaction := range block.Transactions {
		// Create a new coordinate for each transaction
		transactionCoord := coord
		transactionCoord.Index = uint16(i + 1) // Assign the index starting from 1

		// Add the transaction to the Circuit's One map with the Zero key
		c.One[transactionCoord] = transaction
	}

	return nil
}

// Loop loads all data from Redis into a map, starts Present as a goroutine, and returns the map.
func (c *Circuit) Loop(distance time.Duration) error {
	if err := c.Past(); err != nil {
		return fmt.Errorf("failed to load data from Redis: %w", err)
	}
	go c.Present(distance)
	return nil
}

// Present saves the Circuit's One map to Redis every distance until the context is cancelled.
func (c *Circuit) Present(distance time.Duration) {
	ticker := time.NewTicker(distance)
	defer ticker.Stop()
	ctx := c.Ctx
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.Db.Rdb.Do(ctx, "SELECT", 0).Err(); err != nil {
				continue
			}
			c.Db.Rw.RLock()
			for key, value := range c.One {
				keyJSON, err := json.Marshal(key)
				if err != nil {
					continue
				}
				valueJSON, err := json.Marshal(value)
				if err != nil {
					continue
				}
				redisKey := string(keyJSON)
				_ = c.Db.Rdb.Set(ctx, redisKey, valueJSON, 0).Err()
			}
			c.Db.Rw.RUnlock()
		}
	}
}

// Past loads all keys from Redis into the Circuit's One map.
func (c *Circuit) Past() error {
	if err := c.Db.Rdb.Do(c.Ctx, "SELECT", 0).Err(); err != nil {
		return err
	}
	c.Db.Mu.Lock()
	defer c.Db.Mu.Unlock()
	iter := c.Db.Rdb.Scan(c.Ctx, 0, "*", 0).Iterator()
	for iter.Next(c.Ctx) {
		redisKey := iter.Val()
		var key Zero
		if err := json.Unmarshal([]byte(redisKey), &key); err != nil {
			continue
		}
		valueJSON, err := c.Db.Rdb.Get(c.Ctx, redisKey).Bytes()
		if err != nil {
			return err
		}
		var value any
		if err := json.Unmarshal(valueJSON, &value); err != nil {
			return err
		}
		c.One[key] = value
	}
	if err := iter.Err(); err != nil {
		return err
	}
	return nil
}
