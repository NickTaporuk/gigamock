package inMemory

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/NickTaporuk/gigamock/src/fileProvider"
	"github.com/NickTaporuk/gigamock/src/fileType"
	"github.com/NickTaporuk/gigamock/src/fileWalkers"
)

// ScenarioDetailsResponse describes all indexed mock endpoints for the UI.
type ScenarioDetailsResponse struct {
	Endpoints []EndpointDetails `json:"endpoints"`
}

// EndpointDetails is a single mock endpoint with its selectable scenarios.
type EndpointDetails struct {
	Key             string           `json:"key"`
	Path            string           `json:"path"`
	Method          string           `json:"method"`
	Type            string           `json:"type"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description,omitempty"`
	FilePath        string           `json:"filePath"`
	FileName        string           `json:"fileName"`
	Directory       string           `json:"directory"`
	Service         string           `json:"service"`
	CurrentScenario int              `json:"currentScenario"`
	Scenarios       []ScenarioOption `json:"scenarios"`
}

// ScenarioOption is a selectable scenario parsed from a mock config file.
type ScenarioOption struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// Details returns a UI-friendly list of config files and available scenarios.
func (h *Handler) Details(w http.ResponseWriter, r *http.Request) {
	resp, err := h.detailsResponse()
	writeResponseHeaderJson(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	marshaledData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(marshaledData)
}

func (h *Handler) detailsResponse() (*ScenarioDetailsResponse, error) {
	store := *h.store
	keys := make([]string, 0, len(store))
	for key := range store {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	endpoints := make([]EndpointDetails, 0, len(keys))
	for _, key := range keys {
		indexedData := store[key]
		ext, err := fileType.FileExtensionDetection(indexedData.FilePath)
		if err != nil {
			return nil, err
		}

		provider, err := fileProvider.Factory(ext, h.lgr)
		if err != nil {
			return nil, err
		}

		scenario, err := provider.Unmarshal(indexedData.FilePath)
		if err != nil {
			return nil, err
		}
		directory := filepath.Dir(indexedData.FilePath)

		endpoints = append(endpoints, EndpointDetails{
			Key:             fileWalkers.PrepareInMemoryStoreKey(scenario.Path, scenario.Method),
			Path:            scenario.Path,
			Method:          scenario.Method,
			Type:            scenario.Type,
			Name:            scenario.Name,
			Description:     scenario.Description,
			FilePath:        indexedData.FilePath,
			FileName:        filepath.Base(indexedData.FilePath),
			Directory:       directory,
			Service:         filepath.Base(directory),
			CurrentScenario: indexedData.ScenarioNumber,
			Scenarios:       scenarioOptions(scenario.Scenarios),
		})
	}

	return &ScenarioDetailsResponse{Endpoints: endpoints}, nil
}

func scenarioOptions(rawScenarios []map[string]interface{}) []ScenarioOption {
	options := make([]ScenarioOption, 0, len(rawScenarios))
	for index, rawScenario := range rawScenarios {
		name := fmt.Sprintf("Scenario %d", index)
		if rawName, ok := rawScenario["name"]; ok && rawName != nil {
			name = fmt.Sprint(rawName)
		}

		options = append(options, ScenarioOption{
			Index: index,
			Name:  name,
		})
	}

	return options
}
