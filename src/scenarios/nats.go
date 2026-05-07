package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// NATSScenario describes a publish scenario for a NATS subject.
type NATSScenario struct {
	Name    string
	Host    string
	Subject string
	Headers map[string]string
	DryRun  bool
	Message NATSScenarioMessage
}

func (ns NATSScenario) Validate() error {
	return validation.ValidateStruct(
		&ns,
		validation.Field(&ns.Host, validation.Required),
		validation.Field(&ns.Subject, validation.Required),
		validation.Field(&ns.Message, validation.Required),
	)
}

// NATSScenarioMessage is the payload published to NATS.
type NATSScenarioMessage struct {
	Body string
}

func (nsm NATSScenarioMessage) Validate() error {
	return validation.ValidateStruct(
		&nsm,
		validation.Field(&nsm.Body, validation.Required),
	)
}

// NATSScenarios is a list of NATS scenarios.
type NATSScenarios []NATSScenario
