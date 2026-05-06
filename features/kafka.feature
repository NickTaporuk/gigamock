@kafka
Feature: Kafka mock scenarios
  Kafka scenarios should be indexed and should trigger configured producer behavior.

  Background:
    Given Gigamock is running with directory "./examples/kafka"
    And Kafka is available at "0.0.0.0:9092"

  Scenario: Trigger a Kafka producer scenario
    When I send a GET request to "/internal/queue/message-1"
    Then the response status should be 200

  Scenario: List Kafka scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "test-topic.yaml"
    And the response body should contain "kafka"
