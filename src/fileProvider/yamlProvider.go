package fileProvider

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"runtime/debug"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// YAMLProvider
type YAMLProvider struct {
	logger *logrus.Entry
}

func NewYAMLProvider(lgr *logrus.Entry) *YAMLProvider {
	return &YAMLProvider{logger: lgr}
}

func (y YAMLProvider) Validate(scenario scenarios.BaseGigaMockScenario) error {
	return ValidateBaseFileStruct(scenario)
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
		y.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"trace":    string(debug.Stack()),
				"method":   "j.Unmarshal",
				"action":   "yaml.Unmarshal",
			}).
			Error("yaml unmarshal retrieved an error")
		return scenario, err
	}

	err = y.Validate(*scenario)
	if err != nil {
		y.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"scenario": scenario,
				"trace":    string(debug.Stack()),
				"method":   "j.Unmarshal",
				"action":   "yaml.Unmarshal",
			}).
			Error("yaml validation retrieved an error")
		return nil, err
	}

	return scenario, nil
}
