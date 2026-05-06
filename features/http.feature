@http
Feature: HTTP mock scenarios
  HTTP scenario files should be indexed and should return configured responses.

  Background:
    Given Gigamock is running with directory "./examples/rest"

  Scenario: Retrieve the default HTTP response
    When I send a GET request to "/control-ui/users"
    Then the response status should be 200
    And the response body should contain "users"

  Scenario: Switch the active HTTP scenario
    When I set scenario 1 for path "/control-ui/users" and method "GET"
    And I send a GET request to "/control-ui/users"
    Then the response status should be 200
    And the response body should contain "\"users\": []"
