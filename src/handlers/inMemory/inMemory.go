package inMemory

import (
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
)

// Handler is a in-memory store handler
type Handler struct {
	store *map[string]fileWalkers.IndexedData
}

// NewHandler is a constructor for the struct inMemory.Handler
func NewHandler(data *map[string]fileWalkers.IndexedData) *Handler {
	return &Handler{store: data}
}
