# HTTP Scenario Fields

HTTP scenarios return HTTP responses from YAML or JSON config files.

Example:

```yaml
path: "/control-ui/users"
method: GET
type: http
description: "HTTP example with multiple scenarios for testing the control UI"
scenarios:
  - name: "active users"
    request:
    response:
      statusCode: 200
      headers:
        Content-Type: "application/json"
      body: |
        {
          "users": []
        }
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route. Route params are supported, for example `/users/:userID`. |
| `method` | yes | HTTP method such as `GET`, `POST`, `PUT`, or `DELETE`. |
| `type` | yes | Must be `http`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of HTTP scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `request` | no | Request matching data. Currently mostly descriptive for HTTP scenarios. |
| `response` | yes | Response returned by the mock server. |
| `delay` | no | Planned delay value such as `100ms`. |

Request fields:

| Field | Required | Description |
| --- | --- | --- |
| `headers` | no | Expected request headers. |
| `queryStringParameters` | no | Expected query string parameters. |
| `cookies` | no | Expected cookies. |
| `body` | no | Expected request body. |

Response fields:

| Field | Required | Description |
| --- | --- | --- |
| `statusCode` | yes | HTTP status code. Optional when `redirect` is set; defaults to `302`. Redirect responses must use a `3xx` status. |
| `headers` | no | Response headers. |
| `cookies` | no | Response cookies. |
| `body` | no | Response body as a string. |
| `redirect` | no | Redirect target. When set, Gigamock writes the `Location` header and does not write `body`. |

Redirect fields:

| Field | Required | Description |
| --- | --- | --- |
| `url` | yes | Redirect target URL. Existing query string values are preserved. |
| `queryStringParameters` | no | Query string values appended to the redirect URL. Use this for payment/session data. |

Redirect example:

```yaml
path: "/payments"
method: POST
type: http
scenarios:
  - name: "checkout redirect"
    request:
    response:
      statusCode: 303
      redirect:
        url: "https://payments.example/checkout"
        queryStringParameters:
          sessionId: "session-1"
          amount: "4900"
          currency: "USD"
```

Example file:

```text
examples/rest/control-ui-users.yaml
```
