package scenarioType

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/mitchellh/mapstructure"
)

type soapMetrics struct {
	mu        sync.RWMutex
	Endpoints map[string]*soapMetric `json:"endpoints"`
}

type soapMetric struct {
	Calls        int64 `json:"calls"`
	Errors       int64 `json:"errors"`
	Matched      int64 `json:"matched"`
	NotMatched   int64 `json:"notMatched"`
	ActiveForced int64 `json:"activeForced"`
}

var soapRuntimeMetrics = &soapMetrics{Endpoints: map[string]*soapMetric{}}

// SOAPProvider serves SOAP-over-HTTP XML mock responses.
type SOAPProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.SOAPScenarios
}

func NewSOAPProvider(w http.ResponseWriter, req *http.Request) *SOAPProvider {
	return &SOAPProvider{w: w, req: req}
}

func (s *SOAPProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &s.scenarios)
}

func (s SOAPProvider) Validate() error {
	if len(s.scenarios) == 0 {
		return fmt.Errorf("soap scenarios are required")
	}
	for index, scenario := range s.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("soap scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (s *SOAPProvider) Retrieve(scenarioNumber int) {
	endpoint := soapEndpoint(s.req)
	failed := false
	matched := false
	notMatched := false
	activeForced := false
	defer func() {
		recordSOAPMetric(endpoint, soapMetric{
			Calls:        1,
			Errors:       boolToInt64(failed),
			Matched:      boolToInt64(matched),
			NotMatched:   boolToInt64(notMatched),
			ActiveForced: boolToInt64(activeForced),
		})
	}()

	if len(s.scenarios) == 0 {
		failed = true
		s.writeSOAPError(http.StatusNotFound, "SOAP scenario is not found")
		return
	}

	body, err := s.requestBody()
	if err != nil {
		failed = true
		s.writeSOAPError(http.StatusBadRequest, "failed to read SOAP request body")
		return
	}

	scenario := s.scenarioByNumber(scenarioNumber)
	if scenarioNumber > 0 {
		activeForced = true
	} else if !soapRequestMatches(s.req, body, scenario.Request) {
		if matchedScenario, ok := s.matchScenario(body); ok {
			scenario = matchedScenario
			matched = true
		} else {
			notMatched = true
		}
	} else {
		matched = true
	}

	s.writeResponse(scenario)
}

func (s SOAPProvider) scenarioByNumber(scenarioNumber int) scenarios.SOAPScenario {
	if scenarioNumber < 0 || scenarioNumber >= len(s.scenarios) {
		return s.scenarios[0]
	}
	return s.scenarios[scenarioNumber]
}

func (s *SOAPProvider) matchScenario(body string) (scenarios.SOAPScenario, bool) {
	for _, scenario := range s.scenarios {
		if soapRequestMatches(s.req, body, scenario.Request) {
			return scenario, true
		}
	}
	return scenarios.SOAPScenario{}, false
}

func (s *SOAPProvider) requestBody() (string, error) {
	if s.req == nil || s.req.Body == nil {
		return "", nil
	}
	data, err := io.ReadAll(s.req.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *SOAPProvider) writeResponse(scenario scenarios.SOAPScenario) {
	for headerName, headerValue := range scenario.Response.Headers {
		s.w.Header().Set(headerName, headerValue)
	}
	if s.w.Header().Get("Content-Type") == "" {
		s.w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	}

	statusCode := http.StatusOK
	if scenario.Response.StatusCode > 0 {
		statusCode = scenario.Response.StatusCode
	}
	s.w.WriteHeader(statusCode)
	s.w.Write([]byte(scenario.Response.Body))
}

func (s *SOAPProvider) writeSOAPError(statusCode int, message string) {
	s.w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	s.w.WriteHeader(statusCode)
	fmt.Fprintf(s.w, `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><soap:Fault><faultcode>Server</faultcode><faultstring>%s</faultstring></soap:Fault></soap:Body></soap:Envelope>`, message)
}

func soapRequestMatches(req *http.Request, body string, expected scenarios.SOAPScenarioRequest) bool {
	if len(expected.Headers) > 0 && !headersMatch(req, expected.Headers) {
		return false
	}
	if expected.SOAPAction != "" && soapAction(req) != expected.SOAPAction {
		return false
	}
	if expected.BodyContains != "" && !strings.Contains(body, expected.BodyContains) {
		return false
	}
	if expected.SOAPAction == "" && expected.BodyContains == "" && len(expected.Headers) == 0 {
		return true
	}
	return true
}

func soapAction(req *http.Request) string {
	if req == nil {
		return ""
	}
	return strings.Trim(req.Header.Get("SOAPAction"), `"`)
}

func soapEndpoint(req *http.Request) string {
	if req == nil {
		return "unknown"
	}
	return req.URL.Path
}

func recordSOAPMetric(endpoint string, delta soapMetric) {
	if endpoint == "" {
		endpoint = "unknown"
	}
	soapRuntimeMetrics.mu.Lock()
	defer soapRuntimeMetrics.mu.Unlock()

	metric, ok := soapRuntimeMetrics.Endpoints[endpoint]
	if !ok {
		metric = &soapMetric{}
		soapRuntimeMetrics.Endpoints[endpoint] = metric
	}
	metric.Calls += delta.Calls
	metric.Errors += delta.Errors
	metric.Matched += delta.Matched
	metric.NotMatched += delta.NotMatched
	metric.ActiveForced += delta.ActiveForced
}

// SOAPMetricsSnapshot returns a copy of SOAP runtime metrics.
func SOAPMetricsSnapshot() map[string]soapMetric {
	soapRuntimeMetrics.mu.RLock()
	defer soapRuntimeMetrics.mu.RUnlock()

	out := make(map[string]soapMetric, len(soapRuntimeMetrics.Endpoints))
	for endpoint, metric := range soapRuntimeMetrics.Endpoints {
		out[endpoint] = *metric
	}
	return out
}

func (s *SOAPProvider) writeJSONError(statusCode int, message string) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(statusCode)
	json.NewEncoder(s.w).Encode(map[string]string{"error": message})
}
