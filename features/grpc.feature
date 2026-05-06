@grpc
Feature: gRPC mock scenario contracts
  gRPC scenario files should be indexed and visible in the control UI.

  Background:
    Given Gigamock is running with directory "./examples/grpc"

  Scenario: List unary gRPC scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "/customers.CustomersService/GetCustomer"
    And the response body should contain "customer-service-unary.yaml"

  Scenario: List bidirectional gRPC scenario metadata
    When I request "/internal/v1/scenarios"
    Then the response status should be 200
    And the response body should contain "/chat.ChatService/Chat"
    And the response body should contain "chat-service-bidi-stream.yaml"
