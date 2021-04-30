# gigamock
this is a mocking HTTP requests package
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

