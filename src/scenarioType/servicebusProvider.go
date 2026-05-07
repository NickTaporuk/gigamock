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

type serviceBusMetrics struct {
	mu       sync.RWMutex
	Entities map[string]*serviceBusMetric `json:"entities"`
}

type serviceBusMetric struct {
	Calls     int64 `json:"calls"`
	Sent      int64 `json:"sent"`
	Received  int64 `json:"received"`
	Completed int64 `json:"completed"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

type serviceBusMessage struct {
	MessageID    string            `json:"messageId"`
	Queue        string            `json:"queue,omitempty"`
	Topic        string            `json:"topic,omitempty"`
	Subscription string            `json:"subscription,omitempty"`
	Body         string            `json:"body"`
	Properties   map[string]string `json:"properties,omitempty"`
	EnqueuedAt   time.Time         `json:"enqueuedAt"`
	LockToken    string            `json:"lockToken,omitempty"`
}

type serviceBusStore struct {
	mu       sync.Mutex
	Entities map[string][]serviceBusMessage
}

var (
	serviceBusRuntimeMetrics = &serviceBusMetrics{Entities: map[string]*serviceBusMetric{}}
	serviceBusRuntimeStore   = &serviceBusStore{Entities: map[string][]serviceBusMessage{}}
)

// ServiceBusProvider serves an in-memory Azure Service Bus-compatible mock.
type ServiceBusProvider struct {
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.ServiceBusScenarios
	logger    *logrus.Entry
}

func NewServiceBusProvider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry) *ServiceBusProvider {
	return &ServiceBusProvider{w: w, req: req, logger: lgr}
}

func (s *ServiceBusProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &s.scenarios)
}

func (s ServiceBusProvider) Validate() error {
	if len(s.scenarios) == 0 {
		return fmt.Errorf("servicebus scenarios are required")
	}
	for index, scenario := range s.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("servicebus scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (s *ServiceBusProvider) Retrieve(scenarioNumber int) {
	scenario, err := s.scenarioByNumber(scenarioNumber)
	if err != nil {
		s.recordServiceBusMetric("unknown", serviceBusMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusInternalServerError, "servicebus scenario failed")
		return
	}

	entity, action := s.routeParts(scenario)
	if entity == "" {
		s.recordServiceBusMetric("unknown", serviceBusMetric{Calls: 1, Errors: 1})
		s.writeJSONError(http.StatusBadRequest, "servicebus queue or topic is required")
		return
	}
	s.recordServiceBusMetric(entity, serviceBusMetric{Calls: 1})

	if scenario.DryRun {
		s.recordServiceBusMetric(entity, serviceBusMetric{DryRuns: 1})
		s.writeJSON(http.StatusOK, map[string]interface{}{
			"servicebus": true,
			"dryRun":     true,
			"entity":     entity,
			"action":     action,
		})
		return
	}

	switch {
	case s.req.Method == http.MethodPost && action == "send":
		s.send(entity, scenario)
	case s.req.Method == http.MethodPost && action == "receive":
		s.receive(entity, scenario)
	case s.req.Method == http.MethodDelete && action == "purge":
		s.purge(entity)
	default:
		s.recordServiceBusMetric(entity, serviceBusMetric{Errors: 1})
		s.writeJSONError(http.StatusMethodNotAllowed, "servicebus method or action is not supported")
	}
}

func (s *ServiceBusProvider) send(entity string, scenario scenarios.ServiceBusScenario) {
	body := scenario.Message
	if s.req.Body != nil {
		data, err := io.ReadAll(s.req.Body)
		if err != nil {
			s.recordServiceBusMetric(entity, serviceBusMetric{Errors: 1})
			s.writeJSONError(http.StatusBadRequest, "failed to read servicebus message body")
			return
		}
		if len(data) > 0 {
			body = string(data)
		}
	}
	message := serviceBusMessage{
		MessageID:  fmt.Sprintf("sb-%d", time.Now().UnixNano()),
		Queue:      scenario.Queue,
		Topic:      scenario.Topic,
		Body:       body,
		Properties: scenario.Properties,
		EnqueuedAt: time.Now().UTC(),
	}
	if message.Queue == "" && message.Topic == "" {
		message.Queue = entity
	}
	serviceBusRuntimeStore.push(entity, message)
	s.recordServiceBusMetric(entity, serviceBusMetric{Sent: 1})
	s.writeJSON(http.StatusOK, message)
}

func (s *ServiceBusProvider) receive(entity string, scenario scenarios.ServiceBusScenario) {
	message, ok := serviceBusRuntimeStore.pop(entity)
	if !ok {
		s.writeJSON(http.StatusOK, map[string]interface{}{"messages": []serviceBusMessage{}})
		return
	}
	message.Subscription = scenario.Subscription
	message.LockToken = fmt.Sprintf("lock-%d", time.Now().UnixNano())
	s.recordServiceBusMetric(entity, serviceBusMetric{Received: 1})
	s.writeJSON(http.StatusOK, map[string]interface{}{"messages": []serviceBusMessage{message}})
}

func (s *ServiceBusProvider) purge(entity string) {
	deleted := serviceBusRuntimeStore.purge(entity)
	s.recordServiceBusMetric(entity, serviceBusMetric{Completed: int64(deleted)})
	s.writeJSON(http.StatusOK, map[string]interface{}{"entity": entity, "deleted": deleted})
}

func (s ServiceBusProvider) scenarioByNumber(scenarioNumber int) (scenarios.ServiceBusScenario, error) {
	if len(s.scenarios) == 0 {
		return scenarios.ServiceBusScenario{}, fmt.Errorf("servicebus scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(s.scenarios) {
		return s.scenarios[0], nil
	}
	return s.scenarios[scenarioNumber], nil
}

func (s ServiceBusProvider) routeParts(scenario scenarios.ServiceBusScenario) (string, string) {
	entity := defaultString(scenario.Queue, scenario.Topic)
	action := ""
	if s.req == nil {
		return entity, action
	}
	parts := strings.Split(strings.TrimPrefix(s.req.URL.Path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "azure" && (parts[1] == "servicebus" || parts[1] == "servicebus-dry-run") {
		entity = defaultString(entity, parts[2])
		action = parts[3]
	}
	return entity, action
}

func (s *ServiceBusProvider) writeJSON(statusCode int, data interface{}) {
	s.w.Header().Set("Content-Type", "application/json")
	s.w.WriteHeader(statusCode)
	json.NewEncoder(s.w).Encode(data)
}

func (s *ServiceBusProvider) writeJSONError(statusCode int, message string) {
	s.writeJSON(statusCode, map[string]string{"error": message})
}

func (s ServiceBusProvider) recordServiceBusMetric(entity string, delta serviceBusMetric) {
	recordServiceBusMetric(entity, delta)
}

func recordServiceBusMetric(entity string, delta serviceBusMetric) {
	if entity == "" {
		entity = "unknown"
	}
	serviceBusRuntimeMetrics.mu.Lock()
	defer serviceBusRuntimeMetrics.mu.Unlock()

	metric, ok := serviceBusRuntimeMetrics.Entities[entity]
	if !ok {
		metric = &serviceBusMetric{}
		serviceBusRuntimeMetrics.Entities[entity] = metric
	}
	metric.Calls += delta.Calls
	metric.Sent += delta.Sent
	metric.Received += delta.Received
	metric.Completed += delta.Completed
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// ServiceBusMetricsSnapshot returns a copy of Azure Service Bus runtime metrics.
func ServiceBusMetricsSnapshot() map[string]serviceBusMetric {
	serviceBusRuntimeMetrics.mu.RLock()
	defer serviceBusRuntimeMetrics.mu.RUnlock()

	out := make(map[string]serviceBusMetric, len(serviceBusRuntimeMetrics.Entities))
	for entity, metric := range serviceBusRuntimeMetrics.Entities {
		out[entity] = *metric
	}
	return out
}

func (s *serviceBusStore) push(entity string, message serviceBusMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entities[entity] = append(s.Entities[entity], message)
}

func (s *serviceBusStore) pop(entity string) (serviceBusMessage, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	messages := s.Entities[entity]
	if len(messages) == 0 {
		return serviceBusMessage{}, false
	}
	message := messages[0]
	s.Entities[entity] = messages[1:]
	return message, true
}

func (s *serviceBusStore) purge(entity string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	deleted := len(s.Entities[entity])
	s.Entities[entity] = nil
	return deleted
}
