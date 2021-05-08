package fileProvider

import (
	"encoding/json"
	"github.com/NickTaporuk/gigamock/src/scenarios"
	"io/ioutil"
)

// JSONProvider
type JSONProvider struct {}

// Unmarshal
func (j *JSONProvider) Unmarshal(filePath string) (*scenarios.BaseGigaMockScenario, error) {
	scenario := &scenarios.BaseGigaMockScenario{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return scenario, err
	}

	err = json.Unmarshal(yamlFile, &scenario)
	if err != nil {
		return scenario, err
	}

	return scenario, nil
}

// NewJSONProvider
func NewJSONProvider() *JSONProvider {
	return &JSONProvider{}
}
