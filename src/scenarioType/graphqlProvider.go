package scenarioType

import (
	"encoding/json"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

type graphQLRequestPayload struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
}

// GraphQLTypeProvider
type GraphQLTypeProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.GraphQLScenarios `xml:"scenarios",json:"scenarios"`
}

// Validate validates decoded GraphQL scenarios.
func (g GraphQLTypeProvider) Validate() error {
	return validation.ValidateStruct(
		&g,
		validation.Field(&g.scenarios),
	)
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
	if len(g.scenarios) == 0 {
		g.w.WriteHeader(http.StatusNotFound)
		return
	}

	scenario := g.scenarioByNumber(scenarioNumber)
	payload, err := g.requestPayload()
	if err == nil && !graphQLRequestMatches(scenario.Request, payload) {
		if matchedScenario, ok := g.matchScenario(payload); ok {
			scenario = matchedScenario
		}
	}

	g.writeResponse(scenario)
}

func (g *GraphQLTypeProvider) scenarioByNumber(scenarioNumber int) scenarios.GraphQLScenario {
	if scenarioNumber < 0 || scenarioNumber >= len(g.scenarios) {
		return g.scenarios[0]
	}

	return g.scenarios[scenarioNumber]
}

func (g *GraphQLTypeProvider) matchScenario(payload *graphQLRequestPayload) (scenarios.GraphQLScenario, bool) {
	for _, scenario := range g.scenarios {
		if graphQLRequestMatches(scenario.Request, payload) {
			return scenario, true
		}
	}

	return scenarios.GraphQLScenario{}, false
}

func (g *GraphQLTypeProvider) writeResponse(scenario scenarios.GraphQLScenario) {
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
	g.w.Write([]byte(scenario.Response.Body))
}

func (g *GraphQLTypeProvider) requestPayload() (*graphQLRequestPayload, error) {
	payload := &graphQLRequestPayload{}
	if g.req == nil || g.req.Body == nil {
		return payload, nil
	}

	err := json.NewDecoder(g.req.Body).Decode(payload)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

func graphQLRequestMatches(expected scenarios.GraphQLScenarioRequest, actual *graphQLRequestPayload) bool {
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
