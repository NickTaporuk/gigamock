package scenarios

// GigaMockHTTPScenario
type GigaMockHTTPScenario struct {
	Scenarios []HTTPScenario `yaml:"scenarios"`
}

// Scenario
type HTTPScenario struct {
	Request  HTTPScenarioRequest  `yaml:"request"`
	Response HTTPScenarioResponse `yaml:"response"`
	Delay    uint                 `yaml:"delay,omitempty"`
	WebHook  WebHook              `yaml:"webhook,omitempty"`
}

// HTTPScenarios
type HTTPScenarios []HTTPScenario

// HTTPScenarioRequest
type HTTPScenarioRequest struct {
	Headers               map[string]string `yaml:"headers,omitempty"`
	QueryStringParameters map[string]string `yaml:"queryStringParameters,omitempty"`
	Cookies               map[string]string `yaml:"cookies,omitempty"`
}

// HTTPScenarioResponse
type HTTPScenarioResponse struct {
	Body       string            `yaml:"body,omitempty"`
	StatusCode int               `yaml:"statusCode"`
	Headers    map[string]string `yaml:"headers,omitempty"`
	Cookies    map[string]string `yaml:"cookies,omitempty"`
}
