package scenarioType

import (
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// HTTPTypeProvider
type HTTPTypeProvider struct {
	w         http.ResponseWriter
	scenarios scenarios.HTTPScenarios
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

func (hp *HTTPTypeProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	err := mapstructure.Decode(rawScenarios, &hp.scenarios)
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
		response := hp.scenarios[scenarioNumber].Response
		body = response.Body
		if response.StatusCode > 0 {
			statusCode = response.StatusCode
		} else {
			statusCode = http.StatusOK
		}

		if len(response.Headers) > 0 {
			for headerName, headerValue := range response.Headers {
				hp.w.Header().Add(headerName, headerValue)
			}
		}

		if len(response.Cookies) > 0 {
			for cookieName, cookieValue := range response.Cookies {
				http.SetCookie(hp.w, &http.Cookie{Name: cookieName, Value: cookieValue})
			}
		}

		if response.Redirect != nil {
			if response.StatusCode == 0 {
				statusCode = http.StatusFound
			}
			hp.w.Header().Set("Location", redirectLocation(response.Redirect))
			body = ""
		}
	}

	hp.w.WriteHeader(statusCode)
	hp.w.Write([]byte(body))
}

func redirectLocation(redirect *scenarios.HTTPRedirect) string {
	if redirect == nil || len(redirect.QueryStringParameters) == 0 {
		if redirect == nil {
			return ""
		}
		return redirect.URL
	}

	location, err := url.Parse(redirect.URL)
	if err != nil {
		return redirect.URL
	}

	query := location.Query()
	for key, value := range redirect.QueryStringParameters {
		query.Set(key, value)
	}
	location.RawQuery = query.Encode()
	return location.String()
}
