@mqtt
Feature: MQTT mock scenario contracts
  MQTT scenario files should be indexed and visible in the control UI.

  Background:
    Given Gigamock is running with directory "./examples/mqtt"

  Scenario: List MQTT scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "device-telemetry.yaml"
    And the response body should contain "devices/device-1/telemetry"

  Scenario: HTTP fallback explains runtime status
    When I send a POST request to "/internal/mqtt/devices/device-1/telemetry"
    Then the response status should be 501
    And the response body should contain "mqtt mock runtime is not implemented yet"
