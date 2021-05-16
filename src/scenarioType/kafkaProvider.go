package scenarioType

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

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
		conn, err := kafka.DialLeader(
			context.Background(),
			"tcp",
			scenario.Host+":"+scenario.Port,
			scenario.Topic, scenario.Producer.Partition)
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
					"action":   "kafka.DialLeader",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("kafka failed to dial leader:")
		}

		err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
					"action":   "conn.SetWriteDeadline",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("conn.SetWriteDeadline is retrieved an error")
		}

		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(scenario.Producer.Message)},
		)

		if err != nil {
			k.logger.
				WithError(err).
				WithFields(logrus.Fields{
					"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
					"action":   "conn.SetWriteDeadline",
					"stack":    string(debug.Stack()),
					"scenario": scenario,
				}).Error("conn.SetWriteDeadline is retrieved an error")
		}

		//TODO: add logic to stop gorutine with consuler
		if scenario.Consumer.CLI {
			go func() {
				conn, err := kafka.DialLeader(
					context.Background(),
					"tcp",
					scenario.Host+":"+scenario.Port,
					scenario.Topic, scenario.Producer.Partition)
				if err != nil {
					k.logger.
						WithError(err).
						WithFields(logrus.Fields{
							"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
							"action":   "kafka.DialLeader",
							"stack":    string(debug.Stack()),
							"scenario": scenario,
						}).Error("kafka failed to dial leader:")
				}

				err = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					k.logger.
						WithError(err).
						WithFields(logrus.Fields{
							"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
							"action":   "conn.SetReadDeadline(time.Now().Add(10*time.Second))",
							"stack":    string(debug.Stack()),
							"scenario": scenario,
						}).Error("kafka conn.SetReadDeadline is retrieved an error")
				}

				batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

				b := make([]byte, 10e3) // 10KB max per message
				for {
					_, err := batch.Read(b)
					if err != nil {
						break
					}

					k.logger.
						WithFields(logrus.Fields{
							"message": string(b),
						}).Debug("kafka consumer debug")
					fmt.Println(string(b))
				}

				if err := batch.Close(); err != nil {
					k.logger.
						WithError(err).
						WithFields(logrus.Fields{
							"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
							"action":   "batch.Close()",
							"stack":    string(debug.Stack()),
							"scenario": scenario,
						}).Fatal("kafka batch.Close() is failed to close batch")
				}

				if err := conn.Close(); err != nil {
					k.logger.
						WithError(err).
						WithFields(logrus.Fields{
							"method":   "func (k *KafkaProvider) Retrieve(scenarioNumber int) ",
							"action":   "conn.Close()",
							"stack":    string(debug.Stack()),
							"scenario": scenario,
						}).Fatal("kafka conn.Close() is failed to close connection")
				}
			}()
		}

		k.w.WriteHeader(http.StatusOK)
	}
}

func (k KafkaProvider) Validate() error {
	return validation.ValidateStruct(
		&k,
		validation.Field(&k.scenarios),
	)
}
