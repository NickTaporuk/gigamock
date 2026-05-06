package server

import "sync"

type grpcMetrics struct {
	mu      sync.RWMutex
	methods map[string]*grpcMethodMetrics
}

type grpcMethodMetrics struct {
	Calls  int64 `json:"calls"`
	Errors int64 `json:"errors"`
}

func newGRPCMetrics() *grpcMetrics {
	return &grpcMetrics{methods: map[string]*grpcMethodMetrics{}}
}

func (m *grpcMetrics) record(method string, failed bool) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	stats, ok := m.methods[method]
	if !ok {
		stats = &grpcMethodMetrics{}
		m.methods[method] = stats
	}
	stats.Calls++
	if failed {
		stats.Errors++
	}
}

func (m *grpcMetrics) snapshot() map[string]grpcMethodMetrics {
	if m == nil {
		return map[string]grpcMethodMetrics{}
	}
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[string]grpcMethodMetrics, len(m.methods))
	for method, stats := range m.methods {
		out[method] = *stats
	}
	return out
}
