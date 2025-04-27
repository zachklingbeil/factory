package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type State struct {
	Map  map[string]any
	Json *JSON
	Mu   *sync.Mutex
	Rw   *sync.RWMutex
	Data *Data
	Ctx  context.Context
}

func NewState(json *JSON, mu *sync.Mutex, rw *sync.RWMutex, data *Data, ctx context.Context) *State {
	state := &State{
		Map:  make(map[string]any),
		Json: json,
		Mu:   mu,
		Rw:   rw,
		Data: data,
		Ctx:  ctx,
	}

	state.Add("t0", time.Now().Format("08:04:05.0000000"))
	return state
}

func (s *State) Add(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Map[key] = value

	state, err := json.Marshal(s.Map)
	if err != nil {
		return fmt.Errorf("failed to marshal state map: %w", err)
	}

	timestamp := time.Now().Format("08:04:05.0000000")
	if err := s.Data.RB.SAdd(s.Ctx, "state", state, timestamp).Err(); err != nil {
		return fmt.Errorf("failed to add state to Redis: %w", err)
	}

	return nil
}

func (s *State) Get() {
	s.Rw.RLock()
	defer s.Rw.RUnlock()
	s.Json.Print(s.Map)
}
