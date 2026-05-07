package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestPubSubProviderPublishPullPurge(t *testing.T) {
	resetPubSubTestStore()
	logger := logrus.NewEntry(logrus.New())
	publishScenario := scenarios.PubSubScenario{}

	publishReq := httptest.NewRequest(http.MethodPost, "/gcp/pubsub/order-events/publish", strings.NewReader(`{"orderId":"order-1"}`))
	publishRecorder := httptest.NewRecorder()
	publishProvider := NewPubSubProvider(publishRecorder, publishReq, logger)
	publishProvider.scenarios = scenarios.PubSubScenarios{publishScenario}
	publishProvider.Retrieve(0)
	if publishRecorder.Code != http.StatusOK {
		t.Fatalf("expected publish status %d, got %d with body %s", http.StatusOK, publishRecorder.Code, publishRecorder.Body.String())
	}
	if !strings.Contains(publishRecorder.Body.String(), `"messageIds"`) {
		t.Fatalf("publish response must include messageIds: %s", publishRecorder.Body.String())
	}

	pullReq := httptest.NewRequest(http.MethodPost, "/gcp/pubsub/orders-sub/pull", nil)
	pullRecorder := httptest.NewRecorder()
	pullProvider := NewPubSubProvider(pullRecorder, pullReq, logger)
	pullProvider.scenarios = scenarios.PubSubScenarios{{Topic: "order-events", Subscription: "orders-sub"}}
	pullProvider.Retrieve(0)
	if pullRecorder.Code != http.StatusOK {
		t.Fatalf("expected pull status %d, got %d", http.StatusOK, pullRecorder.Code)
	}
	if !strings.Contains(pullRecorder.Body.String(), "order-1") {
		t.Fatalf("pull response must include message body: %s", pullRecorder.Body.String())
	}

	purgeReq := httptest.NewRequest(http.MethodDelete, "/gcp/pubsub/orders-sub/purge", nil)
	purgeRecorder := httptest.NewRecorder()
	purgeProvider := NewPubSubProvider(purgeRecorder, purgeReq, logger)
	purgeProvider.scenarios = scenarios.PubSubScenarios{{Topic: "order-events", Subscription: "orders-sub"}}
	purgeProvider.Retrieve(0)
	if purgeRecorder.Code != http.StatusOK {
		t.Fatalf("expected purge status %d, got %d", http.StatusOK, purgeRecorder.Code)
	}
}

func TestPubSubProviderDryRun(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/gcp/pubsub-dry-run/order-events/publish", nil)
	recorder := httptest.NewRecorder()
	provider := NewPubSubProvider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.PubSubScenarios{{DryRun: true}}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected dry-run status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"dryRun":true`) {
		t.Fatalf("dry-run response is unexpected: %s", recorder.Body.String())
	}
}

func TestPubSubProviderValidateRequiresScenarios(t *testing.T) {
	provider := PubSubProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func resetPubSubTestStore() {
	pubSubRuntimeStore.mu.Lock()
	defer pubSubRuntimeStore.mu.Unlock()
	pubSubRuntimeStore.Topics = map[string][]pubSubMessage{}
}
