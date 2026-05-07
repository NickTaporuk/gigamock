package scenarios

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

// SOAPScenario describes a single SOAP-over-HTTP mock scenario.
type SOAPScenario struct {
	Name     string               `yaml:"name,omitempty"`
	Request  SOAPScenarioRequest  `yaml:"request"`
	Response SOAPScenarioResponse `yaml:"response"`
	Delay    string               `yaml:"delay,omitempty"`
}

func (s SOAPScenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Response),
	)
}

// SOAPScenarios is a list of SOAP scenarios.
type SOAPScenarios []SOAPScenario

// SOAPScenarioRequest describes optional SOAP request matching fields.
type SOAPScenarioRequest struct {
	Headers      map[string]string `yaml:"headers,omitempty"`
	SOAPAction   string            `yaml:"soapAction,omitempty"`
	BodyContains string            `yaml:"bodyContains,omitempty"`
}

// SOAPScenarioResponse describes a SOAP XML HTTP response.
type SOAPScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
}

func (s SOAPScenarioResponse) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.StatusCode, common.CodeStatusValidator...),
	)
}
