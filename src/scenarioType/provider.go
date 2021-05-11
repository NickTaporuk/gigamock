package scenarioType

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// TypeProvider
type TypeProvider interface {
	// Unmarshal
	Unmarshal([]map[string]interface{}) error
	Retrieve(scenarioNumber int)
	validation.Validatable
}
