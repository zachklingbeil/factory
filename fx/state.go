package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type State struct {
	Mu   *sync.Mutex
	Data *Data
	Ctx  context.Context
	Map  map[string]map[string]any
}

func NewState(data *Data, ctx context.Context) *State {
	state := &State{
		Mu:   &sync.Mutex{},
		Data: data,
		Ctx:  ctx,
		Map:  make(map[string]map[string]any),
	}
	state.LoadState()
	return state
}

func (s *State) Add(pkg string, key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if _, exists := s.Map[pkg]; !exists {
		s.Map[pkg] = make(map[string]any)
	}

	s.Map[pkg][key] = value

	state, err := json.Marshal(s.Map)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	result, err := s.Data.RB.ZRevRangeWithScores(s.Ctx, "state", 0, 0).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve highest score from sorted set: %w", err)
	}

	var nextScore float64 = 1
	if len(result) > 0 {
		nextScore = result[0].Score + 1
	}

	if err := s.Data.RB.ZAdd(s.Ctx, "state", redis.Z{
		Score:  nextScore,
		Member: state,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add state to sorted set: %w", err)
	}
	return nil
}

func (s *State) LoadState() error {
	result, err := s.Data.RB.ZRevRangeWithScores(s.Ctx, "state", 0, 0).Result()
	if err != nil {
		return fmt.Errorf("failed to retrieve latest state from Redis: %w", err)
	}
	if len(result) == 0 {
		return fmt.Errorf("no state found in Redis")
	}
	latestState := result[0].Member.(string)
	if err := json.Unmarshal([]byte(latestState), &s.Map); err != nil {
		return fmt.Errorf("failed to unmarshal latest state: %w", err)
	}
	return nil
}

func (s *State) GetValue(pkg string, key string) (any, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	packageMap, exists := s.Map[pkg]
	if !exists {
		return nil, fmt.Errorf("package %s not found in state", pkg)
	}
	value, exists := packageMap[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in package %s", key, pkg)
	}
	return value, nil
}

func (s *State) GetMap(pkg string) (map[string]any, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	packageMap, exists := s.Map[pkg]
	if !exists {
		return nil, fmt.Errorf("package %s not found in state", pkg)
	}
	return packageMap, nil
}
