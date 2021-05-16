package webhookType

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// TypeProvider
type TypeProvider interface {
	// Unmarshal is a parse
	Unmarshal([]map[string]interface{}) error
	// Send
	Send() error
	// Validate
	validation.Validatable
}
