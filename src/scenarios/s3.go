package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// S3Scenario describes an S3-compatible object mock scenario.
type S3Scenario struct {
	Name        string
	Bucket      string
	Key         string
	DryRun      bool
	ContentType string
	Body        string
	Metadata    map[string]string
	Headers     map[string]string
}

func (s S3Scenario) Validate() error {
	return validation.ValidateStruct(
		&s,
		validation.Field(&s.Bucket),
		validation.Field(&s.Key),
	)
}

// S3Scenarios is a list of S3 scenarios.
type S3Scenarios []S3Scenario
