package scenarioType

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/NickTaporuk/gigamock/src/common"
)

func Factory(
	scenarioType string,
	w http.ResponseWriter,
	req *http.Request,
	lgr *logrus.Entry,
) (TypeProvider, error) {
	switch scenarioType {
	case common.HTTPScenarioType:
		return NewHTTPTypeProvider(w), nil
	case common.KafkaScenarioType:
		return NewKafkaProvider(w, lgr), nil
	}

	return nil, errors.New("scenario type provider is not reachable")
}
