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

// Init
func (Y YAMLProvider) Init() error {
	panic("implement me")
}

// Parse
func (Y YAMLProvider) Parse(filePath string) (*scenarios.GigaMockScenario, error) {
	scenario := &scenarios.GigaMockScenario{}
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return scenario, err
	}
	//fmt.Println("YAML FILE==>", string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &scenario)
	if err != nil {
		return scenario, err
	}

	return scenario, nil
}
