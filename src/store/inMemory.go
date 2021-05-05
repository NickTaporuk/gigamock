package store

// InMemoryStore should include name of the route and which a scenario should use
type InMemoryStore struct {
	Data map[string]int
}

func (i *InMemoryStore) SetState(name string, status int) {
	i.Data[name] = status
}

func (i *InMemoryStore) State(name string) int {
	return i.Data[name]
}

func (i *InMemoryStore) Reset(name string) bool {
	delete(i.Data, name)
	return true
}

func (i *InMemoryStore) ResetAll() {
	i.Data = map[string]int{}
}

// NewInMemoryStore is constructor
func NewInMemoryStore(data map[string]int) *InMemoryStore {
	return &InMemoryStore{Data: data}

}
