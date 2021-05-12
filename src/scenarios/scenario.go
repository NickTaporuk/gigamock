package scenarios

import (
	"github.com/NickTaporuk/gigamock/src/webhook"
)

//BaseGigaMockScenario describes a base fields to parse file fields
type BaseGigaMockScenario struct {
	Path   string `yaml:"path",json:"path",xml:"path"`
	Type   string `yaml:"type",json:"type",xml:"type"`
	Name   string `yaml:"name",json:"name",xml:"name"`
	Method string `yaml:"method",json:"method",xml:"method"`

	Scenarios []map[string]interface{} `yaml:"scenarios",json:"scenarios",xml:"scenarios"`
	WebHook   *webhook.WebHook          `xml:"webhook",json:"webhook",yaml:"webhook"`
}
