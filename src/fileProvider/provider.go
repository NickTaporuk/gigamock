package fileProvider

// FileProvider is an interface which will be extended to accomplish
// case with different types of file providers like a yml ot json
// that provides a possibility to use diff file parsing
type FileProvider interface {
	Init() error
	Parse(filePath string) (GigaMockScenario, error)
}

type GigaMockScenario struct {
	Path      string
	Scenarios []Scenario
}

//     - type: http
//      name: "default scenareous"
//      request:
//        method: GET
//      response:
//        body: |
//          {
//            "settings": {
//              "patching": {
//                "deployment_rules": "do_not_wait_for_approval",
//               }
//            }
//          }
type Scenario struct {
	Type    string
	Name    string
	Request struct {
		Method string
	}
	Response struct {
		Body string
	}
}
