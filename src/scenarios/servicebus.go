package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// ServiceBusScenario describes an in-memory Azure Service Bus-compatible mock scenario.
type ServiceBusScenario struct {
	Name         string
	Queue        string
	Topic        string
	Subscription string
	DryRun       bool
	Message      string
	Properties   map[string]string
}

func (s ServiceBusScenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Queue),
		validation.Field(&s.Topic),
		validation.Field(&s.Subscription),
	)
}

// ServiceBusScenarios is a list of Azure Service Bus scenarios.
type ServiceBusScenarios []ServiceBusScenario
