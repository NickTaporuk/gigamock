package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// RabbitMQScenario describes a publish scenario for RabbitMQ.
type RabbitMQScenario struct {
	Name       string
	URL        string
	Exchange   string
	RoutingKey string
	Headers    map[string]string
	DryRun     bool
	Message    RabbitMQScenarioMessage
}

func (rs RabbitMQScenario) Validate() error {
	return validation.ValidateStruct(
		&rs,
		validation.Field(&rs.URL, validation.Required),
		validation.Field(&rs.Exchange, validation.Required),
		validation.Field(&rs.RoutingKey, validation.Required),
		validation.Field(&rs.Message, validation.Required),
	)
}

// RabbitMQScenarioMessage is the payload published to RabbitMQ.
type RabbitMQScenarioMessage struct {
	ContentType string
	Body        string
}

func (rsm RabbitMQScenarioMessage) Validate() error {
	return validation.ValidateStruct(
		&rsm,
		validation.Field(&rsm.Body, validation.Required),
	)
}

// RabbitMQScenarios is a list of RabbitMQ scenarios.
type RabbitMQScenarios []RabbitMQScenario
