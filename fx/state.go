package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type State struct {
	Map     map[string]*One // Map now stores One for each key
	Json    *JSON
	Mu      *sync.Mutex
	Rw      *sync.RWMutex
	Data    *Data
	Ctx     context.Context
	Counter int64 // Incrementing counter for sorted set scores
}

// One represents the current value and its changelog
type One struct {
	Zero  any   `json:"current_value"`
	Delta []any `json:"change_log"` // History of previous values
}

func NewState(json *JSON, mu *sync.Mutex, rw *sync.RWMutex, data *Data, ctx context.Context) *State {
	state := &State{
		Map:  make(map[string]*One),
		Json: json,
		Mu:   mu,
		Rw:   rw,
		Data: data,
		Ctx:  ctx,
	}
	state.LoadMostRecent()
	state.Add("t0", time.Now().Format("04:05.0000000"))
	return state
}

func (s *State) Add(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	// Check if the key already exists
	if entry, exists := s.Map[key]; exists {
		// If the value has changed, update the changelog
		if entry.Zero != value {
			entry.Delta = append(entry.Delta, entry.Zero) // Add the old value to the changelog
			entry.Zero = value                            // Update the current value
			fmt.Printf("Value for key '%s' has changed. Changelog updated.\n", key)
		} else {
			return nil
		}
	} else {
		// If the key doesn't exist, create a new entry
		s.Map[key] = &One{
			Zero:  value,
			Delta: []any{},
		}
	}

	// Increment the counter for the score
	s.Counter++

	// Marshal the state map to JSON
	state, err := json.Marshal(s.Map)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	// Add the state to the sorted set "state" with the counter as the score
	if err := s.Data.RB.ZAdd(s.Ctx, "state", redis.Z{
		Score:  float64(s.Counter),
		Member: state,
	}).Err(); err != nil {
		return fmt.Errorf("failed to add state to sorted set: %w", err)
	}

	return nil
}
func (s *State) Get() {
	s.Rw.RLock()
	defer s.Rw.RUnlock()
	s.Json.Print(s.Map)
}

func (s *State) LoadMostRecent() (map[string]any, error) {
	result, err := s.Data.RB.ZRevRange(s.Ctx, "state", 0, 0).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve most recent state: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no state found in the sorted set")
	}

	var mostRecent map[string]any
	if err := json.Unmarshal([]byte(result[0]), &mostRecent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal most recent state: %w", err)
	}
	return mostRecent, nil
}
