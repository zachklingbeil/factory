package fx

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"maps"

	"github.com/redis/go-redis/v9"
)

type State struct {
	Map        map[string]any
	Json       *JSON
	Mu         *sync.Mutex
	Rw         *sync.RWMutex
	Redis      *redis.Client
	Ctx        context.Context
	MapHistory []map[string]any // Slice to store iterations of the map
}

func NewState(json *JSON, mu *sync.Mutex, rw *sync.RWMutex, redis *redis.Client, ctx context.Context) *State {
	state := &State{
		Map:        make(map[string]any),
		Json:       json,
		Mu:         mu,
		Rw:         rw,
		Redis:      redis,
		Ctx:        ctx,
		MapHistory: []map[string]any{}, // Initialize the slice
	}

	state.Add("t0", time.Now().Format("15:04:05.000"))
	return state
}

func (s *State) Add(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Map[key] = value

	mapCopy := make(map[string]any)
	maps.Copy(mapCopy, s.Map)
	s.MapHistory = append(s.MapHistory, mapCopy)

	t := time.Now().Format("15:04:05.000")
	state, _ := json.Marshal(s.Map)
	s.Redis.SAdd(s.Ctx, "state", t, state, 0)
	s.Get()
	return nil
}

func (s *State) Get() {
	s.Rw.RLock()
	defer s.Rw.RUnlock()
	s.Json.Print(s.Map)
}

func (s *State) Continue() {

}
