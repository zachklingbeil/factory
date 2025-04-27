package fx

import (
	"context"
	"encoding/json"
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

	state.Add("t0", time.Now().Format("15:04:05.000"))
	return state
}

func (s *State) Add(key string, value any) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.Map[key] = value

	t := time.Now().Format("15:04:05.000")
	state, _ := json.Marshal(s.Map)
	s.Data.RB.SAdd(s.Ctx, "state", t, state, 0)
	return nil
}

func (s *State) Get() {
	s.Rw.RLock()
	defer s.Rw.RUnlock()
	s.Json.Print(s.Map)
}
