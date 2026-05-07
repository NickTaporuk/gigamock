package scenarioType

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

type mqttMetrics struct {
	mu     sync.RWMutex
	Topics map[string]*mqttMetric `json:"topics"`
}

type mqttMetric struct {
	Calls     int64 `json:"calls"`
	Published int64 `json:"published"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

var mqttRuntimeMetrics = &mqttMetrics{Topics: map[string]*mqttMetric{}}

// MQTTProvider publishes configured messages to MQTT.
type MQTTProvider struct {
	ctx       context.Context
	w         http.ResponseWriter
	scenarios scenarios.MQTTScenarios
	logger    *logrus.Entry
}

func NewMQTTProvider(w http.ResponseWriter, lgr *logrus.Entry, ctx context.Context) *MQTTProvider {
	return &MQTTProvider{w: w, logger: lgr, ctx: ctx}
}

func (m *MQTTProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &m.scenarios)
}

func (m MQTTProvider) Validate() error {
	if len(m.scenarios) == 0 {
		return fmt.Errorf("mqtt scenarios are required")
	}
	for index, scenario := range m.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("mqtt scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (m *MQTTProvider) Retrieve(scenarioNumber int) {
	started := time.Now()
	scenario, err := m.scenarioByNumber(scenarioNumber)
	if err != nil {
		m.recordMQTTMetric("unknown", mqttMetric{Calls: 1, Errors: 1})
		m.logger.WithError(err).Error("mqtt scenario selection failed")
		m.writeErrorResponse()
		return
	}
	m.recordMQTTMetric(scenario.Topic, mqttMetric{Calls: 1})
	defer func() {
		m.logger.WithFields(logrus.Fields{
			"topic":    scenario.Topic,
			"duration": time.Since(started).String(),
		}).Info("MQTT scenario completed")
	}()

	if scenario.DryRun {
		m.recordMQTTMetric(scenario.Topic, mqttMetric{Published: 1, DryRuns: 1})
		m.writeSuccessResponse(scenario, true)
		return
	}

	options := mqtt.NewClientOptions().
		AddBroker(scenario.Broker).
		SetClientID(scenario.ClientID).
		SetConnectTimeout(10 * time.Second)

	client := mqtt.NewClient(options)
	if token := client.Connect(); !token.WaitTimeout(10*time.Second) || token.Error() != nil {
		m.recordMQTTMetric(scenario.Topic, mqttMetric{Errors: 1})
		err := token.Error()
		if err == nil {
			err = fmt.Errorf("mqtt connect timed out")
		}
		m.logger.WithError(err).WithField("scenario", scenario).Error("mqtt connect failed")
		m.writeErrorResponse()
		return
	}
	defer client.Disconnect(250)

	token := client.Publish(scenario.Topic, scenario.QOS, scenario.Retained, scenario.Message.Body)
	if !token.WaitTimeout(10*time.Second) || token.Error() != nil {
		m.recordMQTTMetric(scenario.Topic, mqttMetric{Errors: 1})
		err := token.Error()
		if err == nil {
			err = fmt.Errorf("mqtt publish timed out")
		}
		m.logger.WithError(err).WithField("scenario", scenario).Error("mqtt publish failed")
		m.writeErrorResponse()
		return
	}

	m.recordMQTTMetric(scenario.Topic, mqttMetric{Published: 1})
	m.writeSuccessResponse(scenario, false)
}

func (m MQTTProvider) scenarioByNumber(scenarioNumber int) (scenarios.MQTTScenario, error) {
	if len(m.scenarios) == 0 {
		return scenarios.MQTTScenario{}, fmt.Errorf("mqtt scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(m.scenarios) {
		return m.scenarios[0], nil
	}
	return m.scenarios[scenarioNumber], nil
}

func (m *MQTTProvider) writeSuccessResponse(scenario scenarios.MQTTScenario, dryRun bool) {
	m.w.Header().Set("Content-Type", "application/json")
	m.w.WriteHeader(http.StatusOK)
	json.NewEncoder(m.w).Encode(map[string]interface{}{
		"topic":     scenario.Topic,
		"published": true,
		"dryRun":    dryRun,
	})
}

func (m *MQTTProvider) writeErrorResponse() {
	m.w.Header().Set("Content-Type", "application/json")
	m.w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(m.w).Encode(map[string]string{"error": "mqtt scenario failed"})
}

func (m MQTTProvider) recordMQTTMetric(topic string, delta mqttMetric) {
	recordMQTTMetric(topic, delta)
}

func recordMQTTMetric(topic string, delta mqttMetric) {
	if topic == "" {
		topic = "unknown"
	}
	mqttRuntimeMetrics.mu.Lock()
	defer mqttRuntimeMetrics.mu.Unlock()

	metric, ok := mqttRuntimeMetrics.Topics[topic]
	if !ok {
		metric = &mqttMetric{}
		mqttRuntimeMetrics.Topics[topic] = metric
	}
	metric.Calls += delta.Calls
	metric.Published += delta.Published
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// MQTTMetricsSnapshot returns a copy of MQTT runtime metrics.
func MQTTMetricsSnapshot() map[string]mqttMetric {
	mqttRuntimeMetrics.mu.RLock()
	defer mqttRuntimeMetrics.mu.RUnlock()

	out := make(map[string]mqttMetric, len(mqttRuntimeMetrics.Topics))
	for topic, metric := range mqttRuntimeMetrics.Topics {
		out[topic] = *metric
	}
	return out
}
