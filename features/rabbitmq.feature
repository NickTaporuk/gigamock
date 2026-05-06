@rabbitmq
Feature: RabbitMQ mock scenario contracts
  RabbitMQ scenario files should be indexed and visible in the control UI.

  Background:
    Given Gigamock is running with directory "./examples/rabbitmq"

  Scenario: List RabbitMQ scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "payment-events.yaml"
    And the response body should contain "payments.captured"

  Scenario: HTTP fallback explains runtime status
    When I send a POST request to "/internal/rabbitmq/payments/payment-1"
    Then the response status should be 501
    And the response body should contain "rabbitmq mock runtime is not implemented yet"
