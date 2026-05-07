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

type snsMetrics struct {
	mu     sync.RWMutex
	Topics map[string]*snsMetric `json:"topics"`
}

type snsMetric struct {
	Calls     int64 `json:"calls"`
	Published int64 `json:"published"`
	Lists     int64 `json:"lists"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

type snsMessage struct {
	MessageID   string            `json:"messageId"`
	Topic       string            `json:"topic"`
	Subject     string            `json:"subject,omitempty"`
	Body        string            `json:"body"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	PublishedAt time.Time         `json:"publishedAt"`
}

type snsTopicStore struct {
	mu     sync.RWMutex
	Topics map[string][]snsMessage
}

var (
	snsRuntimeMetrics = &snsMetrics{Topics: map[string]*snsMetric{}}
	snsRuntimeStore   = &snsTopicStore{Topics: map[string][]snsMessage{}}
)

// SNSProvider serves an in-memory SNS-compatible topic mock.
type SNSProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.SNSScenarios
	logger    *logrus.Entry
}

func NewSNSProvider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry) *SNSProvider {
	return &SNSProvider{w: w, req: req, logger: lgr}
}

func (s *SNSProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &s.scenarios)
}

func (s SNSProvider) Validate() error {
	if len(s.scenarios) == 0 {
		return fmt.Errorf("sns scenarios are required")
	}
	for index, scenario := range s.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("sns scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (s *SNSProvider) Retrieve(scenarioNumber int) {
	scenario, err := s.scenarioByNumber(scenarioNumber)
	if err != nil {
		s.recordSNSMetric("unknown", snsMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusInternalServerError, "sns scenario failed")
		return
	}

	topic := s.topicName(scenario)
	if topic == "" {
		s.recordSNSMetric("unknown", snsMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "sns topic is required")
		return
	}
	s.recordSNSMetric(topic, snsMetric{Calls: 1})

	if scenario.DryRun {
		s.recordSNSMetric(topic, snsMetric{DryRuns: 1})
		s.writeJSON(http.StatusOK, map[string]interface{}{"sns": true, "dryRun": true, "topic": topic})
		return
	}

	switch s.req.Method {
	case http.MethodPost:
		s.publish(topic, scenario)
	case http.MethodGet:
		s.listMessages(topic)
	default:
		s.recordSNSMetric(topic, snsMetric{Errors: 1})
		s.writeJSONError(http.StatusMethodNotAllowed, "sns method is not supported")
	}
}

func (s *SNSProvider) publish(topic string, scenario scenarios.SNSScenario) {
	body := scenario.Message
	if s.req.Body != nil {
		data, err := io.ReadAll(s.req.Body)
		if err != nil {
			s.recordSNSMetric(topic, snsMetric{Errors: 1})
			s.writeJSONError(http.StatusBadRequest, "failed to read sns message body")
			return
		}
		if len(data) > 0 {
			body = string(data)
		}
	}
	message := snsMessage{
		MessageID:   fmt.Sprintf("sns-%d", time.Now().UnixNano()),
		Topic:       topic,
		Subject:     scenario.Subject,
		Body:        body,
		Attributes:  scenario.Attributes,
		PublishedAt: time.Now().UTC(),
	}
	snsRuntimeStore.add(topic, message)
	s.recordSNSMetric(topic, snsMetric{Published: 1})
	s.writeJSON(http.StatusOK, message)
}

func (s *SNSProvider) listMessages(topic string) {
	messages := snsRuntimeStore.list(topic)
	s.recordSNSMetric(topic, snsMetric{Lists: 1})
	s.writeJSON(http.StatusOK, map[string]interface{}{"messages": messages})
}

func (s SNSProvider) scenarioByNumber(scenarioNumber int) (scenarios.SNSScenario, error) {
	if len(s.scenarios) == 0 {
		return scenarios.SNSScenario{}, fmt.Errorf("sns scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(s.scenarios) {
		return s.scenarios[0], nil
	}
	return s.scenarios[scenarioNumber], nil
}

func (s SNSProvider) topicName(scenario scenarios.SNSScenario) string {
	topic := scenario.Topic
	if s.req == nil {
		return topic
	}
	parts := strings.Split(strings.TrimPrefix(s.req.URL.Path, "/"), "/")
	if len(parts) >= 3 && parts[0] == "aws" && (parts[1] == "sns" || parts[1] == "sns-dry-run") {
		topic = defaultString(topic, parts[2])
	}
	return topic
}

func (s *SNSProvider) writeJSON(statusCode int, data interface{}) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(statusCode)
	json.NewEncoder(s.w).Encode(data)
}

func (s *SNSProvider) writeJSONError(statusCode int, message string) {
	s.writeJSON(statusCode, map[string]string{"error": message})
}

func (s SNSProvider) recordSNSMetric(topic string, delta snsMetric) {
	recordSNSMetric(topic, delta)
}

func recordSNSMetric(topic string, delta snsMetric) {
	if topic == "" {
		topic = "unknown"
	}
	snsRuntimeMetrics.mu.Lock()
	defer snsRuntimeMetrics.mu.Unlock()

	metric, ok := snsRuntimeMetrics.Topics[topic]
	if !ok {
		metric = &snsMetric{}
		snsRuntimeMetrics.Topics[topic] = metric
	}
	metric.Calls += delta.Calls
	metric.Published += delta.Published
	metric.Lists += delta.Lists
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// SNSMetricsSnapshot returns a copy of SNS runtime metrics.
func SNSMetricsSnapshot() map[string]snsMetric {
	snsRuntimeMetrics.mu.RLock()
	defer snsRuntimeMetrics.mu.RUnlock()

	out := make(map[string]snsMetric, len(snsRuntimeMetrics.Topics))
	for topic, metric := range snsRuntimeMetrics.Topics {
		out[topic] = *metric
	}
	return out
}

func (s *snsTopicStore) add(topic string, message snsMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Topics[topic] = append(s.Topics[topic], message)
}

func (s *snsTopicStore) list(topic string) []snsMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := s.Topics[topic]
	out := make([]snsMessage, len(messages))
	copy(out, messages)
	return out
}
