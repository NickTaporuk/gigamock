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

func TestMQTTProviderRetrieveDryRun(t *testing.T) {
	recorder := httptest.NewRecorder()
	provider := NewMQTTProvider(recorder, logrus.NewEntry(logrus.New()), context.Background())
	provider.scenarios = scenarios.MQTTScenarios{
		{
			Name:     "dry run",
			Broker:   "tcp://localhost:1883",
			ClientID: "gigamock-test",
			Topic:    "devices/device-1/telemetry",
			QOS:      1,
			DryRun:   true,
			Message: scenarios.MQTTScenarioMessage{
				ContentType: "application/json",
				Body:        `{"ok":true}`,
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
	if response["topic"] != "devices/device-1/telemetry" {
		t.Fatalf("unexpected topic: %#v", response["topic"])
	}
	if response["dryRun"] != true {
		t.Fatalf("dryRun must be true: %#v", response["dryRun"])
	}
}

func TestMQTTProviderValidateRequiresScenarios(t *testing.T) {
	provider := MQTTProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func TestMQTTProviderValidateScenarioFields(t *testing.T) {
	provider := MQTTProvider{
		scenarios: scenarios.MQTTScenarios{
			{
				Topic: "missing-broker-client-id-and-body",
			},
		},
	}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for missing broker, clientID, and message body")
	}
}
