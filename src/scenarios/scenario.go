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
	Delay    uint             `yaml:"delay"`
	WebHook  WebHook
	Control  Control
}

// ScenarioRequest
type ScenarioRequest struct {
	Headers               map[string]string `yaml:"headers"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters"`
	Cookies               map[string]string `yaml:"cookies"`
}

// ScenarioResponse
type ScenarioResponse struct {
	Body       string            `yaml:"body"`
	StatusCode uint              `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers"`
	Cookies    map[string]string `yaml:"cookies"`
}

// WebHook
type WebHook struct {
	URL                   string            `yaml:"url"`
	Method                string            `yaml:"method"`
	Headers               map[string]string `yaml:"headers"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters"`
	Cookies               map[string]string `yaml:"cookies"`
	Type                  string            `yaml:"type"` // can be http or graphql or grpc
}

type Control struct {
	RequiredState []string `yaml:"requiredState"`
	NewState      string   `yaml:"newState"`
}
