package fileProvider

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// ValidateBaseFileStruct is a base file fields validator
// This method is validate base fields path, method and type for a particular file type and format
func ValidateBaseFileStruct(scenario scenarios.BaseGigaMockScenario) error {

	return validation.ValidateStruct(
		&scenario,
		validation.Field(&scenario.Method,
			validation.Required,
			validation.In(http.MethodPost, http.MethodGet, http.MethodPut,
				http.MethodConnect, http.MethodDelete, http.MethodHead,
				http.MethodOptions, http.MethodPatch, http.MethodTrace),
		),
		validation.Field(&scenario.Path,
			validation.Required),
		validation.Field(&scenario.Type,
			validation.Required,
			validation.In(scenarios.HTTPScenarioType, scenarios.GraphQLScenarioType),
		),
		validation.Field(&scenario.Scenarios,
			validation.Required,
		),
	)
}
