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
	Delay    uint             `yaml:"delay,omitempty"`
	WebHook  WebHook
}

// ScenarioRequest
type ScenarioRequest struct {
	Headers               map[string]string `yaml:"headers,omitempty"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
	Cookies               map[string]string `yaml:"cookies,omitempty"`
}

// ScenarioResponse
type ScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
	Cookies    map[string]string `yaml:"cookies,omitempty"`
}

// WebHook
type WebHook struct {
	URL                   string            `yaml:"url"`
	Method                string            `yaml:"method"`
	Headers               map[string]string `yaml:"headers,omitempty"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
	Cookies               map[string]string `yaml:"cookies,omitempty"`
	Type                  string            `yaml:"type"` // can be http or graphql or grpc and so one
}
