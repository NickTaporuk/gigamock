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
}
