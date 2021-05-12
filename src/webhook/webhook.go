package webhook

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

type WebHooker interface {
	validation.Validatable
	Run() error
}

// WebHook
type WebHook struct {
	Scenarios []map[string]interface{} `yaml:"scenarios",json:"scenarios"`
	Type      string                   `yaml:"type",json:"type"`
	Method    string                   `yaml:"method",json:"method"`
}

func (WebHook) Run() error {
	return nil
}

func (w *WebHook) Validate() error {
	if w == nil {
		return nil
	}

	switch w.Type {
	case common.HTTPScenarioType:
		return validation.ValidateStruct(
			w,
			validation.Field(&w.Type, common.ScenarioTypeValidator...),
			validation.Field(&w.Method, common.ScenarioMethodValidator...),
			validation.Field(&w.Scenarios, common.BaseScenariosValidator...),
		)
	}

	return nil
}
