package scenarios

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

// GigaMockHTTPScenario is a struct for parsing yaml file
type GigaMockHTTPScenario struct {
	Scenarios []HTTPScenario `yaml:"scenarios"`
}

// HTTPScenario Scenario is a struct for parsing yaml file
type HTTPScenario struct {
	Request  HTTPScenarioRequest  `yaml:"request"`
	Response HTTPScenarioResponse `yaml:"response"`
	Delay    string               `yaml:"delay,omitempty"`
}

// Validate validates HTTPScenario
func (hs HTTPScenario) Validate() error {
	return validation.ValidateStruct(
		&hs,
		validation.Field(&hs.Response),
		validation.Field(&hs.Request),
	)
}

// HTTPScenarios is a slice of HTTPScenario
type HTTPScenarios []HTTPScenario

// HTTPScenarioRequest is a struct for parsing yaml file
type HTTPScenarioRequest struct {
	Headers               map[string]string `yaml:"headers,omitempty"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
	Cookies               map[string]string `yaml:"cookies,omitempty"`
	Body                  string            `yaml:"body,omitempty"`
}

// HTTPScenarioResponse	is a struct for parsing yaml file
type HTTPScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
	Cookies    map[string]string `yaml:"cookies,omitempty"`
}

// Validate validates HTTPScenarioResponse
func (hsr HTTPScenarioResponse) Validate() error {
	return validation.ValidateStruct(
		&hsr,
		validation.Field(
			&hsr.StatusCode, common.CodeStatusValidator...,
		),
	)
}
