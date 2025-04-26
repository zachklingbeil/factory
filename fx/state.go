package fx

type State struct {
	Map  map[string]any
	Json *JSON
}

func NewState() *State {
	return &State{
		Map: make(map[string]any),
	}
}
func (s *State) Add(key string, value any) {
	s.Map[key] = value
}

func (s *State) Get() {
	s.Json.Print(s.Map)
}
