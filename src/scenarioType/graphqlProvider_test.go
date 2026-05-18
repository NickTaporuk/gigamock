package scenarioType

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

func TestGraphQLTypeProvider_Retrieve(t *testing.T) {
	t.Run("returns active scenario response", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":"GetHero"}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Request: scenarios.GraphQLScenarioRequest{
					OperationName: "GetHero",
				},
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: `{"data":{"hero":{"name":"Luke Skywalker"}}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"data":{"hero":{"name":"Luke Skywalker"}}}`, w.Body.String())
	})

	t.Run("falls back to matching scenario when active scenario does not match request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":"GetVillain"}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{
				{
					Request: scenarios.GraphQLScenarioRequest{
						OperationName: "GetHero",
					},
					Response: scenarios.GraphQLScenarioResponse{
						StatusCode: http.StatusOK,
						Body:       `{"data":{"hero":{"name":"Luke Skywalker"}}}`,
					},
				},
				{
					Request: scenarios.GraphQLScenarioRequest{
						OperationName: "GetVillain",
					},
					Response: scenarios.GraphQLScenarioResponse{
						StatusCode: http.StatusOK,
						Body:       `{"data":{"villain":{"name":"Darth Vader"}}}`,
					},
				},
			},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"villain":{"name":"Darth Vader"}}}`, w.Body.String())
	})

	t.Run("matches scenario by headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":"PrivateHero"}`))
		req.Header.Set("X-Tenant", "jedi")
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Request: scenarios.GraphQLScenarioRequest{
					Headers:       map[string]string{"X-Tenant": "jedi"},
					OperationName: "PrivateHero",
				},
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{"data":{"privateHero":{"name":"Obi-Wan Kenobi"}}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"privateHero":{"name":"Obi-Wan Kenobi"}}}`, w.Body.String())
	})

	t.Run("filters response data to requested fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{
			"query":"query { balance(id: 2) { id } }"
		}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{"data":{"balance":{"id":"balance-1","currency":"USD","available":4900}}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"balance":{"id":"balance-1"}}}`, w.Body.String())
	})

	t.Run("filters response data to requested fields for shorthand field query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{
			"query":"balance(id: 2) { id }"
		}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{"data":{"balance":{"id":"balance-1","currency":"USD","available":4900}}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"balance":{"id":"balance-1"}}}`, w.Body.String())
	})

	t.Run("filters response data to requested fields for operation-like shorthand query", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{
			"query":"query balance(id: 2) { id }"
		}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{"data":{"balance":{"id":"balance-1","currency":"USD","available":4900}}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"balance":{"id":"balance-1"}}}`, w.Body.String())
	})

	t.Run("filters nested response data and lists", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{
			"query":"query { accounts { id balance { available } } }"
		}`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body: `{"data":{"accounts":[
						{"id":"account-1","name":"Checking","balance":{"currency":"USD","available":4900}},
						{"id":"account-2","name":"Savings","balance":{"currency":"EUR","available":1200}}
					]}}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"data":{"accounts":[
			{"id":"account-1","balance":{"available":4900}},
			{"id":"account-2","balance":{"available":1200}}
		]}}`, w.Body.String())
	})

	t.Run("returns bad request for invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`{"operationName":`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{}`,
				},
			}},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"errors":[{"message":"invalid GraphQL JSON request body"}]}`, w.Body.String())
	})

	t.Run("supports batch GraphQL requests", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/graphql", strings.NewReader(`[
			{"operationName":"GetHero","query":"query { hero { name } }"},
			{"operationName":"GetVillain","query":"query { villain { name } }"}
		]`))
		provider := GraphQLTypeProvider{
			w:   w,
			req: req,
			scenarios: scenarios.GraphQLScenarios{
				{
					Request: scenarios.GraphQLScenarioRequest{
						OperationName: "GetHero",
					},
					Response: scenarios.GraphQLScenarioResponse{
						StatusCode: http.StatusOK,
						Body:       `{"data":{"hero":{"name":"Luke Skywalker"}}}`,
					},
				},
				{
					Request: scenarios.GraphQLScenarioRequest{
						OperationName: "GetVillain",
					},
					Response: scenarios.GraphQLScenarioResponse{
						StatusCode: http.StatusOK,
						Body:       `{"data":{"villain":{"name":"Darth Vader"}}}`,
					},
				},
			},
		}

		provider.Retrieve(0)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `[
			{"data":{"hero":{"name":"Luke Skywalker"}}},
			{"data":{"villain":{"name":"Darth Vader"}}}
		]`, w.Body.String())
	})

	t.Run("validates response body JSON", func(t *testing.T) {
		provider := GraphQLTypeProvider{
			scenarios: scenarios.GraphQLScenarios{{
				Response: scenarios.GraphQLScenarioResponse{
					StatusCode: http.StatusOK,
					Body:       `{"data":`,
				},
			}},
		}

		assert.Error(t, provider.Validate())
	})
}
