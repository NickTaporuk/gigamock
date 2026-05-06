package scenarioType

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// MessageBrokerTypeProvider indexes broker mock configs that do not have runtime support yet.
type MessageBrokerTypeProvider struct {
	w            http.ResponseWriter
	scenarioType string
	scenarios    []map[string]interface{}
}

// NewMessageBrokerTypeProvider creates a provider for broker mock scenario files.
func NewMessageBrokerTypeProvider(w http.ResponseWriter, scenarioType string) *MessageBrokerTypeProvider {
	return &MessageBrokerTypeProvider{
		w:            w,
		scenarioType: scenarioType,
	}
}

// Unmarshal stores raw broker scenarios so files can be indexed and shown in the UI.
func (m *MessageBrokerTypeProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	m.scenarios = rawScenarios
	return nil
}

// Validate validates decoded broker scenarios.
func (m MessageBrokerTypeProvider) Validate() error {
	return validation.Validate(m.scenarios, validation.Required)
}

// Retrieve returns a clear response if a broker mock path is called through HTTP.
func (m *MessageBrokerTypeProvider) Retrieve(scenarioNumber int) {
	if m.w == nil {
		return
	}

	m.w.Header().Set("Content-Type", "application/json")
	m.w.WriteHeader(http.StatusNotImplemented)
	m.w.Write([]byte(fmt.Sprintf(
		`{"error":"%s mock runtime is not implemented yet; this scenario is indexed for the control UI"}`,
		m.scenarioType,
	)))
}
