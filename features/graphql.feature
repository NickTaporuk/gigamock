@graphql
Feature: GraphQL mock scenarios
  GraphQL scenarios should be matched by operationName, query, and variables.

  Background:
    Given Gigamock is running with directory "./examples/graphql"

  Scenario: Match a GraphQL operation by operationName and variables
    When I send a POST request to "/graphql" with JSON body:
      """
      {
        "operationName": "GetHero",
        "query": "query GetHero($episode: String!) { hero(episode: $episode) { id name } }",
        "variables": {
          "episode": "NEWHOPE"
        }
      }
      """
    Then the response status should be 200
    And the response body should contain "Luke Skywalker"

  Scenario: Match another GraphQL operation on the same endpoint
    When I send a POST request to "/graphql" with JSON body:
      """
      {
        "operationName": "GetVillain",
        "query": "query GetVillain { villain { id name } }"
      }
      """
    Then the response status should be 200
    And the response body should contain "Darth Vader"
