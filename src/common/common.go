package common

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	// HTTPScenarioType represents type "http" for a scenario
	HTTPScenarioType = "http"
	// GraphQLScenarioType represents type "graphql" for a scenario
	GraphQLScenarioType = "graphql"
)

var (
	// ScenarioTypeValidator is a validator rule for the type of a scenario
	// can be http or graphql
	// in the feature should add kafka and grpc
	ScenarioTypeValidator = []validation.Rule{
		validation.Required,
		validation.In(HTTPScenarioType, GraphQLScenarioType),
	}
	// ScenarioMethodValidator is a validator rule for the method type of a scenario
	// must be any type of HTTP methods
	ScenarioMethodValidator = []validation.Rule{
		validation.Required,
		validation.In(http.MethodPost, http.MethodGet, http.MethodPut,
			http.MethodConnect, http.MethodDelete, http.MethodHead,
			http.MethodOptions, http.MethodPatch, http.MethodTrace),
	}
	// ScenarioPathValidator is a validator rule for the method type of a scenario
	// must be required
	ScenarioPathValidator = []validation.Rule{
		validation.Required,
	}
	// BaseScenariosValidator is a validator rule for a base validation of the field scenarios
	// must be required
	BaseScenariosValidator = []validation.Rule{
		validation.Required,
	}
	// CodeStatus is a validator rule for the field codeStatus
	// must be required
	CodeStatusValidator = []validation.Rule{
		validation.Required,
		validation.Min(http.StatusOK),
		validation.Max(http.StatusNetworkAuthenticationRequired),
	}
)
