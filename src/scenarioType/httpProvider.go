package scenarioType

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// HTTPTypeProvider
type HTTPTypeProvider struct {
	w         http.ResponseWriter
	scenarios scenarios.HTTPScenarios
}

func checkEachScenario(scenarious interface{}) error {
	return nil
}

// Validate
func (hp HTTPTypeProvider) Validate() error {
	// need to validate request
	// need to validate response
	// validate delay
	// validate webhook

	err := validation.ValidateStruct(
		&hp,
		validation.Field(&hp.scenarios),
	)

	if err != nil {
		return err
	}

	return nil

}

// NewHTTPTypeProvider
func NewHTTPTypeProvider(w http.ResponseWriter) *HTTPTypeProvider {
	return &HTTPTypeProvider{w: w}
}

func (hp *HTTPTypeProvider) Unmarshal(rawScenario []map[string]interface{}) error {
	err := mapstructure.Decode(rawScenario, &hp.scenarios)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve
func (hp *HTTPTypeProvider) Retrieve(scenarioNumber int) {
	var body string
	var statusCode int
	if len(hp.scenarios) > 0 {
		body = hp.scenarios[scenarioNumber].Response.Body
		if hp.scenarios[scenarioNumber].Response.StatusCode > 0 {
			statusCode = hp.scenarios[scenarioNumber].Response.StatusCode
		} else {
			statusCode = http.StatusOK
		}

		if len(hp.scenarios[scenarioNumber].Response.Headers) > 0 {
			for headerName, headerValue := range hp.scenarios[scenarioNumber].Response.Headers {
				hp.w.Header().Add(headerName, headerValue)
			}
		}
	}

	hp.w.WriteHeader(statusCode)
	hp.w.Write([]byte(body))
}
