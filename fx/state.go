package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type State struct {
	Mu   *sync.Mutex
	Data *Data
	Ctx  context.Context
	Map  map[string]any
}

func NewState(data *Data, ctx context.Context) *State {
	state := &State{
		Mu:   &sync.Mutex{},
		Data: data,
		Ctx:  ctx,
		Map:  make(map[string]any),
	}
	state.LoadState()
	return state
}

func (s *State) LoadState() error {
	result, err := s.Data.RB.Get(s.Ctx, "state").Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve state from Redis: %w", err)
	}
	if result == "" {
		return fmt.Errorf("no state found in Redis")
	}
	if err := json.Unmarshal([]byte(result), &s.Map); err != nil {
		return fmt.Errorf("failed to unmarshal state: %w", err)
	}
	return nil
}

func (s *State) Get(key string) (any, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	val, ok := s.Map[key]
	return val, ok
}

func (s *State) Count(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Map[key] = value
	state, err := json.Marshal(s.Map)
	if err != nil {
		return err
	}
	return s.Data.RB.Set(s.Ctx, "state", state, 0).Err()
}

func (s *State) Up(ctx context.Context, key string, from int64) <-chan int64 {
	ch := make(chan int64)
	go func() {
		defer close(ch)
		for i := from; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				s.Count(key, i)
			}
		}
	}()
	return ch
}

func (s *State) Down(ctx context.Context, key string, from int64) <-chan int64 {
	ch := make(chan int64)
	go func() {
		defer close(ch)
		for i := from; i >= 1; i-- {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
				s.Count(key, i)
			}
		}
	}()
	return ch
}

// this will become left and right, for adding and subtracting
// func (m *Math) Int(value string) *big.Int {
// 	bigIntValue := new(big.Int)
// 	if _, ok := bigIntValue.SetString(value, 10); !ok {
// 		log.Error("Failed to convert string to big.Int: %s", value)
// 	}
// 	return bigIntValue
// }

// func (m *Math) StringToInt(value string) *big.Int {
// 	i := new(big.Int)
// 	_, ok := i.SetString(value, 10)
// 	if !ok {
// 		log.Error("Failed to convert string to big.Int: %s", value)
// 		return nil
// 	}
// 	return i
// }
