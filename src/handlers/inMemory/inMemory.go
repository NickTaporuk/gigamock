package inMemory

import (
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
)

type Handler struct {
	store *map[string]fileWalkers.IndexedData
}

func NewHandler(data *map[string]fileWalkers.IndexedData) *Handler {
	return &Handler{store: data}
}
