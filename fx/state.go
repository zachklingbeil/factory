package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
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
	result, err := s.Data.RB.HGetAll(s.Ctx, "state").Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve state from Redis: %w", err)
	}
	if len(result) == 0 {
		return fmt.Errorf("no state found in Redis")
	}
	var latestTimestamp string
	for ts := range result {
		if ts > latestTimestamp {
			latestTimestamp = ts
		}
	}
	latestState := result[latestTimestamp]
	if err := json.Unmarshal([]byte(latestState), &s.Map); err != nil {
		return fmt.Errorf("failed to unmarshal latest state: %w", err)
	}
	return nil
}

func (s *State) Count(key string, value any, persist bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.Map[key] = value
	if !persist {
		return nil
	}

	state, err := json.Marshal(s.Map)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	timestamp := strconv.FormatInt(time.Now().UnixMicro(), 10)
	if err := s.Data.RB.HSet(s.Ctx, "state", timestamp, state).Err(); err != nil {
		return fmt.Errorf("failed to add state to hash: %w", err)
	}
	return nil
}

func (s *State) Read(key string) (any, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	value, exists := s.Map[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in state", key)
	}
	return value, nil
}
