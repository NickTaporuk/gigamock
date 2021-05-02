package fileProvider

import "github.com/NickTaporuk/gigamock/src/scenarios"

// FileProvider is an interface which will be extended to accomplish
// case with different types of file providers like a yml ot json
// that provides a possibility to use diff file parsing
type FileProvider interface {
	Init() error
	Parse(filePath string) (*scenarios.GigaMockScenario, error)
}
