package webhookType

import (
	"github.com/sirupsen/logrus"

	"github.com/NickTaporuk/gigamock/src/common"
	"github.com/NickTaporuk/gigamock/src/webhook"
)

// Factory represents a pattern factory to select right provider by web hook type parameter
// which can be included to a description for a scenario file
// by example file examples/rest/create-user.yaml include the section webhook with all fields
func Factory(
	webhook *webhook.WebHook,
	lgr *logrus.Entry,
	scenarioNumber int,
	) (TypeProvider, error) {
	switch webhook.Type {
	case common.HTTPScenarioType:
		return NewHTTPProvider(lgr, webhook, scenarioNumber), nil
	default:
		return nil, nil
	}
}
