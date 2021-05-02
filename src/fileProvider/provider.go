package fileProvider

// FileProvider is an interface which will be extended to accomplish
// case with different types of file providers like a yml ot json
// that provides a possibility to use diff file parsing
type FileProvider interface {
	Init() error
	Parse(filePath string) (*GigaMockScenario, error)
}

type GigaMockScenario struct {
	Path   string `yaml:"path"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	Method string `yaml:"method"`

	Scenarios []Scenario `yaml:"scenarios"`
}

// Scenario
type Scenario struct {
	Request  ScenarioRequest  `yaml:"request"`
	Response ScenarioResponse `yaml:"response"`
}

// ScenarioRequest
type ScenarioRequest struct {
	Method string `yaml:"method"`
}

// ScenarioResponse
type ScenarioResponse struct {
	Body string `yaml:"body"`
}
