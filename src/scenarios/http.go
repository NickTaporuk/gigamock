package scenarios

import (
	"errors"
	"net/http"

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
	Redirect   *HTTPRedirect     `yaml:"redirect,omitempty"`
}

// HTTPRedirect describes an HTTP redirect response target.
type HTTPRedirect struct {
	URL                   string            `yaml:"url"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
}

// Validate validates HTTPScenarioResponse
func (hsr HTTPScenarioResponse) Validate() error {
	if hsr.Redirect == nil {
		return validation.ValidateStruct(
			&hsr,
			validation.Field(
				&hsr.StatusCode, common.CodeStatusValidator...,
			),
		)
	}

	statusCode := hsr.StatusCode
	if statusCode == 0 {
		statusCode = http.StatusFound
	}

	if statusCode < http.StatusMultipleChoices || statusCode >= http.StatusBadRequest {
		return errors.New("redirect response statusCode must be a 3xx code")
	}

	return validation.ValidateStruct(
		hsr.Redirect,
		validation.Field(&hsr.Redirect.URL, validation.Required),
	)
}
