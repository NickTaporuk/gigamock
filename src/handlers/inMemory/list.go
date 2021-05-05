package inMemory

import (
	"encoding/json"
	"net/http"
)


func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	marshaledData, err := json.Marshal(h.store)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Conent-Type", "application/json")
	w.Write(marshaledData)
}
