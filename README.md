# gigamock
Gigamock is a generic utility to be used as mocked server to provide a simplified and consistent API over various network data sources such as http REST API or grpc or graphql services or as a mocking messaging systems like a kafka via mocking response or send message to some message system.
## Conception
### Parse YML files
    example.yaml
    
    path: "/test/test/:test"
    scenarious:
        - type: http
          name: "default scenareous"
          request:
            method: GET
          response:
            body: |
              {
                "settings": {
                  "patching": {
                    "deployment_rules": "do_not_wait_for_approval",
                   }
                }
              }
            
###
### Scenarios
### Feature
grpc api mock
parse swagger api to mock scenarios

