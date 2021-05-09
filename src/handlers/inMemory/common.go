package inMemory

import (
	"net/http"
)

// writeResponseHeaderJson
func writeResponseHeaderJson(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}
