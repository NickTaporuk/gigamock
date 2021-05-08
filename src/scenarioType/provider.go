package scenarioType

// TypeProvider
type TypeProvider interface {
	// Unmarshal
	Unmarshal([]map[string]interface{}) error
	Retrieve(scenarioNumber int)
}

type TypeUnmarshaller interface {
	Data() TypeUnmarshaller
}
