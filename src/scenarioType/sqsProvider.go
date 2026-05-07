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

type sqsMetrics struct {
	mu     sync.RWMutex
	Queues map[string]*sqsMetric `json:"queues"`
}

type sqsMetric struct {
	Calls    int64 `json:"calls"`
	Sent     int64 `json:"sent"`
	Received int64 `json:"received"`
	Deleted  int64 `json:"deleted"`
	Errors   int64 `json:"errors"`
	DryRuns  int64 `json:"dryRuns"`
}

type sqsMessage struct {
	MessageID     string            `json:"messageId"`
	ReceiptHandle string            `json:"receiptHandle"`
	Body          string            `json:"body"`
	Attributes    map[string]string `json:"attributes,omitempty"`
	SentAt        time.Time         `json:"sentAt"`
}

type sqsQueueStore struct {
	mu     sync.Mutex
	Queues map[string][]sqsMessage
}

var (
	sqsRuntimeMetrics = &sqsMetrics{Queues: map[string]*sqsMetric{}}
	sqsRuntimeStore   = &sqsQueueStore{Queues: map[string][]sqsMessage{}}
)

// SQSProvider serves an in-memory SQS-compatible queue mock.
type SQSProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.SQSScenarios
	logger    *logrus.Entry
}

func NewSQSProvider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry) *SQSProvider {
	return &SQSProvider{w: w, req: req, logger: lgr}
}

func (s *SQSProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &s.scenarios)
}

func (s SQSProvider) Validate() error {
	if len(s.scenarios) == 0 {
		return fmt.Errorf("sqs scenarios are required")
	}
	for index, scenario := range s.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("sqs scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (s *SQSProvider) Retrieve(scenarioNumber int) {
	scenario, err := s.scenarioByNumber(scenarioNumber)
	if err != nil {
		s.recordSQSMetric("unknown", sqsMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusInternalServerError, "sqs scenario failed")
		return
	}

	queue := s.queueName(scenario)
	if queue == "" {
		s.recordSQSMetric("unknown", sqsMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "sqs queue is required")
		return
	}
	s.recordSQSMetric(queue, sqsMetric{Calls: 1})

	if scenario.DryRun {
		s.recordSQSMetric(queue, sqsMetric{DryRuns: 1})
		s.writeJSON(http.StatusOK, map[string]interface{}{"sqs": true, "dryRun": true, "queue": queue})
		return
	}

	switch s.req.Method {
	case http.MethodPost:
		s.sendMessage(queue, scenario)
	case http.MethodGet:
		s.receiveMessage(queue)
	case http.MethodDelete:
		s.purgeQueue(queue)
	default:
		s.recordSQSMetric(queue, sqsMetric{Errors: 1})
		s.writeJSONError(http.StatusMethodNotAllowed, "sqs method is not supported")
	}
}

func (s *SQSProvider) sendMessage(queue string, scenario scenarios.SQSScenario) {
	body := scenario.Message
	if s.req.Body != nil {
		data, err := io.ReadAll(s.req.Body)
		if err != nil {
			s.recordSQSMetric(queue, sqsMetric{Errors: 1})
			s.writeJSONError(http.StatusBadRequest, "failed to read sqs message body")
			return
		}
		if len(data) > 0 {
			body = string(data)
		}
	}
	message := sqsMessage{
		MessageID:     fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		ReceiptHandle: fmt.Sprintf("receipt-%d", time.Now().UnixNano()),
		Body:          body,
		Attributes:    scenario.Attributes,
		SentAt:        time.Now().UTC(),
	}
	sqsRuntimeStore.push(queue, message)
	s.recordSQSMetric(queue, sqsMetric{Sent: 1})
	s.writeJSON(http.StatusOK, message)
}

func (s *SQSProvider) receiveMessage(queue string) {
	message, ok := sqsRuntimeStore.pop(queue)
	if !ok {
		s.writeJSON(http.StatusOK, map[string]interface{}{"messages": []sqsMessage{}})
		return
	}
	s.recordSQSMetric(queue, sqsMetric{Received: 1})
	s.writeJSON(http.StatusOK, map[string]interface{}{"messages": []sqsMessage{message}})
}

func (s *SQSProvider) purgeQueue(queue string) {
	deleted := sqsRuntimeStore.purge(queue)
	s.recordSQSMetric(queue, sqsMetric{Deleted: int64(deleted)})
	s.writeJSON(http.StatusOK, map[string]interface{}{"queue": queue, "deleted": deleted})
}

func (s SQSProvider) scenarioByNumber(scenarioNumber int) (scenarios.SQSScenario, error) {
	if len(s.scenarios) == 0 {
		return scenarios.SQSScenario{}, fmt.Errorf("sqs scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(s.scenarios) {
		return s.scenarios[0], nil
	}
	return s.scenarios[scenarioNumber], nil
}

func (s SQSProvider) queueName(scenario scenarios.SQSScenario) string {
	queue := scenario.Queue
	if s.req == nil {
		return queue
	}
	parts := strings.Split(strings.TrimPrefix(s.req.URL.Path, "/"), "/")
	if len(parts) >= 3 && parts[0] == "aws" && (parts[1] == "sqs" || parts[1] == "sqs-dry-run") {
		queue = defaultString(queue, parts[2])
	}
	return queue
}

func (s *SQSProvider) writeJSON(statusCode int, data interface{}) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(statusCode)
	json.NewEncoder(s.w).Encode(data)
}

func (s *SQSProvider) writeJSONError(statusCode int, message string) {
	s.writeJSON(statusCode, map[string]string{"error": message})
}

func (s SQSProvider) recordSQSMetric(queue string, delta sqsMetric) {
	recordSQSMetric(queue, delta)
}

func recordSQSMetric(queue string, delta sqsMetric) {
	if queue == "" {
		queue = "unknown"
	}
	sqsRuntimeMetrics.mu.Lock()
	defer sqsRuntimeMetrics.mu.Unlock()

	metric, ok := sqsRuntimeMetrics.Queues[queue]
	if !ok {
		metric = &sqsMetric{}
		sqsRuntimeMetrics.Queues[queue] = metric
	}
	metric.Calls += delta.Calls
	metric.Sent += delta.Sent
	metric.Received += delta.Received
	metric.Deleted += delta.Deleted
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// SQSMetricsSnapshot returns a copy of SQS runtime metrics.
func SQSMetricsSnapshot() map[string]sqsMetric {
	sqsRuntimeMetrics.mu.RLock()
	defer sqsRuntimeMetrics.mu.RUnlock()

	out := make(map[string]sqsMetric, len(sqsRuntimeMetrics.Queues))
	for queue, metric := range sqsRuntimeMetrics.Queues {
		out[queue] = *metric
	}
	return out
}

func (s *sqsQueueStore) push(queue string, message sqsMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Queues[queue] = append(s.Queues[queue], message)
}

func (s *sqsQueueStore) pop(queue string) (sqsMessage, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	messages := s.Queues[queue]
	if len(messages) == 0 {
		return sqsMessage{}, false
	}
	message := messages[0]
	s.Queues[queue] = messages[1:]
	return message, true
}

func (s *sqsQueueStore) purge(queue string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	deleted := len(s.Queues[queue])
	s.Queues[queue] = nil
	return deleted
}
