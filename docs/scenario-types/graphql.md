# GraphQL Scenario Fields

GraphQL scenarios are served over HTTP and support matching by GraphQL payload.
This is useful because many GraphQL operations share a single endpoint such as
`/graphql`.

Example:

```yaml
path: "/graphql"
method: POST
type: graphql
description: "GraphQL mock with operationName and variables matching"
scenarios:
  - name: "get hero by episode"
    request:
      operationName: "GetHero"
      query: |
        query GetHero($episode: String!) {
          hero(episode: $episode) {
            id
            name
          }
        }
      variables:
        episode: "NEWHOPE"
    response:
      statusCode: 200
      headers:
        Content-Type: "application/json"
      body: |
        {
          "data": {
            "hero": {
              "id": "1000",
              "name": "Luke Skywalker"
            }
          }
        }
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP endpoint for GraphQL requests. |
| `method` | yes | Usually `POST`. |
| `type` | yes | Must be `graphql`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of GraphQL scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `request` | no | GraphQL request matcher. |
| `response` | yes | HTTP response returned to the GraphQL client. |
| `delay` | no | Planned delay value. |

Request fields:

| Field | Required | Description |
| --- | --- | --- |
| `operationName` | no | GraphQL operation name to match. |
| `query` | no | GraphQL query. Whitespace is normalized before matching. |
| `variables` | no | JSON-like variables map to match. |
| `headers` | no | Expected request headers. Header matching is enforced when configured. |
| `body` | no | Raw body, mostly descriptive. |

Response fields:

| Field | Required | Description |
| --- | --- | --- |
| `statusCode` | yes | HTTP status code. GraphQL errors usually still use `200`. |
| `headers` | no | Response headers. |
| `body` | no | GraphQL JSON response body. |

Matching behavior:

- Gigamock first tries the currently active scenario.
- If the active scenario is `0` and its request does not match the incoming GraphQL payload,
  Gigamock searches the other scenarios for a matching `operationName`, `query`,
- `query`, `variables`, and configured `headers`.
- If a non-zero scenario is selected through the control API/UI, that scenario
  is returned directly. This makes manual error-state testing predictable.
- Batched GraphQL HTTP requests are supported when the body is a JSON array.
- Invalid JSON returns `400` with a GraphQL-style `errors` response.

Metrics:

```bash
curl http://localhost:7777/internal/v1/graphql/metrics
```

Metrics are grouped by GraphQL HTTP endpoint and include calls, errors,
matched requests, unmatched requests, batch calls, invalid JSON, and forced
active scenario selections.

Example file:

```text
examples/graphql/starwars-operations.yaml
```

Batch request example:

```bash
curl -X POST http://localhost:7777/graphql \
  -H "Content-Type: application/json" \
  -d '[
    {
      "operationName": "GetHero",
      "query": "query GetHero($episode: String!) { hero(episode: $episode) { id name } }",
      "variables": {
        "episode": "NEWHOPE"
      }
    },
    {
      "operationName": "GetVillain",
      "query": "query GetVillain { villain { id name } }"
    }
  ]'
```
