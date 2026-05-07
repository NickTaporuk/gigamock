package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NickTaporuk/gigamock/src/scenarios"
	"github.com/sirupsen/logrus"
)

func TestS3ProviderPutGetDeleteObject(t *testing.T) {
	resetS3TestStore()
	logger := logrus.NewEntry(logrus.New())
	scenario := scenarios.S3Scenario{ContentType: "application/json"}

	putReq := httptest.NewRequest(http.MethodPut, "/s3/test-bucket/customer.json", strings.NewReader(`{"name":"Ada"}`))
	putReq.Header.Set("Content-Type", "application/json")
	putRecorder := httptest.NewRecorder()
	putProvider := NewS3Provider(putRecorder, putReq, logger)
	putProvider.scenarios = scenarios.S3Scenarios{scenario}
	putProvider.Retrieve(0)
	if putRecorder.Code != http.StatusOK {
		t.Fatalf("expected PUT status %d, got %d with body %s", http.StatusOK, putRecorder.Code, putRecorder.Body.String())
	}
	if putRecorder.Header().Get("ETag") == "" {
		t.Fatal("PUT response must include ETag")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/s3/test-bucket/customer.json", nil)
	getRecorder := httptest.NewRecorder()
	getProvider := NewS3Provider(getRecorder, getReq, logger)
	getProvider.scenarios = scenarios.S3Scenarios{scenario}
	getProvider.Retrieve(0)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected GET status %d, got %d with body %s", http.StatusOK, getRecorder.Code, getRecorder.Body.String())
	}
	if getRecorder.Body.String() != `{"name":"Ada"}` {
		t.Fatalf("unexpected GET body: %s", getRecorder.Body.String())
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/s3/test-bucket/customer.json", nil)
	deleteRecorder := httptest.NewRecorder()
	deleteProvider := NewS3Provider(deleteRecorder, deleteReq, logger)
	deleteProvider.scenarios = scenarios.S3Scenarios{scenario}
	deleteProvider.Retrieve(0)
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected DELETE status %d, got %d", http.StatusNoContent, deleteRecorder.Code)
	}
}

func TestS3ProviderListBucket(t *testing.T) {
	resetS3TestStore()
	s3RuntimeStore.put("docs", s3Object{
		Key:          "readme.txt",
		Body:         []byte("hello"),
		ContentType:  "text/plain",
		ETag:         s3ETag([]byte("hello")),
		Metadata:     map[string]string{},
		LastModified: testTime(),
	})

	req := httptest.NewRequest(http.MethodGet, "/s3/docs", nil)
	recorder := httptest.NewRecorder()
	provider := NewS3Provider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.S3Scenarios{{Bucket: "docs"}}
	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d with body %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "<Key>readme.txt</Key>") {
		t.Fatalf("bucket list must include object key: %s", recorder.Body.String())
	}
}

func TestS3ProviderDryRun(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/s3/dry-run/example.txt", nil)
	recorder := httptest.NewRecorder()
	provider := NewS3Provider(recorder, req, logrus.NewEntry(logrus.New()))
	provider.scenarios = scenarios.S3Scenarios{{DryRun: true}}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected dry-run status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), `"dryRun":true`) {
		t.Fatalf("dry-run response is unexpected: %s", recorder.Body.String())
	}
}

func TestS3ProviderValidateRequiresScenarios(t *testing.T) {
	provider := S3Provider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}

func resetS3TestStore() {
	s3RuntimeStore.mu.Lock()
	defer s3RuntimeStore.mu.Unlock()
	s3RuntimeStore.Buckets = map[string]map[string]s3Object{}
}

func testTime() time.Time {
	return time.Date(2026, time.May, 7, 0, 0, 0, 0, time.UTC)
}
