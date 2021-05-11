package scenarioType

import (
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
}
