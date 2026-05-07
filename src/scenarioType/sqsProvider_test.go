package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestSQSProviderSendReceivePurge(t *testing.T) {
	resetSQSTestStore()
	logger := logrus.NewEntry(logrus.New())
	scenario := scenarios.SQSScenario{}

	sendReq := httptest.NewRequest(http.MethodPost, "/aws/sqs/orders", strings.NewReader(`{"orderId":"order-1"}`))
	sendRecorder := httptest.NewRecorder()
	sendProvider := NewSQSProvider(sendRecorder, sendReq, logger)
	sendProvider.scenarios = scenarios.SQSScenarios{scenario}
	sendProvider.Retrieve(0)
	if sendRecorder.Code != http.StatusOK {
		t.Fatalf("expected send status %d, got %d with body %s", http.StatusOK, sendRecorder.Code, sendRecorder.Body.String())
	}
	if !strings.Contains(sendRecorder.Body.String(), `"messageId"`) {
		t.Fatalf("send response must include messageId: %s", sendRecorder.Body.String())
	}

	receiveReq := httptest.NewRequest(http.MethodGet, "/aws/sqs/orders", nil)
	receiveRecorder := httptest.NewRecorder()
	receiveProvider := NewSQSProvider(receiveRecorder, receiveReq, logger)
	receiveProvider.scenarios = scenarios.SQSScenarios{scenario}
	receiveProvider.Retrieve(0)
	if receiveRecorder.Code != http.StatusOK {
		t.Fatalf("expected receive status %d, got %d", http.StatusOK, receiveRecorder.Code)
	}
	if !strings.Contains(receiveRecorder.Body.String(), "order-1") {
		t.Fatalf("receive response must include message body: %s", receiveRecorder.Body.String())
	}

	purgeReq := httptest.NewRequest(http.MethodDelete, "/aws/sqs/orders", nil)
	purgeRecorder := httptest.NewRecorder()
	purgeProvider := NewSQSProvider(purgeRecorder, purgeReq, logger)
	purgeProvider.scenarios = scenarios.SQSScenarios{scenario}
	purgeProvider.Retrieve(0)
	if purgeRecorder.Code != http.StatusOK {
		t.Fatalf("expected purge status %d, got %d", http.StatusOK, purgeRecorder.Code)
	}
}

func TestSQSProviderDryRun(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/aws/sqs-dry-run/orders", nil)
	recorder := httptest.NewRecorder()
	provider := NewSQSProvider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.SQSScenarios{{DryRun: true}}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected dry-run status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"dryRun":true`) {
		t.Fatalf("dry-run response is unexpected: %s", recorder.Body.String())
	}
}

func TestSQSProviderValidateRequiresScenarios(t *testing.T) {
	provider := SQSProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func resetSQSTestStore() {
	sqsRuntimeStore.mu.Lock()
	defer sqsRuntimeStore.mu.Unlock()
	sqsRuntimeStore.Queues = map[string][]sqsMessage{}
}
