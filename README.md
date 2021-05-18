# gigamock
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_shield)

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



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_large)