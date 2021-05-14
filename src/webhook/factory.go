package webhook

import (
	"github.com/NickTaporuk/gigamock/src/common"
)

// Factory represents a pattern factory to select right provider by web hook type parameter
// which can be included to a description for a scenario file
// by example file examples/rest/create-user.yaml include the section webhook with all fields
func Factory(webHookType string) (WebHookTypeProvider, error) {
	switch webHookType {
	case common.HTTPScenarioType:
		return NewHTTPProvider(), nil
	default:
		return nil, nil
	}
}
