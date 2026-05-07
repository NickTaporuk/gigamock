package scenarioType

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

type graphQLRequestPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

type graphQLMetrics struct {
	mu        sync.RWMutex
	Endpoints map[string]*graphQLMetric `json:"endpoints"`
}

type graphQLMetric struct {
	Calls        int64 `json:"calls"`
	Errors       int64 `json:"errors"`
	Matched      int64 `json:"matched"`
	NotMatched   int64 `json:"notMatched"`
	BatchCalls   int64 `json:"batchCalls"`
	InvalidJSON  int64 `json:"invalidJson"`
	ActiveForced int64 `json:"activeForced"`
}

var graphQLRuntimeMetrics = &graphQLMetrics{Endpoints: map[string]*graphQLMetric{}}

// GraphQLTypeProvider
type GraphQLTypeProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.GraphQLScenarios `xml:"scenarios",json:"scenarios"`
}

// Validate validates decoded GraphQL scenarios.
func (g GraphQLTypeProvider) Validate() error {
	if len(g.scenarios) == 0 {
		return fmt.Errorf("scenarios is required")
	}

	for index, scenario := range g.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("graphql scenario %d is invalid: %w", index, err)
		}
		if scenario.Response.Body != "" && !json.Valid([]byte(scenario.Response.Body)) {
			return fmt.Errorf("graphql scenario %d response.body must be valid JSON", index)
		}
	}

	return nil
}

// NewGraphQLTypeProvider
func NewGraphQLTypeProvider(w http.ResponseWriter, req *http.Request) *GraphQLTypeProvider {
	return &GraphQLTypeProvider{w: w, req: req}
}

func (g *GraphQLTypeProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	err := mapstructure.Decode(rawScenarios, &g.scenarios)
	if err != nil {
		return err
	}

	return nil
}

func (g *GraphQLTypeProvider) Retrieve(scenarioNumber int) {
	endpoint := graphQLEndpoint(g.req)
	failed := false
	matched := false
	notMatched := false
	batchCall := false
	invalidJSON := false
	activeForced := false
	defer func() {
		recordGraphQLMetric(endpoint, graphQLMetric{
			Calls:        1,
			Errors:       boolToInt64(failed),
			Matched:      boolToInt64(matched),
			NotMatched:   boolToInt64(notMatched),
			BatchCalls:   boolToInt64(batchCall),
			InvalidJSON:  boolToInt64(invalidJSON),
			ActiveForced: boolToInt64(activeForced),
		})
	}()

	if len(g.scenarios) == 0 {
		g.w.WriteHeader(http.StatusNotFound)
		return
	}

	payloads, err := g.requestPayloads()
	if err != nil {
		invalidJSON = true
		failed = true
		g.writeJSONError(http.StatusBadRequest, "invalid GraphQL JSON request body")
		return
	}
	if len(payloads) > 1 {
		batchCall = true
	}

	responses := make([]json.RawMessage, 0, len(payloads))
	for _, payload := range payloads {
		scenario := g.scenarioByNumber(scenarioNumber)
		if scenarioNumber > 0 {
			activeForced = true
		} else if !graphQLRequestMatches(g.req, scenario.Request, payload) {
			if matchedScenario, ok := g.matchScenario(payload); ok {
				scenario = matchedScenario
				matched = true
			} else {
				notMatched = true
			}
		} else {
			matched = true
		}

		rawBody := scenario.Response.Body
		if rawBody == "" {
			rawBody = `{}`
		}
		responses = append(responses, json.RawMessage(rawBody))
		if len(payloads) == 1 {
			g.writeResponse(scenario, rawBody)
			return
		}
	}

	g.writeBatchResponse(responses)
}

func (g *GraphQLTypeProvider) scenarioByNumber(scenarioNumber int) scenarios.GraphQLScenario {
	if scenarioNumber < 0 || scenarioNumber >= len(g.scenarios) {
		return g.scenarios[0]
	}

	return g.scenarios[scenarioNumber]
}

func (g *GraphQLTypeProvider) matchScenario(payload *graphQLRequestPayload) (scenarios.GraphQLScenario, bool) {
	for _, scenario := range g.scenarios {
		if graphQLRequestMatches(g.req, scenario.Request, payload) {
			return scenario, true
		}
	}

	return scenarios.GraphQLScenario{}, false
}

func (g *GraphQLTypeProvider) writeResponse(scenario scenarios.GraphQLScenario, body string) {
	if len(scenario.Response.Headers) > 0 {
		for headerName, headerValue := range scenario.Response.Headers {
			g.w.Header().Add(headerName, headerValue)
		}
	}

	statusCode := http.StatusOK
	if scenario.Response.StatusCode > 0 {
		statusCode = scenario.Response.StatusCode
	}

	g.w.WriteHeader(statusCode)
	g.w.Write([]byte(body))
}

