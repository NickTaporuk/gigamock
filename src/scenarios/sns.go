package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// SNSScenario describes an in-memory SNS-compatible topic mock scenario.
type SNSScenario struct {
	Name       string
	Topic      string
	DryRun     bool
	Message    string
	Subject    string
	Attributes map[string]string
}

func (s SNSScenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Topic),
	)
}

// SNSScenarios is a list of SNS scenarios.
type SNSScenarios []SNSScenario
