package fileProvider

import (
	"encoding/json"
	"io/ioutil"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

// JSONProvider
type JSONProvider struct {
	logger *logrus.Entry
}

// NewJSONProvider
func NewJSONProvider(lgr *logrus.Entry) *JSONProvider {
	return &JSONProvider{logger: lgr}
}

//  Validate
func (j JSONProvider) Validate(scenario scenarios.BaseGigaMockScenario) error {
	return ValidateBaseFileStruct(scenario)
}

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

	err = j.Validate(*scenario)
	if err != nil {
		return nil, err
	}

	return scenario, nil
}
