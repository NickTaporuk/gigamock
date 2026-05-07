package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestSNSProviderPublishAndList(t *testing.T) {
	resetSNSTestStore()
	logger := logrus.NewEntry(logrus.New())
	scenario := scenarios.SNSScenario{Subject: "order created"}

	publishReq := httptest.NewRequest(http.MethodPost, "/aws/sns/order-events", strings.NewReader(`{"orderId":"order-1"}`))
	publishRecorder := httptest.NewRecorder()
	publishProvider := NewSNSProvider(publishRecorder, publishReq, logger)
	publishProvider.scenarios = scenarios.SNSScenarios{scenario}
	publishProvider.Retrieve(0)
	if publishRecorder.Code != http.StatusOK {
		t.Fatalf("expected publish status %d, got %d with body %s", http.StatusOK, publishRecorder.Code, publishRecorder.Body.String())
	}
	if !strings.Contains(publishRecorder.Body.String(), `"messageId"`) {
		t.Fatalf("publish response must include messageId: %s", publishRecorder.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/aws/sns/order-events", nil)
	listRecorder := httptest.NewRecorder()
	listProvider := NewSNSProvider(listRecorder, listReq, logger)
	listProvider.scenarios = scenarios.SNSScenarios{scenario}
	listProvider.Retrieve(0)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listRecorder.Code)
	}
	if !strings.Contains(listRecorder.Body.String(), "order-1") {
		t.Fatalf("list response must include published message body: %s", listRecorder.Body.String())
	}
}

func TestSNSProviderDryRun(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/aws/sns-dry-run/order-events", nil)
	recorder := httptest.NewRecorder()
	provider := NewSNSProvider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.SNSScenarios{{DryRun: true}}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected dry-run status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"dryRun":true`) {
		t.Fatalf("dry-run response is unexpected: %s", recorder.Body.String())
	}
}

func TestSNSProviderValidateRequiresScenarios(t *testing.T) {
	provider := SNSProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func resetSNSTestStore() {
	snsRuntimeStore.mu.Lock()
	defer snsRuntimeStore.mu.Unlock()
	snsRuntimeStore.Topics = map[string][]snsMessage{}
}
