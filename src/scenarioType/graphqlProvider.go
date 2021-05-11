package scenarioType

import (
	"net/http"

	"github.com/mitchellh/mapstructure"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// GraphQLTypeProvider
type GraphQLTypeProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.HTTPScenarios `xml:"scenarios",json:"scenarios"`
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
	// need to compare post request payloads
	// to be sure we use right endpoint with particular params
	// we use the same endpoint just only different post params
	panic("implement me")
}
