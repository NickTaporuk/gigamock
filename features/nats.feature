@nats
Feature: NATS mock scenario contracts
  NATS scenario files should be indexed and visible in the control UI.

  Background:
    Given Gigamock is running with directory "./examples/nats"

  Scenario: List NATS scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "order-created.yaml"
    And the response body should contain "orders.created"

  Scenario: HTTP fallback explains runtime status
    When I send a POST request to "/internal/nats/orders/order-1"
    Then the response status should be 501
    And the response body should contain "nats mock runtime is not implemented yet"
