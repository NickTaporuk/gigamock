package scenarios

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

// GigaMockHTTPScenario
type GigaMockHTTPScenario struct {
	Scenarios []HTTPScenario `yaml:"scenarios"`
}

// Scenario
type HTTPScenario struct {
	Request  HTTPScenarioRequest  `yaml:"request"`
	Response HTTPScenarioResponse `yaml:"response"`
	Delay    string               `yaml:"delay,omitempty"`
}

func (hs HTTPScenario) Validate() error {
	return validation.ValidateStruct(
		&hs,
		validation.Field(&hs.Response),
		validation.Field(&hs.Request),
	)
}

// HTTPScenarios
type HTTPScenarios []HTTPScenario

// HTTPScenarioRequest
type HTTPScenarioRequest struct {
	Headers               map[string]string `yaml:"headers,omitempty"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
	Cookies               map[string]string `yaml:"cookies,omitempty"`
}

// HTTPScenarioResponse
type HTTPScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
	Cookies    map[string]string `yaml:"cookies,omitempty"`
}

func (hsr HTTPScenarioResponse) Validate() error {
	return validation.ValidateStruct(
		&hsr,
		validation.Field(
			&hsr.StatusCode,
			common.CodeStatusValidator...,
		),

	)
}
