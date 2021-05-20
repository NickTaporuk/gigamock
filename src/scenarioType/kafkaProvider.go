package scenarioType

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

var runnedConsumers map[string]bool

type KafkaProvider struct {
	w         http.ResponseWriter
	scenarios scenarios.KafkaScenarios
	logger    *logrus.Entry
}

func NewKafkaProvider(w http.ResponseWriter, lgr *logrus.Entry) *KafkaProvider {
	return &KafkaProvider{w: w, logger: lgr}
}

func (k *KafkaProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	err := mapstructure.Decode(rawScenarios, &k.scenarios)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProvider) Retrieve(scenarioNumber int) {
	// if we have a type producer
	// I should to send a message some times
	// else if I have the type equal consumer
	// I should to show in console a message which was retrieved from kafka server

	scenario := k.scenarios[scenarioNumber]

	if scenario.Producer != nil {
		err := k.prepareTopic(&scenario)
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
					"action":   "k.prepareTopic(&scenario)",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("k.prepareTopic is retrieved an error")
			k.retrieveErrorResponse()

			return
		}

		msg := kafka.Message{}
		k.prepareMessage(scenario, &msg)
		k.prepareHeaders(scenario, &msg)

		err = k.writeMessage(scenario, msg)
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
					"action":   "k.writeMessage",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("k.writeMessage is retrieved an error")
			k.retrieveErrorResponse()

			return
		}

		//TODO: add logic to stop gorutine with consumer
		if scenario.Consumer != nil {
			go func() {
				if runnedConsumers == nil {
					runnedConsumers = map[string]interface{}{}
				}

				if _, ok := runnedConsumers[scenario.Topic]; !ok {
					err = k.prepareTopic(&scenario)
					if err != nil {
						k.logger.
							WithError(err).
							WithFields(logrus.Fields{
								"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
								"action":   "k.prepareTopic(&scenario)",
								"stack":    string(debug.Stack()),
								"scenario": scenario,
							}).Error("k.prepareTopic is retrieved an error")

						return
					}
					// make a new reader that consumes from topic-A
					r := kafka.NewReader(kafka.ReaderConfig{
						Brokers:  []string{k.host(&scenario)},
						GroupID:  "consumer-group-id",
						Topic:    scenario.Topic,
						MinBytes: 10e3, // 10KB
						MaxBytes: 10e6, // 10MB
					})

					for {
						m, err := r.ReadMessage(context.Background())
						if err != nil {
							break
						}

						k.logger.Info(fmt.Sprintf("message at topic/partition/offset/headers %v/%v/%v: %s = %s, %v\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value), m.Headers))
					}

					if err := r.Close(); err != nil {
						k.logger.
							WithError(err).
							WithFields(logrus.Fields{
								"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
								"action":   "kafka.DialLeader",
								"stack":    string(debug.Stack()),
								"scenario": scenario,
							}).Fatal("failed to close reader")
					}
					runnedConsumers[scenario.Topic] = true
				}
			}()
		}

		k.w.WriteHeader(http.StatusOK)
	}
}

func (k *KafkaProvider) prepareTopic(scenario *scenarios.KafkaScenario) error {
	conn, err := kafka.Dial("tcp", k.host(scenario))
	if err != nil {
		k.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"method":   "prepareTopic(scenario *scenarios.KafkaScenario) error",
				"action":   "kafka.Dial",
				"stack":    string(debug.Stack()),
				"scenario": scenario,
			}).Error("kafka.Dial failed to dial")

		return err
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "prepareTopic(scenario *scenarios.KafkaScenario) error",
					"action":   "conn.Close()",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("conn.Close failed to dial")
			return
		}
	}()

	exist := k.listTopics(conn, scenario.Topic)

	if !exist {
		err = k.createTopic(conn, scenario.Topic)
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "prepareTopic(scenario *scenarios.KafkaScenario) error",
					"action":   "k.createTopic(conn, scenario.Topic)",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("k.createTopic failed to dial")

			return err
		}
	}

	return nil
}

// host is a host name preparation
func (k *KafkaProvider) host(scenario *scenarios.KafkaScenario) string {
	return scenario.Host + ":" + scenario.Port
}

// writeMessage
func (k *KafkaProvider) writeMessage(scenario scenarios.KafkaScenario, msg kafka.Message) error {

	w := &kafka.Writer{
		Addr:  kafka.TCP(k.host(&scenario)),
		Topic: scenario.Topic,
		// NOTE: When Topic is not defined here, each Message must define it instead.
		Balancer: &kafka.LeastBytes{},
	}

	err := w.WriteMessages(context.Background(),
		// NOTE: Each Message has Topic defined, otherwise an error is returned.
		msg,
	)

	if err != nil {
		k.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
				"action":   "w.WriteMessages",
				"stack":    string(debug.Stack()),
				"scenario": scenario,
			}).Error("w.WriteMessages is retrieved an error")

		return err
	}
	if err := w.Close(); err != nil {
		k.logger.
			WithError(err).
			WithFields(logrus.Fields{
				"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
				"action":   "w.Close()",
				"stack":    string(debug.Stack()),
				"scenario": scenario,
			}).Error("w.Close failed to close writer")
		k.retrieveErrorResponse()

		return err
	}

	return nil
}

// prepareMessage presents a flow to put all needed values to a message
func (k *KafkaProvider) prepareMessage(scenario scenarios.KafkaScenario, msg *kafka.Message) {
	msg.Value = []byte(scenario.Producer.Message.Value)
	if scenario.Producer.Message.Key != "" {
		msg.Key = []byte(scenario.Producer.Message.Key)
	}
}

func (k *KafkaProvider) prepareHeaders(scenario scenarios.KafkaScenario, msg *kafka.Message) {
	if len(scenario.Producer.Headers) > 0 {
		for headerKey, headerValue := range scenario.Producer.Headers {
			msg.Headers = append(msg.Headers, kafka.Header{Key: headerKey, Value: []byte(headerValue)})
		}
	}

}

// list topics retrieves exist a topic or not
func (k *KafkaProvider) listTopics(conn *kafka.Conn, topic string) bool {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}

	for k := range m {
		if k == topic {
			return true
		}
	}

	return false
}

func (k *KafkaProvider) createTopic(conn *kafka.Conn, topic string) error {
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}

	return nil
}

func (k *KafkaProvider) retrieveErrorResponse() {
	k.w.WriteHeader(http.StatusInternalServerError)
}
func (k KafkaProvider) Validate() error {
	return validation.ValidateStruct(
		&k,
		validation.Field(&k.scenarios),
	)
}
