package scenarioType

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

var runnedConsumers map[string]bool

type kafkaMetrics struct {
	mu     sync.RWMutex
	Topics map[string]*kafkaMetric `json:"topics"`
}

type kafkaMetric struct {
	Calls         int64 `json:"calls"`
	Produced      int64 `json:"produced"`
	Consumed      int64 `json:"consumed"`
	Errors        int64 `json:"errors"`
	TopicsCreated int64 `json:"topicsCreated"`
}

var kafkaRuntimeMetrics = &kafkaMetrics{Topics: map[string]*kafkaMetric{}}

// KafkaProvider is an implementation of  the interface TypeProvider
// this provider is responsible to run a scenario with type kafka
type KafkaProvider struct {
	ctx       context.Context
	w         http.ResponseWriter
	scenarios scenarios.KafkaScenarios
	logger    *logrus.Entry
}

func NewKafkaProvider(
	w http.ResponseWriter,
	lgr *logrus.Entry,
	ctx context.Context,
) *KafkaProvider {
	return &KafkaProvider{w: w, logger: lgr, ctx: ctx}
}

func (k *KafkaProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	err := mapstructure.Decode(rawScenarios, &k.scenarios)
	if err != nil {
		return err
	}

	return nil
}

// Retrieve is a main function to run a scenario
func (k *KafkaProvider) Retrieve(scenarioNumber int) {
	// if we have a type producer
	// I should to send a message some times
	// else if I have the type equal consumer
	// I should to show in console a message which was retrieved from kafka server

	started := time.Now()
	scenario, err := k.scenarioByNumber(scenarioNumber)
	if err != nil {
		k.recordKafkaMetric("unknown", kafkaMetric{Calls: 1, Errors: 1})
		k.logger.WithError(err).Error("kafka scenario selection failed")
		k.retrieveErrorResponse()
		return
	}
	k.recordKafkaMetric(scenario.Topic, kafkaMetric{Calls: 1})
	defer func() {
		k.logger.WithFields(logrus.Fields{
			"topic":    scenario.Topic,
			"duration": time.Since(started).String(),
		}).Info("Kafka scenario completed")
	}()

	if scenario.Producer != nil {
		if scenario.DryRun {
			k.recordKafkaMetric(scenario.Topic, kafkaMetric{Produced: 1})
			k.writeSuccessResponse(scenario, true)
			return
		}

		err := k.prepareTopic(&scenario)
		if err != nil {
			k.recordKafkaMetric(scenario.Topic, kafkaMetric{Errors: 1})
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
			k.recordKafkaMetric(scenario.Topic, kafkaMetric{Errors: 1})
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

		//TODO: add logic to stop goroutine with consumer
		if scenario.Consumer != nil {
			go func() {
				if runnedConsumers == nil {
					runnedConsumers = map[string]bool{}
				}

				if _, ok := runnedConsumers[scenario.Topic]; !ok {
					err = k.prepareTopic(&scenario)
					if err != nil {
						k.recordKafkaMetric(scenario.Topic, kafkaMetric{Errors: 1})
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
					reader := kafka.NewReader(kafka.ReaderConfig{
						Brokers:  []string{k.host(&scenario)},
						GroupID:  "consumer-group-id",
						Topic:    scenario.Topic,
						MinBytes: 10e3, // 10KB
						MaxBytes: 10e6, // 10MB
					})

					for {
						m, err := reader.ReadMessage(k.context())
						if err != nil {
							break
						}
						k.recordKafkaMetric(scenario.Topic, kafkaMetric{Consumed: 1})

						k.logger.Info(fmt.Sprintf("message at topic/partition/offset/headers %v/%v/%v: %s = %s, %#v\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value), m.Headers))
					}

					if err := reader.Close(); err != nil {
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

		k.writeSuccessResponse(scenario, false)
	}
}

func (k KafkaProvider) scenarioByNumber(scenarioNumber int) (scenarios.KafkaScenario, error) {
	if len(k.scenarios) == 0 {
		return scenarios.KafkaScenario{}, fmt.Errorf("kafka scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(k.scenarios) {
		return k.scenarios[0], nil
	}
	return k.scenarios[scenarioNumber], nil
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
		k.recordKafkaMetric(scenario.Topic, kafkaMetric{TopicsCreated: 1})
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

	err := w.WriteMessages(k.context(),
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
	k.recordKafkaMetric(scenario.Topic, kafkaMetric{Produced: 1})
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
		k.logger.WithError(err).Error("kafka ReadPartitions failed")
		return false
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
	k.w.Header().Set("Content-Type", "application/json")
	k.w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(k.w).Encode(map[string]string{"error": "kafka scenario failed"})
}

func (k *KafkaProvider) writeSuccessResponse(scenario scenarios.KafkaScenario, dryRun bool) {
	k.w.Header().Set("Content-Type", "application/json")
	k.w.WriteHeader(http.StatusOK)
	json.NewEncoder(k.w).Encode(map[string]interface{}{
		"topic":    scenario.Topic,
		"produced": true,
		"dryRun":   dryRun,
	})
}

func (k KafkaProvider) Validate() error {
	if len(k.scenarios) == 0 {
		return fmt.Errorf("kafka scenarios are required")
	}
	for index, scenario := range k.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("kafka scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (k KafkaProvider) context() context.Context {
	if k.ctx == nil {
		return context.Background()
	}
	return k.ctx
}

func (k KafkaProvider) recordKafkaMetric(topic string, delta kafkaMetric) {
	recordKafkaMetric(topic, delta)
}

func recordKafkaMetric(topic string, delta kafkaMetric) {
	if topic == "" {
		topic = "unknown"
	}
	kafkaRuntimeMetrics.mu.Lock()
	defer kafkaRuntimeMetrics.mu.Unlock()

	metric, ok := kafkaRuntimeMetrics.Topics[topic]
	if !ok {
		metric = &kafkaMetric{}
		kafkaRuntimeMetrics.Topics[topic] = metric
	}
	metric.Calls += delta.Calls
	metric.Produced += delta.Produced
	metric.Consumed += delta.Consumed
	metric.Errors += delta.Errors
	metric.TopicsCreated += delta.TopicsCreated
}

// KafkaMetricsSnapshot returns a copy of Kafka runtime metrics.
func KafkaMetricsSnapshot() map[string]kafkaMetric {
	kafkaRuntimeMetrics.mu.RLock()
	defer kafkaRuntimeMetrics.mu.RUnlock()

	out := make(map[string]kafkaMetric, len(kafkaRuntimeMetrics.Topics))
	for topic, metric := range kafkaRuntimeMetrics.Topics {
		out[topic] = *metric
	}
	return out
}
