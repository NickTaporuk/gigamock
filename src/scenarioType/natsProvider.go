package scenarioType

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/mitchellh/mapstructure"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type natsMetrics struct {
	mu       sync.RWMutex
	Subjects map[string]*natsMetric `json:"subjects"`
}

type natsMetric struct {
	Calls     int64 `json:"calls"`
	Published int64 `json:"published"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

var natsRuntimeMetrics = &natsMetrics{Subjects: map[string]*natsMetric{}}

// NATSProvider publishes configured messages to a NATS subject.
type NATSProvider struct {
	ctx       context.Context
	w         http.ResponseWriter
	scenarios scenarios.NATSScenarios
	logger    *logrus.Entry
}

func NewNATSProvider(w http.ResponseWriter, lgr *logrus.Entry, ctx context.Context) *NATSProvider {
	return &NATSProvider{w: w, logger: lgr, ctx: ctx}
}

func (n *NATSProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &n.scenarios)
}

func (n NATSProvider) Validate() error {
	if len(n.scenarios) == 0 {
		return fmt.Errorf("nats scenarios are required")
	}
	for index, scenario := range n.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("nats scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (n *NATSProvider) Retrieve(scenarioNumber int) {
	started := time.Now()
	scenario, err := n.scenarioByNumber(scenarioNumber)
	if err != nil {
		n.recordNATSMetric("unknown", natsMetric{Calls: 1, Errors: 1})
		n.logger.WithError(err).Error("nats scenario selection failed")
		n.writeErrorResponse()
		return
	}
	n.recordNATSMetric(scenario.Subject, natsMetric{Calls: 1})
	defer func() {
		n.logger.WithFields(logrus.Fields{
			"subject":  scenario.Subject,
			"duration": time.Since(started).String(),
		}).Info("NATS scenario completed")
	}()

	if scenario.DryRun {
		n.recordNATSMetric(scenario.Subject, natsMetric{Published: 1, DryRuns: 1})
		n.writeSuccessResponse(scenario, true)
		return
	}

	conn, err := nats.Connect(scenario.Host)
	if err != nil {
		n.recordNATSMetric(scenario.Subject, natsMetric{Errors: 1})
		n.logger.WithError(err).WithField("scenario", scenario).Error("nats connect failed")
		n.writeErrorResponse()
		return
	}
	defer conn.Close()

	msg := &nats.Msg{
		Subject: scenario.Subject,
		Data:    []byte(scenario.Message.Body),
		Header:  nats.Header{},
	}
	for key, value := range scenario.Headers {
		msg.Header.Set(key, value)
	}

	if err := conn.PublishMsg(msg); err != nil {
		n.recordNATSMetric(scenario.Subject, natsMetric{Errors: 1})
		n.logger.WithError(err).WithField("scenario", scenario).Error("nats publish failed")
		n.writeErrorResponse()
		return
	}
	if err := conn.FlushWithContext(n.context()); err != nil {
		n.recordNATSMetric(scenario.Subject, natsMetric{Errors: 1})
		n.logger.WithError(err).WithField("scenario", scenario).Error("nats flush failed")
		n.writeErrorResponse()
		return
	}

	n.recordNATSMetric(scenario.Subject, natsMetric{Published: 1})
	n.writeSuccessResponse(scenario, false)
}

func (n NATSProvider) scenarioByNumber(scenarioNumber int) (scenarios.NATSScenario, error) {
	if len(n.scenarios) == 0 {
		return scenarios.NATSScenario{}, fmt.Errorf("nats scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(n.scenarios) {
		return n.scenarios[0], nil
	}
	return n.scenarios[scenarioNumber], nil
}

func (n NATSProvider) context() context.Context {
	if n.ctx == nil {
		return context.Background()
	}
	return n.ctx
}

func (n *NATSProvider) writeSuccessResponse(scenario scenarios.NATSScenario, dryRun bool) {
	n.w.Header().Set("Content-Type", "application/json")
	n.w.WriteHeader(http.StatusOK)
	json.NewEncoder(n.w).Encode(map[string]interface{}{
		"subject":   scenario.Subject,
		"published": true,
		"dryRun":    dryRun,
	})
}

func (n *NATSProvider) writeErrorResponse() {
	n.w.Header().Set("Content-Type", "application/json")
	n.w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(n.w).Encode(map[string]string{"error": "nats scenario failed"})
}

func (n NATSProvider) recordNATSMetric(subject string, delta natsMetric) {
	recordNATSMetric(subject, delta)
}

func recordNATSMetric(subject string, delta natsMetric) {
	if subject == "" {
		subject = "unknown"
	}
	natsRuntimeMetrics.mu.Lock()
	defer natsRuntimeMetrics.mu.Unlock()

	metric, ok := natsRuntimeMetrics.Subjects[subject]
	if !ok {
		metric = &natsMetric{}
		natsRuntimeMetrics.Subjects[subject] = metric
	}
	metric.Calls += delta.Calls
	metric.Published += delta.Published
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// NATSMetricsSnapshot returns a copy of NATS runtime metrics.
func NATSMetricsSnapshot() map[string]natsMetric {
	natsRuntimeMetrics.mu.RLock()
	defer natsRuntimeMetrics.mu.RUnlock()

	out := make(map[string]natsMetric, len(natsRuntimeMetrics.Subjects))
	for subject, metric := range natsRuntimeMetrics.Subjects {
		out[subject] = *metric
	}
	return out
}
