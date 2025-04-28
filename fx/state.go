package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"maps"

	"github.com/redis/go-redis/v9"
)

type State struct {
	Mu        *sync.Mutex
	Data      *Data
	Ctx       context.Context
	Current   map[string]any
	ChangeLog []map[string]any
}

func NewState(data *Data, ctx context.Context) *State {
	state := &State{
		Mu:        &sync.Mutex{},
		Data:      data,
		Ctx:       ctx,
		Current:   make(map[string]any),
		ChangeLog: []map[string]any{},
	}
	return state
}

func (s *State) Add(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if existingValue, exists := s.Current[key]; exists && existingValue == value {
		return nil
	}

	if len(s.Current) > 0 {
		change := make(map[string]any)
		maps.Copy(change, s.Current)
		s.ChangeLog = append([]map[string]any{change}, s.ChangeLog...)
	}

	s.Current[key] = value

	state, err := json.Marshal(s.Current)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	score := float64(time.Now().UnixNano())

	if err := s.Data.RB.ZAdd(s.Ctx, "state", redis.Z{
		Score:  score,
		Member: state,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add state to sorted set: %w", err)
	}
	return nil
}

func (s *State) LoadLatestState() error {
	result, err := s.Data.RB.ZRevRangeWithScores(s.Ctx, "state", 0, 0).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve latest state from Redis: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no state found in Redis")
	}

	latestState := result[0].Member.(string)
	if err := json.Unmarshal([]byte(latestState), &s.Current); err != nil {
		return fmt.Errorf("failed to unmarshal latest state: %w", err)
	}
	return nil
}
