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

func TestRabbitMQProviderRetrieveDryRun(t *testing.T) {
	recorder := httptest.NewRecorder()
	provider := NewRabbitMQProvider(recorder, logrus.NewEntry(logrus.New()), context.Background())
	provider.scenarios = scenarios.RabbitMQScenarios{
		{
			Name:       "dry run",
			URL:        "amqp://guest:guest@localhost:5672/",
			Exchange:   "payments",
			RoutingKey: "payments.captured",
			DryRun:     true,
			Message: scenarios.RabbitMQScenarioMessage{
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
	if response["exchange"] != "payments" {
		t.Fatalf("unexpected exchange: %#v", response["exchange"])
	}
	if response["routingKey"] != "payments.captured" {
		t.Fatalf("unexpected routing key: %#v", response["routingKey"])
	}
	if response["dryRun"] != true {
		t.Fatalf("dryRun must be true: %#v", response["dryRun"])
	}
}

func TestRabbitMQProviderValidateRequiresScenarios(t *testing.T) {
	provider := RabbitMQProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func TestRabbitMQProviderValidateScenarioFields(t *testing.T) {
	provider := RabbitMQProvider{
		scenarios: scenarios.RabbitMQScenarios{
			{
				Exchange: "missing-url-routing-key-and-body",
			},
		},
	}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for missing url, routing key, and message body")
	}
}
