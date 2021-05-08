package scenarios

// WebHook
type WebHook struct {
	Scenarios map[string]interface{} `yaml:"scenarios"`
	Type      string                 `yaml:"type"` // can be http or graphql or grpc and so one
}
