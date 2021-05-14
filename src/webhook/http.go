package webhook

// HTTPProvider
type HTTPProvider struct{}

// NewHTTPProvider
func NewHTTPProvider() *HTTPProvider {
	return &HTTPProvider{}
}

// Unmarshal
func (h *HTTPProvider) Unmarshal(i []map[string]interface{}) error {
	panic("implement me")
}

// Send
func (h *HTTPProvider) Send() error {
	panic("implement me")
}

// Validate
func (h HTTPProvider) Validate() error {
	panic("implement me")
}
