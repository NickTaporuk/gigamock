package webhook

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/NickTaporuk/gigamock/src/common"
)

type WebHooker interface {
	Scenarios() []map[string]interface{}
	Type() string
	Method() string

	validation.Validatable
}

// WebHook
type WebHook struct {
	Scenarios []map[string]interface{} `yaml:"scenarios",json:"scenarios"`
	Type      string                   `yaml:"type",json:"type"`
	Method    string                   `yaml:"method",json:"method"`
	Path      string                   `yaml:"path",json:"path"`
}

func (w *WebHook) Validate() error {
	if w == nil {
		return nil
	}

	switch w.Type {
	case common.HTTPScenarioType:
		pntr := *w
		return validation.ValidateStruct(
			&pntr,
			validation.Field(&pntr.Type, common.ScenarioTypeValidator...),
			validation.Field(&pntr.Method, common.ScenarioMethodValidator...),
			validation.Field(&pntr.Path, common.URLPathValidator...),
			validation.Field(&pntr.Scenarios, common.BaseScenariosValidator...),
		)
	}

	return nil
}
