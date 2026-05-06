package scenarioType

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

// GRPCTypeProvider is a placeholder HTTP-facing provider for indexed gRPC mock configs.
// A real gRPC runtime should serve these scenarios from a grpc.Server.
type GRPCTypeProvider struct {
	w         http.ResponseWriter
	scenarios []map[string]interface{}
}

// NewGRPCTypeProvider creates a provider for gRPC mock scenario files.
func NewGRPCTypeProvider(w http.ResponseWriter) *GRPCTypeProvider {
	return &GRPCTypeProvider{w: w}
}

// Unmarshal stores raw gRPC scenarios so files can be indexed and shown in the UI.
func (g *GRPCTypeProvider) Unmarshal(rawScenarios []map[string]interface{}) error {
	g.scenarios = rawScenarios
	return nil
}

// Validate validates decoded gRPC scenarios.
func (g GRPCTypeProvider) Validate() error {
	return validation.Validate(g.scenarios, validation.Required)
}

// Retrieve returns a clear response if a gRPC mock path is called through HTTP.
func (g *GRPCTypeProvider) Retrieve(scenarioNumber int) {
	if g.w == nil {
		return
	}

	g.w.Header().Set("Content-Type", "application/json")
	g.w.WriteHeader(http.StatusNotImplemented)
	g.w.Write([]byte(`{"error":"grpc mock runtime is not implemented yet; this scenario is indexed for the control UI"}`))
}
