package fileProvider

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// YAMLProvider
type YAMLProvider struct{}

func NewYAMLProvider() *YAMLProvider {
	return &YAMLProvider{}
}

// Parse
func (y YAMLProvider) Unmarshal(filePath string) (*scenarios.BaseGigaMockScenario, error) {
	scenario := &scenarios.BaseGigaMockScenario{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return scenario, err
	}

	err = yaml.Unmarshal(yamlFile, &scenario)
	if err != nil {
		return scenario, err
	}

	return scenario, nil
}
