package inMemory

import (
	"encoding/json"
	"net/http"
)


func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	marshaledData, err := json.Marshal(h.store)
	writeResponseHeaderJson(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(marshaledData)
	w.WriteHeader(http.StatusOK)
}
