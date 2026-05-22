package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

func TestHTTPTypeProvider_Validate(t *testing.T) {
	t.Run("success validation", func(t *testing.T) {
		provider := HTTPTypeProvider{
			w: nil,
			scenarios: scenarios.HTTPScenarios{{
				Request: scenarios.HTTPScenarioRequest{
					Headers:               nil,
					QueryStringParameters: nil,
					Cookies:               nil,
				},
				Response: scenarios.HTTPScenarioResponse{
					Body:       "test",
					StatusCode: 200,
					Headers:    nil,
					Cookies:    nil,
				},
			}},
		}

		err := provider.Validate()
		assert.NoError(t, err)
	})

	t.Run("failure validation", func(t *testing.T) {
		provider := HTTPTypeProvider{
			w: nil,
			scenarios: scenarios.HTTPScenarios{
				{
					Request: scenarios.HTTPScenarioRequest{
						Headers:               nil,
						QueryStringParameters: nil,
						Cookies:               nil,
					},
					Response: scenarios.HTTPScenarioResponse{
						Body:       "test",
						StatusCode: 0,
						Headers:    nil,
						Cookies:    nil,
					},
				},
				{
					Request: scenarios.HTTPScenarioRequest{
						Headers:               nil,
						QueryStringParameters: nil,
						Cookies:               nil,
					},
					Response: scenarios.HTTPScenarioResponse{},
				},
			},
		}

		err := provider.Validate()
		assert.Error(t, err)
	})

	t.Run("success validation for redirect without explicit status code", func(t *testing.T) {
		provider := HTTPTypeProvider{
			w: nil,
			scenarios: scenarios.HTTPScenarios{{
				Response: scenarios.HTTPScenarioResponse{
					Redirect: &scenarios.HTTPRedirect{
						URL: "https://payments.example/checkout",
					},
				},
			}},
		}

		err := provider.Validate()
		assert.NoError(t, err)
	})

	t.Run("failure validation for redirect with non-3xx status code", func(t *testing.T) {
		provider := HTTPTypeProvider{
			w: nil,
			scenarios: scenarios.HTTPScenarios{{
				Response: scenarios.HTTPScenarioResponse{
					StatusCode: http.StatusOK,
					Redirect: &scenarios.HTTPRedirect{
						URL: "https://payments.example/checkout",
					},
				},
			}},
		}

		err := provider.Validate()
		assert.Error(t, err)
	})
}

func TestHTTPTypeProvider_Retrieve(t *testing.T) {
	t.Run("returns configured redirect with query data", func(t *testing.T) {
		w := httptest.NewRecorder()
		provider := HTTPTypeProvider{
			w: w,
			scenarios: scenarios.HTTPScenarios{{
				Response: scenarios.HTTPScenarioResponse{
					StatusCode: http.StatusSeeOther,
					Headers: map[string]string{
						"X-Payment-Mock": "enabled",
					},
					Cookies: map[string]string{
						"payment_session": "session-1",
					},
					Body: "ignored for redirects",
					Redirect: &scenarios.HTTPRedirect{
						URL: "https://payments.example/checkout?locale=en",
						QueryStringParameters: map[string]string{
							"amount":    "4900",
							"currency":  "USD",
							"sessionId": "session-1",
						},
					},
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		assert.Equal(t, "enabled", w.Header().Get("X-Payment-Mock"))
		assert.Equal(t, "https://payments.example/checkout?amount=4900&currency=USD&locale=en&sessionId=session-1", w.Header().Get("Location"))
		assert.Contains(t, w.Header().Values("Set-Cookie"), "payment_session=session-1")
		assert.Empty(t, w.Body.String())
	})

	t.Run("defaults redirect status code to found", func(t *testing.T) {
		w := httptest.NewRecorder()
		provider := HTTPTypeProvider{
			w: w,
			scenarios: scenarios.HTTPScenarios{{
				Response: scenarios.HTTPScenarioResponse{
					Redirect: &scenarios.HTTPRedirect{
						URL: "https://payments.example/checkout",
					},
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "https://payments.example/checkout", w.Header().Get("Location"))
		assert.Empty(t, w.Body.String())
	})
}
