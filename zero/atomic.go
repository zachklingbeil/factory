package zero

import (
	"encoding/json"
	"os"
	"sync/atomic"
)

func (z *Zero) Add(key string, val int) {
	z.Lock()
	defer z.Unlock()
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
}

func (z *Zero) Observe(key string) <-chan int {
	ch := make(chan int, 1)
	z.RLock()
	if v, exists := z.Map[key]; exists {
		if i, ok := v.Load().(int); ok {
			ch <- i
		}
	}
	z.RUnlock()
	return ch
}

func (z *Zero) Subtract(key string) {
	z.Lock()
	defer z.Unlock()
	delete(z.Map, key)
	for _, ch := range z.watchers[key] {
		select {
		case ch <- 0:
		default:
		}
	}
	z.Cond.Broadcast()
	z.save()
}

// save the Map map to the JSON file.
func (z *Zero) save() {
	plain := make(map[string]int)
	for k, v := range z.Map {
		if i, ok := v.Load().(int); ok {
			plain[k] = i
		}
	}
	data, _ := json.MarshalIndent(plain, "", "  ")
	_ = os.WriteFile("factory/atomic.json", data, 0644)
}

// // Load the Map map from the JSON file, creating the file if needed.
// func (z *Zero) Load() {
// 	z.Lock()
// 	defer z.Unlock()
// 	const filePath = "factory/atomic.json"
// 	if _, err := os.Stat(filePath); os.IsNotExist(err) {
// 		_ = os.WriteFile(filePath, []byte("{}"), 0644)
// 	}

// 	data, err := os.ReadFile(filePath)
// 	if err != nil || len(data) == 0 {
// 		return
// 	}
// 	var plain map[string]any
// 	if err := json.Unmarshal(data, &plain); err != nil {
// 		return
// 	}
// 	for k, v := range plain {
// 		av := &atomic.Value{}
// 		av.Store(v)
// 		z.Map[k] = av
// 	}
// }

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
