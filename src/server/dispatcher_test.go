package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	urlrouter "github.com/azer/url-router"
	"github.com/sirupsen/logrus"
)

func TestDispatcherUnknownRouteReturnsNotFoundWithoutPanic(t *testing.T) {
	dispatcher := NewDispatcher(
		context.Background(),
		nil,
		urlrouter.New(),
		logrus.NewEntry(logrus.New()),
		GRPCServerConfig{},
	)
	req := httptest.NewRequest(http.MethodGet, "/not-indexed", nil)
	recorder := httptest.NewRecorder()

	dispatcher.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}
