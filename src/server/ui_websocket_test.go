package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestMockUIWebSocketPublishesMetricsSnapshot(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dispatcher := &Dispatcher{
		ctx:         ctx,
		grpcMetrics: newGRPCMetrics(),
	}
	server := httptest.NewServer(http.HandlerFunc(dispatcher.serveMockUIWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/internal/v1/mock-ui/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial mock UI websocket: %v", err)
	}
	defer conn.Close()

	payload := map[string]interface{}{}
	if err := conn.ReadJSON(&payload); err != nil {
		t.Fatalf("read mock UI metrics payload: %v", err)
	}

	if _, ok := payload["grpc"]; !ok {
		t.Fatalf("expected grpc metrics in payload: %#v", payload)
	}
	if _, ok := payload["graphql"]; !ok {
		t.Fatalf("expected graphql metrics in payload: %#v", payload)
	}
	if _, ok := payload["soap"]; !ok {
		t.Fatalf("expected soap metrics in payload: %#v", payload)
	}
}
