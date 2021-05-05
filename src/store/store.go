package store

type Store interface {
	SetState(name, status string)
	State(name string) string
	Reset(name string) bool
	ResetAll()
}
