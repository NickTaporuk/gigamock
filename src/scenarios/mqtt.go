package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// MQTTScenario describes a publish scenario for MQTT.
type MQTTScenario struct {
	Name     string
	Broker   string
	ClientID string
	Topic    string
	QOS      byte
	Retained bool
	DryRun   bool
	Message  MQTTScenarioMessage
}

func (ms MQTTScenario) Validate() error {
	return validation.ValidateStruct(
		&ms,
		validation.Field(&ms.Broker, validation.Required),
		validation.Field(&ms.ClientID, validation.Required),
		validation.Field(&ms.Topic, validation.Required),
		validation.Field(&ms.Message, validation.Required),
	)
}

// MQTTScenarioMessage is the payload published to MQTT.
type MQTTScenarioMessage struct {
	ContentType string
	Body        string
}

func (msm MQTTScenarioMessage) Validate() error {
	return validation.ValidateStruct(
		&msm,
		validation.Field(&msm.Body, validation.Required),
	)
}

// MQTTScenarios is a list of MQTT scenarios.
type MQTTScenarios []MQTTScenario
