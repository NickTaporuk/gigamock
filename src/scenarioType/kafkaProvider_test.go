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

func TestKafkaProviderRetrieveDryRun(t *testing.T) {
	recorder := httptest.NewRecorder()
	provider := NewKafkaProvider(recorder, logrus.NewEntry(logrus.New()), context.Background())
	provider.scenarios = scenarios.KafkaScenarios{
		{
			Name:   "dry run",
			Host:   "localhost",
			Port:   "9092",
			Topic:  "test-dry-run",
			DryRun: true,
			Producer: &scenarios.KafkaScenarioProducer{
				Message: scenarios.KafkaScenarioProducerMessage{
					Key:   "key",
					Value: `{"ok":true}`,
				},
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
	if response["topic"] != "test-dry-run" {
		t.Fatalf("unexpected topic: %#v", response["topic"])
	}
	if response["dryRun"] != true {
		t.Fatalf("dryRun must be true: %#v", response["dryRun"])
	}
}

func TestKafkaProviderValidateRequiresScenarios(t *testing.T) {
	provider := KafkaProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func TestKafkaProviderValidateScenarioFields(t *testing.T) {
	provider := KafkaProvider{
		scenarios: scenarios.KafkaScenarios{
			{
				Topic: "missing-host-and-port",
				Producer: &scenarios.KafkaScenarioProducer{
					Message: scenarios.KafkaScenarioProducerMessage{Value: "payload"},
				},
			},
		},
	}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for missing host and port")
	}
}
