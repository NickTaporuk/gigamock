package fileProvider

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
	"github.com/NickTaporuk/gigamock/src/scenarios"
)

// ValidateBaseFileStruct is a base file fields validator
// This method is validate base fields path, method and type for a particular file type and format
func ValidateBaseFileStruct(scenario scenarios.BaseGigaMockScenario) error {

	return validation.ValidateStruct(
		&scenario,
		validation.Field(&scenario.Method, common.ScenarioMethodValidator...,
		),
		validation.Field(&scenario.Path,
			validation.Required),
		validation.Field(&scenario.Type,
			common.ScenarioTypeValidator...,
		),
		validation.Field(&scenario.Scenarios,
			validation.Required,
		),
	)
}
