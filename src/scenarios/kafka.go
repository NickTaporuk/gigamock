package scenarios

import validation "github.com/go-ozzo/ozzo-validation"

// KafkaScenario is a kafka scenario
// which should include Host, Port, Topic as a required fields
type KafkaScenario struct {
	Name  string // a name of a scenario
	Host  string // a particular host to send or consume a message
	Port  string // port for a kafka server
	Topic string // name of a topic which the app should send a message
	Delay string // time duration between sending next message

	Producer *KafkaScenarioProducer
	Consumer *KafkaScenarioConsumer
}

func (ks KafkaScenario) Validate() error {
	return validation.ValidateStruct(
		&ks,
		validation.Field(&ks.Host, validation.Required),
		validation.Field(&ks.Port, validation.Required),
		validation.Field(&ks.Topic, validation.Required),
		validation.Field(&ks.Producer),
	)
}

type KafkaScenarioProducer struct {
	Headers map[string]string
	Message string
	Partition int // kafka partition
	Retry   uint32
}

func (ksp KafkaScenarioProducer) Validate() error {
	return validation.ValidateStruct(
		&ksp,
		validation.Field(&ksp.Message, validation.Required),
	)
}

type KafkaScenarioConsumer struct {
	CLI bool
}

// KafkaScenarios is a list of scenarios
type KafkaScenarios []KafkaScenario
