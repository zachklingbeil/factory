package zero

import (
	"encoding/json"
	"os"
	"sync/atomic"
)

func (z *Zero) Add(key string, val any) {
	z.Lock()
	if _, exists := z.Map[key]; !exists {
		z.Map[key] = &atomic.Value{}
	}
	z.Map[key].Store(val)
	for _, ch := range z.watchers[key] {
		select {
		case ch <- val:
		default:
		}
	}
	z.Cond.Broadcast()
	z.save()
	z.Unlock()
}

// Observe returns a channel that receives updates for the key.
func (z *Zero) Observe(key string) <-chan any {
	ch := make(chan any, 1)
	z.Lock()
	z.watchers[key] = append(z.watchers[key], ch)
	if v, exists := z.Map[key]; exists {
		ch <- v.Load()
	}
	z.Unlock()
	return ch
}

func (z *Zero) Subtract(key string) {
	z.Lock()
	delete(z.Map, key)
	for _, ch := range z.watchers[key] {
		select {
		case ch <- nil:
		default:
		}
	}
	z.Cond.Broadcast()
	z.save()
	z.Unlock()
}

// save the Map map to the JSON file.
func (z *Zero) save() {
	plain := make(map[string]any)
	for k, v := range z.Map {
		plain[k] = v.Load()
	}
	data, _ := json.MarshalIndent(plain, "", "  ")
	_ = os.WriteFile("factory/atomic.json", data, 0644)
}

// Load the Map map from the JSON file, creating the file if needed.
func (z *Zero) Load() {
	z.Lock()
	defer z.Unlock()
	const filePath = "factory/atomic.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.WriteFile(filePath, []byte("{}"), 0644)
	}

	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return
	}
	var plain map[string]any
	if err := json.Unmarshal(data, &plain); err != nil {
		return
	}
	z.Map = make(map[string]*atomic.Value)
	for k, v := range plain {
		av := &atomic.Value{}
		av.Store(v)
		z.Map[k] = av
	}
}

// WaitForCondition waits until the provided condition function returns true.
// It must be called with the lock held.
func (z *Zero) WaitForCondition(condFn func() bool) {
	z.Lock()
	defer z.Unlock()
	for !condFn() {
		z.Cond.Wait()
	}
}

// SignalOne signals one waiting goroutine.
func (z *Zero) SignalOne() {
	z.Lock()
	z.Cond.Signal()
	z.Unlock()
}

// BroadcastAll signals all waiting goroutines.
func (z *Zero) BroadcastAll() {
	z.Lock()
	z.Cond.Broadcast()
	z.Unlock()
}
