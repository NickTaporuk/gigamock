package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestServiceBusProviderSendReceivePurge(t *testing.T) {
	resetServiceBusTestStore()
	logger := logrus.NewEntry(logrus.New())
	scenario := scenarios.ServiceBusScenario{}

	sendReq := httptest.NewRequest(http.MethodPost, "/azure/servicebus/orders/send", strings.NewReader(`{"orderId":"order-1"}`))
	sendRecorder := httptest.NewRecorder()
	sendProvider := NewServiceBusProvider(sendRecorder, sendReq, logger)
	sendProvider.scenarios = scenarios.ServiceBusScenarios{scenario}
	sendProvider.Retrieve(0)
	if sendRecorder.Code != http.StatusOK {
		t.Fatalf("expected send status %d, got %d with body %s", http.StatusOK, sendRecorder.Code, sendRecorder.Body.String())
	}
	if !strings.Contains(sendRecorder.Body.String(), `"messageId"`) {
		t.Fatalf("send response must include messageId: %s", sendRecorder.Body.String())
	}

	receiveReq := httptest.NewRequest(http.MethodPost, "/azure/servicebus/orders/receive", nil)
	receiveRecorder := httptest.NewRecorder()
	receiveProvider := NewServiceBusProvider(receiveRecorder, receiveReq, logger)
	receiveProvider.scenarios = scenarios.ServiceBusScenarios{scenario}
	receiveProvider.Retrieve(0)
	if receiveRecorder.Code != http.StatusOK {
		t.Fatalf("expected receive status %d, got %d", http.StatusOK, receiveRecorder.Code)
	}
	if !strings.Contains(receiveRecorder.Body.String(), "order-1") {
		t.Fatalf("receive response must include message body: %s", receiveRecorder.Body.String())
	}

	purgeReq := httptest.NewRequest(http.MethodDelete, "/azure/servicebus/orders/purge", nil)
	purgeRecorder := httptest.NewRecorder()
	purgeProvider := NewServiceBusProvider(purgeRecorder, purgeReq, logger)
	purgeProvider.scenarios = scenarios.ServiceBusScenarios{scenario}
	purgeProvider.Retrieve(0)
	if purgeRecorder.Code != http.StatusOK {
		t.Fatalf("expected purge status %d, got %d", http.StatusOK, purgeRecorder.Code)
	}
}

func TestServiceBusProviderDryRun(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/azure/servicebus-dry-run/orders/send", nil)
	recorder := httptest.NewRecorder()
	provider := NewServiceBusProvider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.ServiceBusScenarios{{DryRun: true}}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected dry-run status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"dryRun":true`) {
		t.Fatalf("dry-run response is unexpected: %s", recorder.Body.String())
	}
}

func TestServiceBusProviderValidateRequiresScenarios(t *testing.T) {
	provider := ServiceBusProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func resetServiceBusTestStore() {
	serviceBusRuntimeStore.mu.Lock()
	defer serviceBusRuntimeStore.mu.Unlock()
	serviceBusRuntimeStore.Entities = map[string][]serviceBusMessage{}
}
