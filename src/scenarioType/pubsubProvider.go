package scenarioType

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type pubSubMetrics struct {
	mu     sync.RWMutex
	Topics map[string]*pubSubMetric `json:"topics"`
}

type pubSubMetric struct {
	Calls     int64 `json:"calls"`
	Published int64 `json:"published"`
	Pulled    int64 `json:"pulled"`
	Acked     int64 `json:"acked"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

type pubSubMessage struct {
	MessageID    string            `json:"messageId"`
	Topic        string            `json:"topic"`
	Subscription string            `json:"subscription,omitempty"`
	Data         string            `json:"data"`
	Attributes   map[string]string `json:"attributes,omitempty"`
	PublishTime  time.Time         `json:"publishTime"`
	AckID        string            `json:"ackId,omitempty"`
}

type pubSubTopicStore struct {
	mu     sync.Mutex
	Topics map[string][]pubSubMessage
}

var (
	pubSubRuntimeMetrics = &pubSubMetrics{Topics: map[string]*pubSubMetric{}}
	pubSubRuntimeStore   = &pubSubTopicStore{Topics: map[string][]pubSubMessage{}}
)

// PubSubProvider serves an in-memory Google Pub/Sub-compatible mock.
type PubSubProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.PubSubScenarios
	logger    *logrus.Entry
}

func NewPubSubProvider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry) *PubSubProvider {
	return &PubSubProvider{w: w, req: req, logger: lgr}
}

func (p *PubSubProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &p.scenarios)
}

func (p PubSubProvider) Validate() error {
	if len(p.scenarios) == 0 {
		return fmt.Errorf("pubsub scenarios are required")
	}
	for index, scenario := range p.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("pubsub scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (p *PubSubProvider) Retrieve(scenarioNumber int) {
	scenario, err := p.scenarioByNumber(scenarioNumber)
	if err != nil {
		p.recordPubSubMetric("unknown", pubSubMetric{Calls: 1, Errors: 1})
		p.writeJSONError(http.StatusInternalServerError, "pubsub scenario failed")
		return
	}

	topic, subscription, action := p.routeParts(scenario)
	metricKey := defaultString(topic, subscription)
	if metricKey == "" {
		p.recordPubSubMetric("unknown", pubSubMetric{Calls: 1, Errors: 1})
		p.writeJSONError(http.StatusBadRequest, "pubsub topic or subscription is required")
		return
	}
	p.recordPubSubMetric(metricKey, pubSubMetric{Calls: 1})

	if scenario.DryRun {
		p.recordPubSubMetric(metricKey, pubSubMetric{DryRuns: 1})
		p.writeJSON(http.StatusOK, map[string]interface{}{
			"pubsub":       true,
			"dryRun":       true,
			"topic":        topic,
			"subscription": subscription,
			"action":       action,
		})
		return
	}

	switch {
	case p.req.Method == http.MethodPost && action == "publish":
		p.publish(topic, scenario)
	case p.req.Method == http.MethodPost && action == "pull":
		p.pull(topic, subscription)
	case p.req.Method == http.MethodDelete && action == "purge":
		p.purge(topic, subscription)
	default:
		p.recordPubSubMetric(metricKey, pubSubMetric{Errors: 1})
		p.writeJSONError(http.StatusMethodNotAllowed, "pubsub method or action is not supported")
	}
}

func (p *PubSubProvider) publish(topic string, scenario scenarios.PubSubScenario) {
	body := scenario.Message
	if p.req.Body != nil {
		data, err := io.ReadAll(p.req.Body)
		if err != nil {
			p.recordPubSubMetric(topic, pubSubMetric{Errors: 1})
			p.writeJSONError(http.StatusBadRequest, "failed to read pubsub message body")
			return
		}
		if len(data) > 0 {
			body = string(data)
		}
	}
	message := pubSubMessage{
		MessageID:   fmt.Sprintf("pubsub-%d", time.Now().UnixNano()),
		Topic:       topic,
		Data:        body,
		Attributes:  scenario.Attributes,
		PublishTime: time.Now().UTC(),
	}
	pubSubRuntimeStore.add(topic, message)
	p.recordPubSubMetric(topic, pubSubMetric{Published: 1})
	p.writeJSON(http.StatusOK, map[string]interface{}{"messageIds": []string{message.MessageID}, "message": message})
}

func (p *PubSubProvider) pull(topic string, subscription string) {
	message, ok := pubSubRuntimeStore.pop(topic)
	if !ok {
		p.writeJSON(http.StatusOK, map[string]interface{}{"receivedMessages": []pubSubMessage{}})
		return
	}
	message.Subscription = subscription
	message.AckID = fmt.Sprintf("ack-%d", time.Now().UnixNano())
	p.recordPubSubMetric(topic, pubSubMetric{Pulled: 1})
	p.writeJSON(http.StatusOK, map[string]interface{}{"receivedMessages": []pubSubMessage{message}})
}

func (p *PubSubProvider) purge(topic string, subscription string) {
	deleted := pubSubRuntimeStore.purge(topic)
	p.recordPubSubMetric(topic, pubSubMetric{Acked: int64(deleted)})
	p.writeJSON(http.StatusOK, map[string]interface{}{"topic": topic, "subscription": subscription, "deleted": deleted})
}

func (p PubSubProvider) scenarioByNumber(scenarioNumber int) (scenarios.PubSubScenario, error) {
	if len(p.scenarios) == 0 {
		return scenarios.PubSubScenario{}, fmt.Errorf("pubsub scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(p.scenarios) {
		return p.scenarios[0], nil
	}
	return p.scenarios[scenarioNumber], nil
}

func (p PubSubProvider) routeParts(scenario scenarios.PubSubScenario) (string, string, string) {
	topic := scenario.Topic
	subscription := scenario.Subscription
	action := ""
	if p.req == nil {
		return topic, subscription, action
	}
	parts := strings.Split(strings.TrimPrefix(p.req.URL.Path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "gcp" && (parts[1] == "pubsub" || parts[1] == "pubsub-dry-run") {
		switch parts[3] {
		case "publish":
			topic = defaultString(topic, parts[2])
			action = "publish"
		case "pull":
			subscription = defaultString(subscription, parts[2])
			topic = defaultString(topic, subscription)
			action = "pull"
		case "purge":
			subscription = defaultString(subscription, parts[2])
			topic = defaultString(topic, subscription)
			action = "purge"
		}
	}
	return topic, subscription, action
}

func (p *PubSubProvider) writeJSON(statusCode int, data interface{}) {
	p.w.Header().Set("Content-Type", "application/json")
	p.w.WriteHeader(statusCode)
	json.NewEncoder(p.w).Encode(data)
}

func (p *PubSubProvider) writeJSONError(statusCode int, message string) {
	p.writeJSON(statusCode, map[string]string{"error": message})
}

func (p PubSubProvider) recordPubSubMetric(topic string, delta pubSubMetric) {
	recordPubSubMetric(topic, delta)
}

func recordPubSubMetric(topic string, delta pubSubMetric) {
	if topic == "" {
		topic = "unknown"
	}
	pubSubRuntimeMetrics.mu.Lock()
	defer pubSubRuntimeMetrics.mu.Unlock()

	metric, ok := pubSubRuntimeMetrics.Topics[topic]
	if !ok {
		metric = &pubSubMetric{}
		pubSubRuntimeMetrics.Topics[topic] = metric
	}
	metric.Calls += delta.Calls
	metric.Published += delta.Published
	metric.Pulled += delta.Pulled
	metric.Acked += delta.Acked
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// PubSubMetricsSnapshot returns a copy of Pub/Sub runtime metrics.
func PubSubMetricsSnapshot() map[string]pubSubMetric {
	pubSubRuntimeMetrics.mu.RLock()
	defer pubSubRuntimeMetrics.mu.RUnlock()

	out := make(map[string]pubSubMetric, len(pubSubRuntimeMetrics.Topics))
	for topic, metric := range pubSubRuntimeMetrics.Topics {
		out[topic] = *metric
	}
	return out
}

func (s *pubSubTopicStore) add(topic string, message pubSubMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Topics[topic] = append(s.Topics[topic], message)
}

func (s *pubSubTopicStore) pop(topic string) (pubSubMessage, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	messages := s.Topics[topic]
	if len(messages) == 0 {
		return pubSubMessage{}, false
	}
	message := messages[0]
	s.Topics[topic] = messages[1:]
	return message, true
}

func (s *pubSubTopicStore) purge(topic string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	deleted := len(s.Topics[topic])
	s.Topics[topic] = nil
	return deleted
}
