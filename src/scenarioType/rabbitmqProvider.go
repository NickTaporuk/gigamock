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
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type rabbitMQMetrics struct {
	mu     sync.RWMutex
	Routes map[string]*rabbitMQMetric `json:"routes"`
}

type rabbitMQMetric struct {
	Calls     int64 `json:"calls"`
	Published int64 `json:"published"`
	Errors    int64 `json:"errors"`
	DryRuns   int64 `json:"dryRuns"`
}

var rabbitMQRuntimeMetrics = &rabbitMQMetrics{Routes: map[string]*rabbitMQMetric{}}

// RabbitMQProvider publishes configured messages to RabbitMQ.
type RabbitMQProvider struct {
	ctx       context.Context
	w         http.ResponseWriter
	scenarios scenarios.RabbitMQScenarios
	logger    *logrus.Entry
}

func NewRabbitMQProvider(w http.ResponseWriter, lgr *logrus.Entry, ctx context.Context) *RabbitMQProvider {
	return &RabbitMQProvider{w: w, logger: lgr, ctx: ctx}
}

func (r *RabbitMQProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &r.scenarios)
}

func (r RabbitMQProvider) Validate() error {
	if len(r.scenarios) == 0 {
		return fmt.Errorf("rabbitmq scenarios are required")
	}
	for index, scenario := range r.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("rabbitmq scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (r *RabbitMQProvider) Retrieve(scenarioNumber int) {
	started := time.Now()
	scenario, err := r.scenarioByNumber(scenarioNumber)
	if err != nil {
		r.recordRabbitMQMetric("unknown", rabbitMQMetric{Calls: 1, Errors: 1})
		r.logger.WithError(err).Error("rabbitmq scenario selection failed")
		r.writeErrorResponse()
		return
	}
	route := rabbitMQRouteKey(scenario)
	r.recordRabbitMQMetric(route, rabbitMQMetric{Calls: 1})
	defer func() {
		r.logger.WithFields(logrus.Fields{
			"exchange":   scenario.Exchange,
			"routingKey": scenario.RoutingKey,
			"duration":   time.Since(started).String(),
		}).Info("RabbitMQ scenario completed")
	}()

	if scenario.DryRun {
		r.recordRabbitMQMetric(route, rabbitMQMetric{Published: 1, DryRuns: 1})
		r.writeSuccessResponse(scenario, true)
		return
	}

	conn, err := amqp.Dial(scenario.URL)
	if err != nil {
		r.recordRabbitMQMetric(route, rabbitMQMetric{Errors: 1})
		r.logger.WithError(err).WithField("scenario", scenario).Error("rabbitmq dial failed")
		r.writeErrorResponse()
		return
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		r.recordRabbitMQMetric(route, rabbitMQMetric{Errors: 1})
		r.logger.WithError(err).WithField("scenario", scenario).Error("rabbitmq channel failed")
		r.writeErrorResponse()
		return
	}
	defer channel.Close()

	headers := amqp.Table{}
	for key, value := range scenario.Headers {
		headers[key] = value
	}

	err = channel.PublishWithContext(
		r.context(),
		scenario.Exchange,
		scenario.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: scenario.Message.ContentType,
			Body:        []byte(scenario.Message.Body),
			Headers:     headers,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		r.recordRabbitMQMetric(route, rabbitMQMetric{Errors: 1})
		r.logger.WithError(err).WithField("scenario", scenario).Error("rabbitmq publish failed")
		r.writeErrorResponse()
		return
	}

	r.recordRabbitMQMetric(route, rabbitMQMetric{Published: 1})
	r.writeSuccessResponse(scenario, false)
}

func (r RabbitMQProvider) scenarioByNumber(scenarioNumber int) (scenarios.RabbitMQScenario, error) {
	if len(r.scenarios) == 0 {
		return scenarios.RabbitMQScenario{}, fmt.Errorf("rabbitmq scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(r.scenarios) {
		return r.scenarios[0], nil
	}
	return r.scenarios[scenarioNumber], nil
}

func (r RabbitMQProvider) context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}

func (r *RabbitMQProvider) writeSuccessResponse(scenario scenarios.RabbitMQScenario, dryRun bool) {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(http.StatusOK)
	json.NewEncoder(r.w).Encode(map[string]interface{}{
		"exchange":   scenario.Exchange,
		"routingKey": scenario.RoutingKey,
		"published":  true,
		"dryRun":     dryRun,
	})
}

func (r *RabbitMQProvider) writeErrorResponse() {
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(r.w).Encode(map[string]string{"error": "rabbitmq scenario failed"})
}

func (r RabbitMQProvider) recordRabbitMQMetric(route string, delta rabbitMQMetric) {
	recordRabbitMQMetric(route, delta)
}

func recordRabbitMQMetric(route string, delta rabbitMQMetric) {
	if route == "" {
		route = "unknown"
	}
	rabbitMQRuntimeMetrics.mu.Lock()
	defer rabbitMQRuntimeMetrics.mu.Unlock()

	metric, ok := rabbitMQRuntimeMetrics.Routes[route]
	if !ok {
		metric = &rabbitMQMetric{}
		rabbitMQRuntimeMetrics.Routes[route] = metric
	}
	metric.Calls += delta.Calls
	metric.Published += delta.Published
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

func rabbitMQRouteKey(scenario scenarios.RabbitMQScenario) string {
	return scenario.Exchange + "/" + scenario.RoutingKey
}

// RabbitMQMetricsSnapshot returns a copy of RabbitMQ runtime metrics.
func RabbitMQMetricsSnapshot() map[string]rabbitMQMetric {
	rabbitMQRuntimeMetrics.mu.RLock()
	defer rabbitMQRuntimeMetrics.mu.RUnlock()

	out := make(map[string]rabbitMQMetric, len(rabbitMQRuntimeMetrics.Routes))
	for route, metric := range rabbitMQRuntimeMetrics.Routes {
		out[route] = *metric
	}
	return out
}
