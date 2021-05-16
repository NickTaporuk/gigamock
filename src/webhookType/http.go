package webhookType

import (
	"io"
	"net/http"
	"strings"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/NickTaporuk/gigamock/src/webhook"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// HTTPProvider
type HTTPProvider struct {
	logger         *logrus.Entry
	scenarios      scenarios.WebHookHTTPScenarios
	webHook        *webhook.WebHook
	scenarioNumber int
}

// NewHTTPProvider
func NewHTTPProvider(
	lgr *logrus.Entry,
	wbh *webhook.WebHook,
	scenarioNumber int,
) *HTTPProvider {
	return &HTTPProvider{
		logger:         lgr,
		webHook:        wbh,
		scenarioNumber: scenarioNumber,
	}
}

// Unmarshal
func (h *HTTPProvider) Unmarshal(rawData []map[string]interface{}) error {
	err := mapstructure.Decode(rawData, &h.scenarios)

	if err != nil {
		return err
	}

	return nil
}

// Send
func (h *HTTPProvider) Send() error {
	client := &http.Client{}

	var body io.Reader
	scenario := h.scenarios[h.scenarioNumber]

	if scenario.Request.Body != "" {
		body = strings.NewReader(scenario.Request.Body)
	}

	req, err := http.NewRequest(h.webHook.Method, h.webHook.Path, body)
	if err != nil {
		return err
	}

	for headerKey, headerValue := range scenario.Request.Headers {
		req.Header.Add(headerKey, headerValue)
	}

	for cookieKey, cookieValue := range scenario.Request.Cookies {
		req.AddCookie(&http.Cookie{
			Name:  cookieKey,
			Value: cookieValue,
			Path:  h.webHook.Path,
		})
	}

	_, err = client.Do(req)

	return err

}

// Validate
func (h HTTPProvider) Validate() error {
	return nil
}
