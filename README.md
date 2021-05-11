# gigamock
Gigamock is a generic utility to be used as mocked server to provide a simplified and consistent API over various network data sources such as http REST API or grpc or graphql services via mocking response.
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
            
http rest api mock
###
### Scenarious
### Feature
grpc api mock
swagger api to mock

