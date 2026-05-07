package scenarioType

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type webSocketMetrics struct {
	mu        sync.RWMutex
	Endpoints map[string]*webSocketMetric `json:"endpoints"`
}

type webSocketMetric struct {
	Calls            int64 `json:"calls"`
	Connections      int64 `json:"connections"`
	MessagesSent     int64 `json:"messagesSent"`
	MessagesReceived int64 `json:"messagesReceived"`
	Closed           int64 `json:"closed"`
	Errors           int64 `json:"errors"`
	DryRuns          int64 `json:"dryRuns"`
}

var webSocketRuntimeMetrics = &webSocketMetrics{Endpoints: map[string]*webSocketMetric{}}

// WebSocketProvider serves scripted WebSocket mock scenarios.
type WebSocketProvider struct {
	ctx       context.Context
	w         http.ResponseWriter
	req       *http.Request
	scenarios scenarios.WebSocketScenarios
	logger    *logrus.Entry
}

func NewWebSocketProvider(w http.ResponseWriter, req *http.Request, lgr *logrus.Entry, ctx context.Context) *WebSocketProvider {
	return &WebSocketProvider{w: w, req: req, logger: lgr, ctx: ctx}
}

func (ws *WebSocketProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	return mapstructure.Decode(rawScenarios, &ws.scenarios)
}

func (ws WebSocketProvider) Validate() error {
	if len(ws.scenarios) == 0 {
		return fmt.Errorf("websocket scenarios are required")
	}
	for index, scenario := range ws.scenarios {
		if err := scenario.Validate(); err != nil {
			return fmt.Errorf("websocket scenario %d is invalid: %w", index, err)
		}
	}
	return nil
}

func (ws *WebSocketProvider) Retrieve(scenarioNumber int) {
	endpoint := ws.endpoint()
	started := time.Now()
	scenario, err := ws.scenarioByNumber(scenarioNumber)
	if err != nil {
		ws.recordWebSocketMetric(endpoint, webSocketMetric{Calls: 1, Errors: 1})
		ws.logger.WithError(err).Error("websocket scenario selection failed")
		ws.writeErrorResponse(http.StatusInternalServerError, "websocket scenario failed")
		return
	}
	ws.recordWebSocketMetric(endpoint, webSocketMetric{Calls: 1})
	defer func() {
		ws.logger.WithFields(logrus.Fields{
			"endpoint": endpoint,
			"duration": time.Since(started).String(),
		}).Info("WebSocket scenario completed")
	}()

	if scenario.DryRun {
		ws.recordWebSocketMetric(endpoint, webSocketMetric{DryRuns: 1})
		ws.writeDryRunResponse(scenario)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(ws.w, ws.req, nil)
	if err != nil {
		ws.recordWebSocketMetric(endpoint, webSocketMetric{Errors: 1})
		ws.logger.WithError(err).WithField("endpoint", endpoint).Error("websocket upgrade failed")
		return
	}
	defer conn.Close()
	ws.recordWebSocketMetric(endpoint, webSocketMetric{Connections: 1})

	if err := ws.runScenario(conn, endpoint, scenario); err != nil {
		ws.recordWebSocketMetric(endpoint, webSocketMetric{Errors: 1})
		ws.logger.WithError(err).WithField("endpoint", endpoint).Error("websocket scenario failed")
		return
	}
}

func (ws *WebSocketProvider) runScenario(conn *websocket.Conn, endpoint string, scenario scenarios.WebSocketScenario) error {
	for _, message := range scenario.SendOnConnect {
		if err := ws.writeMessage(conn, endpoint, message); err != nil {
			return err
		}
	}

	for _, step := range scenario.Steps {
		if step.Delay != "" {
			delay, err := time.ParseDuration(step.Delay)
			if err != nil {
				return err
			}
			timer := time.NewTimer(delay)
			select {
			case <-timer.C:
			case <-ws.context().Done():
				timer.Stop()
				return ws.context().Err()
			}
		}
		if step.Receive != nil {
			if err := ws.readExpectedMessage(conn, endpoint, *step.Receive); err != nil {
				return err
			}
		}
		if step.Send != nil {
			if err := ws.writeMessage(conn, endpoint, *step.Send); err != nil {
				return err
			}
		}
		if step.Close != nil {
			code := step.Close.Code
			if code == 0 {
				code = websocket.CloseNormalClosure
			}
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, step.Close.Reason))
			ws.recordWebSocketMetric(endpoint, webSocketMetric{Closed: 1})
			return err
		}
	}

	ws.recordWebSocketMetric(endpoint, webSocketMetric{Closed: 1})
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "scenario completed"))
}

