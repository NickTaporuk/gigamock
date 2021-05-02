package scenarios

// GigaMockScenario
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
