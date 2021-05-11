package inMemory

import (
	"encoding/json"
	"errors"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
	"net/http"
)

var (
	errParticularRouteNotFound = errors.New("particular route isn't found")
)

// Request is used to parse request to insert new record to in-memory store
type Request struct {
	Path           string `json:"path"`
	ScenarioNumber int    `json:"scenarioNumber"`
	Method         string `json:"method"`
}

// AddRecordResponse
type AddRecordResponse struct {
	Route  string `json:"route"`
	Method string `json:"method"`
	Error  string `json:"error, omitempty"`
}

// AddRecord is a handler to create a new record or update existed one for in-memory store data
func (h *Handler) AddRecord(w http.ResponseWriter, r *http.Request) {

	req := Request{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key := fileWalkers.PrepareImMemoryStoreKey(req.Path, req.Method)
	store := *h.store
	if _, ok := store[key]; ok {
		store[key] = fileWalkers.IndexedData{
			FilePath:       store[key].FilePath,
			ScenarioNumber: req.ScenarioNumber,
		}
	} else {
		addRecordResponseError(w, req.Path, req.Method, http.StatusBadRequest, errParticularRouteNotFound)
	}

	writeResponseHeaderJson(w)
	w.WriteHeader(http.StatusOK)
}

func addRecordResponseError(
	w http.ResponseWriter,
	route string,
	method string,
	httpStatus int,
	err error,
) {
	resp := AddRecordResponse{
		Route:  route,
		Method: method,
	}

	resp.Error = err.Error()

	marshaled, errParsing := json.Marshal(resp)
	if errParsing != nil {
		resp.Error = errParsing.Error()
	}

	writeResponseHeaderJson(w)
	w.Write(marshaled)
	w.WriteHeader(httpStatus)
}