func (ws *WebSocketProvider) readExpectedMessage(conn *websocket.Conn, endpoint string, expected scenarios.WebSocketMessage) error {
	messageType, data, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	ws.recordWebSocketMetric(endpoint, webSocketMetric{MessagesReceived: 1})

	expectedType := webSocketMessageType(expected)
	if messageType != expectedType {
		return fmt.Errorf("unexpected websocket message type: got %d, want %d", messageType, expectedType)
	}
	if string(data) != expected.Text {
		return fmt.Errorf("unexpected websocket message: got %q, want %q", string(data), expected.Text)
	}
	return nil
}

func (ws *WebSocketProvider) writeMessage(conn *websocket.Conn, endpoint string, message scenarios.WebSocketMessage) error {
	err := conn.WriteMessage(webSocketMessageType(message), []byte(message.Text))
	if err != nil {
		return err
	}
	ws.recordWebSocketMetric(endpoint, webSocketMetric{MessagesSent: 1})
	return nil
}

func webSocketMessageType(message scenarios.WebSocketMessage) int {
	if message.Type == "binary" {
		return websocket.BinaryMessage
	}
	return websocket.TextMessage
}

func (ws WebSocketProvider) scenarioByNumber(scenarioNumber int) (scenarios.WebSocketScenario, error) {
	if len(ws.scenarios) == 0 {
		return scenarios.WebSocketScenario{}, fmt.Errorf("websocket scenarios are empty")
	}
	if scenarioNumber < 0 || scenarioNumber >= len(ws.scenarios) {
		return ws.scenarios[0], nil
	}
	return ws.scenarios[scenarioNumber], nil
}

func (ws WebSocketProvider) endpoint() string {
	if ws.req == nil {
		return "unknown"
	}
	return ws.req.URL.Path
}

func (ws WebSocketProvider) context() context.Context {
	if ws.ctx == nil {
		return context.Background()
	}
	return ws.ctx
}

func (ws *WebSocketProvider) writeDryRunResponse(scenario scenarios.WebSocketScenario) {
	ws.w.Header().Set("Content-Type", "application/json")
	ws.w.WriteHeader(http.StatusOK)
	json.NewEncoder(ws.w).Encode(map[string]interface{}{
		"websocket":     true,
		"dryRun":        true,
		"sendOnConnect": len(scenario.SendOnConnect),
		"steps":         len(scenario.Steps),
	})
}

func (ws *WebSocketProvider) writeErrorResponse(statusCode int, message string) {
	ws.w.Header().Set("Content-Type", "application/json")
	ws.w.WriteHeader(statusCode)
	json.NewEncoder(ws.w).Encode(map[string]string{"error": message})
}

func (ws WebSocketProvider) recordWebSocketMetric(endpoint string, delta webSocketMetric) {
	recordWebSocketMetric(endpoint, delta)
}

func recordWebSocketMetric(endpoint string, delta webSocketMetric) {
	if endpoint == "" {
		endpoint = "unknown"
	}
	webSocketRuntimeMetrics.mu.Lock()
	defer webSocketRuntimeMetrics.mu.Unlock()

	metric, ok := webSocketRuntimeMetrics.Endpoints[endpoint]
	if !ok {
		metric = &webSocketMetric{}
		webSocketRuntimeMetrics.Endpoints[endpoint] = metric
	}
	metric.Calls += delta.Calls
	metric.Connections += delta.Connections
	metric.MessagesSent += delta.MessagesSent
	metric.MessagesReceived += delta.MessagesReceived
	metric.Closed += delta.Closed
	metric.Errors += delta.Errors
	metric.DryRuns += delta.DryRuns
}

// WebSocketMetricsSnapshot returns a copy of WebSocket runtime metrics.
func WebSocketMetricsSnapshot() map[string]webSocketMetric {
	webSocketRuntimeMetrics.mu.RLock()
	defer webSocketRuntimeMetrics.mu.RUnlock()

	out := make(map[string]webSocketMetric, len(webSocketRuntimeMetrics.Endpoints))
	for endpoint, metric := range webSocketRuntimeMetrics.Endpoints {
		out[endpoint] = *metric
	}
	return out
}
