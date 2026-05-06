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
	case common.GraphQLScenarioType:
		return NewGraphQLTypeProvider(w, req), nil
	case common.KafkaScenarioType:
		return NewKafkaProvider(w, lgr), nil
	case common.GRPCScenarioType:
		return NewGRPCTypeProvider(w), nil
	case common.NATSScenarioType:
		return NewMessageBrokerTypeProvider(w, common.NATSScenarioType), nil
	case common.RabbitMQScenarioType:
		return NewMessageBrokerTypeProvider(w, common.RabbitMQScenarioType), nil
	case common.MQTTScenarioType:
		return NewMessageBrokerTypeProvider(w, common.MQTTScenarioType), nil
	}

	return nil, errors.New("scenario type provider is not reachable")
}
