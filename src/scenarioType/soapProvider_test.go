package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

func TestSOAPProviderMatchesByActionAndBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/soap/customers", strings.NewReader(`<customerId>customer-1</customerId>`))
	req.Header.Set("SOAPAction", `"GetCustomer"`)
	recorder := httptest.NewRecorder()
	provider := NewSOAPProvider(recorder, req)
	provider.scenarios = scenarios.SOAPScenarios{
		{
			Request: scenarios.SOAPScenarioRequest{SOAPAction: "GetCustomer", BodyContains: "<customerId>customer-1</customerId>"},
			Response: scenarios.SOAPScenarioResponse{
				StatusCode: http.StatusOK,
				Body:       `<Envelope><Body><name>Ada</name></Body></Envelope>`,
			},
		},
	}

	provider.Retrieve(0)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), "Ada") {
		t.Fatalf("expected SOAP response body: %s", recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "text/xml; charset=utf-8" {
		t.Fatalf("expected SOAP content type, got %s", recorder.Header().Get("Content-Type"))
	}
}

func TestSOAPProviderFindsMatchingScenario(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/soap/customers", strings.NewReader(`<customerId>missing</customerId>`))
	req.Header.Set("SOAPAction", "GetCustomer")
	recorder := httptest.NewRecorder()
	provider := NewSOAPProvider(recorder, req)
	provider.scenarios = scenarios.SOAPScenarios{
		{
			Request:  scenarios.SOAPScenarioRequest{SOAPAction: "GetCustomer", BodyContains: "<customerId>customer-1</customerId>"},
			Response: scenarios.SOAPScenarioResponse{StatusCode: http.StatusOK, Body: `<ok/>`},
		},
		{
			Request:  scenarios.SOAPScenarioRequest{SOAPAction: "GetCustomer", BodyContains: "<customerId>missing</customerId>"},
			Response: scenarios.SOAPScenarioResponse{StatusCode: http.StatusInternalServerError, Body: `<fault>customer is not found</fault>`},
		},
	}

	provider.Retrieve(0)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "customer is not found") {
		t.Fatalf("expected matched SOAP fault body: %s", recorder.Body.String())
	}
}

func TestSOAPProviderValidateRequiresScenarios(t *testing.T) {
	provider := SOAPProvider{}

	if err := provider.Validate(); err == nil {
		t.Fatal("expected validation error for empty scenarios")
	}
}
