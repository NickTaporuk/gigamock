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
		rawBody = filterGraphQLResponseBody(rawBody, payload.Query)
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

type graphQLSelection map[string]graphQLSelection

type graphQLSelectionParser struct {
	tokens []string
	pos    int
}

func filterGraphQLResponseBody(body string, query string) string {
	if strings.TrimSpace(query) == "" {
		return body
	}

	selection, ok := parseGraphQLSelection(query)
	if !ok || len(selection) == 0 {
		return body
	}

	var response map[string]interface{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return body
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		return body
	}

	filteredData, matched := filterGraphQLObject(data, selection)
	if !matched {
		return body
	}

	response["data"] = filteredData
	filteredBody, err := json.Marshal(response)
	if err != nil {
		return body
	}

	return string(filteredBody)
}

func parseGraphQLSelection(query string) (graphQLSelection, bool) {
	tokens := tokenizeGraphQLQuery(query)
	if len(tokens) == 0 {
		return nil, false
	}

	parser := &graphQLSelectionParser{tokens: tokens}
	switch parser.peek() {
	case "query", "mutation", "subscription":
		parser.next()
		operationName := ""
		if parser.isName(parser.peek()) {
			operationName = parser.next()
		}
		hasOperationArguments := parser.peek() == "("
		if hasOperationArguments {
			parser.skipBalanced("(", ")")
		}
		parser.skipDirectives()
		if parser.peek() != "{" {
			return nil, false
		}
		selection, ok := parser.parseSelectionSet()
		if hasOperationArguments && operationName != "" && len(selection) > 0 && !strings.Contains(query, "$") {
			return graphQLSelection{operationName: selection}, ok
		}
		return selection, ok
	case "{":
		return parser.parseSelectionSet()
	default:
		if !parser.isName(parser.peek()) {
			return nil, false
		}
		key, child, ok := parser.parseField()
		if !ok {
			return nil, false
		}
		return graphQLSelection{key: child}, true
	}
}

func tokenizeGraphQLQuery(query string) []string {
	tokens := []string{}
	for i := 0; i < len(query); {
		ch := query[i]
		if ch == '#' {
			for i < len(query) && query[i] != '\n' && query[i] != '\r' {
				i++
			}
			continue
		}
		if ch == '"' {
			i++
			for i < len(query) {
				if query[i] == '\\' {
					i += 2
					continue
				}
				if query[i] == '"' {
					i++
					break
				}
				i++
			}
			continue
		}
		if isGraphQLNameStart(ch) {
			start := i
			i++
			for i < len(query) && isGraphQLNameContinue(query[i]) {
				i++
			}
			tokens = append(tokens, query[start:i])
			continue
		}
		if strings.ContainsRune("{}():![]=@", rune(ch)) {
			tokens = append(tokens, string(ch))
		}
		if ch == '.' && i+2 < len(query) && query[i+1] == '.' && query[i+2] == '.' {
			tokens = append(tokens, "...")
			i += 3
			continue
		}
		i++
	}
	return tokens
}

func isGraphQLNameStart(ch byte) bool {
	return ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isGraphQLNameContinue(ch byte) bool {
	return isGraphQLNameStart(ch) || (ch >= '0' && ch <= '9')
}

func (p *graphQLSelectionParser) parseSelectionSet() (graphQLSelection, bool) {
	if p.next() != "{" {
		return nil, false
	}

	selection := graphQLSelection{}
	for p.pos < len(p.tokens) && p.peek() != "}" {
		if p.peek() == "..." {
			p.next()
			if p.peek() == "on" {
				p.next()
				if p.isName(p.peek()) {
					p.next()
				}
				p.skipDirectives()
				if p.peek() == "{" {
					inlineSelection, ok := p.parseSelectionSet()
					if !ok {
						return nil, false
					}
					mergeGraphQLSelections(selection, inlineSelection)
				}
			} else if p.isName(p.peek()) {
				p.next()
			}
			continue
		}

		key, child, ok := p.parseField()
		if !ok {
			return nil, false
		}
		selection[key] = child
	}

	if p.next() != "}" {
		return nil, false
	}
	return selection, true
}

func (p *graphQLSelectionParser) parseField() (string, graphQLSelection, bool) {
	if !p.isName(p.peek()) {
		return "", nil, false
	}

	responseKey := p.next()
	if p.peek() == ":" {
		p.next()
		if !p.isName(p.peek()) {
			return "", nil, false
		}
		p.next()
	}

	if p.peek() == "(" {
		p.skipBalanced("(", ")")
	}
	p.skipDirectives()

	if p.peek() == "{" {
		child, ok := p.parseSelectionSet()
		return responseKey, child, ok
	}

	return responseKey, graphQLSelection{}, true
}

func (p *graphQLSelectionParser) skipDirectives() {
	for p.peek() == "@" {
		p.next()
		if p.isName(p.peek()) {
			p.next()
		}
		if p.peek() == "(" {
			p.skipBalanced("(", ")")
		}
	}
}

func (p *graphQLSelectionParser) skipBalanced(open string, close string) {
	if p.peek() != open {
		return
	}

	depth := 0
	for p.pos < len(p.tokens) {
		token := p.next()
		if token == open {
			depth++
		}
		if token == close {
			depth--
			if depth == 0 {
				return
			}
		}
	}
}

func (p *graphQLSelectionParser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *graphQLSelectionParser) next() string {
	token := p.peek()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return token
}

func (p *graphQLSelectionParser) isName(token string) bool {
	return token != "" && isGraphQLNameStart(token[0])
}

func mergeGraphQLSelections(dst graphQLSelection, src graphQLSelection) {
	for key, child := range src {
		dst[key] = child
	}
}

func filterGraphQLObject(source map[string]interface{}, selection graphQLSelection) (map[string]interface{}, bool) {
	filtered := map[string]interface{}{}
	for key, childSelection := range selection {
		value, ok := source[key]
		if !ok {
			continue
		}
		filtered[key] = filterGraphQLValue(value, childSelection)
	}
	return filtered, len(filtered) > 0
}

func filterGraphQLValue(value interface{}, selection graphQLSelection) interface{} {
	if len(selection) == 0 {
		return value
	}

	switch typedValue := value.(type) {
	case map[string]interface{}:
		filtered, matched := filterGraphQLObject(typedValue, selection)
		if !matched {
			return typedValue
		}
		return filtered
	case []interface{}:
		filtered := make([]interface{}, 0, len(typedValue))
		for _, item := range typedValue {
			filtered = append(filtered, filterGraphQLValue(item, selection))
		}
		return filtered
	default:
		return value
	}
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
