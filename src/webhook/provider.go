package webhook

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// WebHookTypeProvider
type WebHookTypeProvider interface {
	// Unmarshal
	Unmarshal([]map[string]interface{}) error
	// Send
	Send() error
	validation.Validatable
}
