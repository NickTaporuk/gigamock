package inMemory

import (
	"encoding/json"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"net/http"
)

// Request is used to parse request to insert new record to in-memory store
type Request struct {
	Path           string `json:"path"`
	ScenarioNumber int    `json:"scenarioNumber"`
	Method         string `json:"method"`
}

// AddRecord is a hanler to create a new record or update existed one for in-memory store data
func (h *Handler) AddRecord(w http.ResponseWriter, r *http.Request) {

	req := Request{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := req.Path + "|" + req.Method
	store := *h.store
	if _, ok := store[key]; ok {
		store[key] = fileWalkers.IndexedData{
			FilePath:       store[key].FilePath,
			ScenarioNumber: req.ScenarioNumber,
		}
	}

	w.WriteHeader(http.StatusOK)
}
