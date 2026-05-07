package scenarioType

import (
	"context"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/common"
)

func Factory(
	scenarioType string,
	w http.ResponseWriter,
	req *http.Request,
	lgr *logrus.Entry,
	ctx context.Context,
) (TypeProvider, error) {
	switch scenarioType {
	case common.HTTPScenarioType:
		return NewHTTPTypeProvider(w), nil
	case common.GraphQLScenarioType:
		return NewGraphQLTypeProvider(w, req), nil
	case common.KafkaScenarioType:
		return NewKafkaProvider(w, lgr, ctx), nil
	case common.GRPCScenarioType:
		return NewGRPCTypeProvider(w), nil
	case common.NATSScenarioType:
		return NewNATSProvider(w, lgr, ctx), nil
	case common.RabbitMQScenarioType:
		return NewRabbitMQProvider(w, lgr, ctx), nil
	case common.MQTTScenarioType:
		return NewMQTTProvider(w, lgr, ctx), nil
	case common.WebSocketScenarioType:
		return NewWebSocketProvider(w, req, lgr, ctx), nil
	}

	return nil, errors.New("scenario type provider is not reachable")
}
