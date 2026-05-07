package scenarioType

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestNATSProviderRetrieveDryRun(t *testing.T) {
	recorder := httptest.NewRecorder()
	provider := NewNATSProvider(recorder, logrus.NewEntry(logrus.New()), context.Background())
	provider.scenarios = scenarios.NATSScenarios{
		{
			Name:    "dry run",
			Host:    "nats://localhost:4222",
			Subject: "orders.created",
			DryRun:  true,
			Message: scenarios.NATSScenarioMessage{Body: `{"ok":true}`},
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
	if response["subject"] != "orders.created" {
		t.Fatalf("unexpected subject: %#v", response["subject"])
	}
	if response["dryRun"] != true {
		t.Fatalf("dryRun must be true: %#v", response["dryRun"])
	}
}

func TestNATSProviderValidateRequiresScenarios(t *testing.T) {
	provider := NATSProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func TestNATSProviderValidateScenarioFields(t *testing.T) {
	provider := NATSProvider{
		scenarios: scenarios.NATSScenarios{
			{
				Subject: "missing-host-and-body",
			},
		},
	}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for missing host and message body")
	}
}
