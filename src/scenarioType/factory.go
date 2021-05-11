package scenarioType

import (
	"errors"
	"net/http"

	"github.com/NickTaporuk/gigamock/src/scenarios"
)

func Factory(scenarioType string, w http.ResponseWriter, req *http.Request) (TypeProvider, error) {
	switch scenarioType {
	case scenarios.HTTPScenarioType:
		return NewHTTPTypeProvider(w), nil
	}

	return nil, errors.New("scenario type provider is not reachable")
}
