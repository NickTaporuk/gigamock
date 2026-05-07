package scenarioType

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func TestWebSocketProviderRetrieveDryRun(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/ws/chat", nil)
	recorder := httptest.NewRecorder()
	provider := NewWebSocketProvider(recorder, request, logrus.NewEntry(logrus.New()), context.Background())
	provider.scenarios = scenarios.WebSocketScenarios{
		{
			Name:   "dry run",
			DryRun: true,
			Steps: []scenarios.WebSocketStep{
				{Send: &scenarios.WebSocketMessage{Text: "hello"}},
			},
		},
	}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("response must be valid JSON: %v", err)
	}
	if response["websocket"] != true {
		t.Fatalf("websocket must be true: %#v", response["websocket"])
	}
	if response["dryRun"] != true {
		t.Fatalf("dryRun must be true: %#v", response["dryRun"])
	}
}

func TestWebSocketProviderRetrieveScript(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		provider := NewWebSocketProvider(w, r, logger, context.Background())
		provider.scenarios = scenarios.WebSocketScenarios{
			{
				Name: "chat",
				SendOnConnect: []scenarios.WebSocketMessage{
					{Text: "connected"},
				},
				Steps: []scenarios.WebSocketStep{
					{Receive: &scenarios.WebSocketMessage{Text: "ping"}},
					{Send: &scenarios.WebSocketMessage{Text: "pong"}},
					{Close: &scenarios.WebSocketClose{Code: websocket.CloseNormalClosure, Reason: "done"}},
				},
			},
		}
		provider.Retrieve(0)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read connect message: %v", err)
	}
	if string(message) != "connected" {
		t.Fatalf("unexpected connect message: %q", string(message))
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
		t.Fatalf("write ping: %v", err)
	}

	_, message, err = conn.ReadMessage()
	if err != nil {
		t.Fatalf("read pong: %v", err)
	}
	if string(message) != "pong" {
		t.Fatalf("unexpected pong message: %q", string(message))
	}
}

func TestWebSocketProviderValidateRequiresScenarios(t *testing.T) {
	provider := WebSocketProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func TestWebSocketProviderValidateScenarioFields(t *testing.T) {
	provider := WebSocketProvider{
		scenarios: scenarios.WebSocketScenarios{
			{},
		},
	}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty websocket scenario")
	}
}