func (g *GraphQLTypeProvider) writeBatchResponse(responses []json.RawMessage) {
	g.w.Header().Set("Content-Type", "application/json")
	g.w.WriteHeader(http.StatusOK)
	json.NewEncoder(g.w).Encode(responses)
}

func (g *GraphQLTypeProvider) writeJSONError(statusCode int, message string) {
	g.w.Header().Set("Content-Type", "application/json")
	g.w.WriteHeader(statusCode)
	json.NewEncoder(g.w).Encode(map[string]interface{}{
		"errors": []map[string]string{{"message": message}},
	})
}

func (g *GraphQLTypeProvider) requestPayloads() ([]*graphQLRequestPayload, error) {
	if g.req == nil || g.req.Body == nil {
		return []*graphQLRequestPayload{{}}, nil
	}

	var raw json.RawMessage
	if err := json.NewDecoder(g.req.Body).Decode(&raw); err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return []*graphQLRequestPayload{{}}, nil
	}
	if raw[0] == '[' {
		payloads := []*graphQLRequestPayload{}
		if err := json.Unmarshal(raw, &payloads); err != nil {
			return nil, err
		}
		if len(payloads) == 0 {
			return []*graphQLRequestPayload{{}}, nil
		}
		return payloads, nil
	}

	payload := &graphQLRequestPayload{}
	if err := json.Unmarshal(raw, payload); err != nil {
		return nil, err
	}
	return []*graphQLRequestPayload{payload}, nil
}

func graphQLRequestMatches(req *http.Request, expected scenarios.GraphQLScenarioRequest, actual *graphQLRequestPayload) bool {
	if len(expected.Headers) > 0 && !headersMatch(req, expected.Headers) {
		return false
	}
	return graphQLPayloadMatches(expected, actual)
}

func graphQLPayloadMatches(expected scenarios.GraphQLScenarioRequest, actual *graphQLRequestPayload) bool {
	if expected.OperationName != "" && expected.OperationName != actual.OperationName {
		return false
	}

	if expected.Query != "" && normalizeGraphQLQuery(expected.Query) != normalizeGraphQLQuery(actual.Query) {
		return false
	}

	if len(expected.Variables) > 0 && !jsonPayloadEqual(expected.Variables, actual.Variables) {
		return false
	}

	if expected.OperationName == "" && expected.Query == "" && len(expected.Variables) == 0 {
		return true
	}

	return true
}

func normalizeGraphQLQuery(query string) string {
	return strings.Join(strings.Fields(query), " ")
}

func jsonPayloadEqual(expected interface{}, actual interface{}) bool {
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		return false
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		return false
	}

	return string(expectedJSON) == string(actualJSON)
}

func headersMatch(req *http.Request, expected map[string]string) bool {
	if req == nil {
		return len(expected) == 0
	}
	for key, expectedValue := range expected {
		if req.Header.Get(key) != expectedValue {
			return false
		}
	}
	return true
}

func graphQLEndpoint(req *http.Request) string {
	if req == nil {
		return "unknown"
	}
	return req.URL.Path
}

func boolToInt64(value bool) int64 {
	if value {
		return 1
	}
	return 0
}

func recordGraphQLMetric(endpoint string, delta graphQLMetric) {
	graphQLRuntimeMetrics.mu.Lock()
	defer graphQLRuntimeMetrics.mu.Unlock()

	metric, ok := graphQLRuntimeMetrics.Endpoints[endpoint]
	if !ok {
		metric = &graphQLMetric{}
		graphQLRuntimeMetrics.Endpoints[endpoint] = metric
	}
	metric.Calls += delta.Calls
	metric.Errors += delta.Errors
	metric.Matched += delta.Matched
	metric.NotMatched += delta.NotMatched
	metric.BatchCalls += delta.BatchCalls
	metric.InvalidJSON += delta.InvalidJSON
	metric.ActiveForced += delta.ActiveForced
}

// GraphQLMetricsSnapshot returns a copy of GraphQL runtime metrics.
func GraphQLMetricsSnapshot() map[string]graphQLMetric {
	graphQLRuntimeMetrics.mu.RLock()
	defer graphQLRuntimeMetrics.mu.RUnlock()

	out := make(map[string]graphQLMetric, len(graphQLRuntimeMetrics.Endpoints))
	for endpoint, metric := range graphQLRuntimeMetrics.Endpoints {
		out[endpoint] = *metric
	}
	return out
}
