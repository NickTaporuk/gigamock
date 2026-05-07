package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// SQSScenario describes an in-memory SQS-compatible queue mock scenario.
type SQSScenario struct {
	Name       string
	Queue      string
	DryRun     bool
	Message    string
	Attributes map[string]string
}

func (s SQSScenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Queue),
	)
}

// SQSScenarios is a list of SQS scenarios.
type SQSScenarios []SQSScenario
