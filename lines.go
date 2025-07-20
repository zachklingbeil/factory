package factory

import (
	"log"
	"net/http"
)

func (f *Factory) InitPathless(cssPath, scriptPath string) {
	f.SetPathless(cssPath, scriptPath)
	f.HandleFunc("/", f.One).Methods("GET")
	go func() {
		log.Println("Starting pathless on :1001")
		http.ListenAndServe(":1001", f.Router)
	}()
}

// SetMapValue sets a value in the Factory's Map.
func (f *Factory) Input(key string, value any) {
	f.RWMutex.Lock()
	defer f.RWMutex.Unlock()
	f.Map[key] = &value
}

// GetMapValue retrieves a value from the Factory's Map.
func (f *Factory) Output(key string) (any, bool) {
	f.RWMutex.RLock()
	defer f.RWMutex.RUnlock()
	valPtr, ok := f.Map[key]
	if !ok || valPtr == nil {
		return nil, false
	}
	return *valPtr, true
}
