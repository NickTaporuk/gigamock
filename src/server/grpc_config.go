package server

// GRPCServerConfig contains production hardening settings for the dynamic gRPC runtime.
type GRPCServerConfig struct {
	Addr                 string
	TLSCertFile          string
	TLSKeyFile           string
	TLSClientCAFile      string
	MaxStreamMessages    int
	StreamTimeoutSeconds int
}
