package inMemory

import (
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"github.com/sirupsen/logrus"
)

// Handler is a in-memory store handler
type Handler struct {
	store *map[string]fileWalkers.IndexedData
	lgr   *logrus.Entry
}

// NewHandler is a constructor for the struct inMemory.Handler
func NewHandler(data *map[string]fileWalkers.IndexedData, lgr ...*logrus.Entry) *Handler {
	logger := logrus.NewEntry(logrus.New())
	if len(lgr) > 0 && lgr[0] != nil {
		logger = lgr[0]
	}

	return &Handler{store: data, lgr: logger}
}
