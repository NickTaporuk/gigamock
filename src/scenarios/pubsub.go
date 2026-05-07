package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// PubSubScenario describes an in-memory Google Pub/Sub-compatible mock scenario.
type PubSubScenario struct {
	Name         string
	Topic        string
	Subscription string
	DryRun       bool
	Message      string
	Attributes   map[string]string
}

func (s PubSubScenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Topic),
		validation.Field(&s.Subscription),
	)
}

// PubSubScenarios is a list of Pub/Sub scenarios.
type PubSubScenarios []PubSubScenario
