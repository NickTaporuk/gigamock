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
| `statusCode` | yes | HTTP status code. |
| `headers` | no | Response headers. |
| `cookies` | no | Response cookies. |
| `body` | no | Response body as a string. |

Example file:

```text
examples/rest/control-ui-users.yaml
```
