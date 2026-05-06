package scenarios

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

// GraphQLScenario describes a single GraphQL mock scenario.
type GraphQLScenario struct {
	Request  GraphQLScenarioRequest  `yaml:"request"`
	Response GraphQLScenarioResponse `yaml:"response"`
	Delay    string                  `yaml:"delay,omitempty"`
}

// Validate validates GraphQLScenario.
func (gs GraphQLScenario) Validate() error {
	return validation.ValidateStruct(
		&gs,
		validation.Field(&gs.Response),
	)
}

// GraphQLScenarios is a slice of GraphQLScenario.
type GraphQLScenarios []GraphQLScenario

// GraphQLScenarioRequest describes optional request matching fields.
type GraphQLScenarioRequest struct {
	Headers       map[string]string      `yaml:"headers,omitempty"`
	OperationName string                 `yaml:"operationName,omitempty"`
	Query         string                 `yaml:"query,omitempty"`
	Variables     map[string]interface{} `yaml:"variables,omitempty"`
	Body          string                 `yaml:"body,omitempty"`
}

// GraphQLScenarioResponse describes a GraphQL HTTP response.
type GraphQLScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
}

// Validate validates GraphQLScenarioResponse.
func (gsr GraphQLScenarioResponse) Validate() error {
	return validation.ValidateStruct(
		&gsr,
		validation.Field(
			&gsr.StatusCode, common.CodeStatusValidator...,
		),
	)
}
