# NATS Scenario Fields

NATS scenarios publish configured messages to a NATS subject. For local route/UI
testing without a NATS broker, use `dryRun: true`; the mock will validate the
scenario, skip network calls, and return a successful JSON response.

Example:

```yaml
path: "/internal/nats/orders/:orderID"
method: POST
type: nats
description: "NATS mock scenario for publishing an order-created event"
scenarios:
  - name: "publish order created"
    host: "nats://localhost:4222"
    subject: "orders.created"
    dryRun: false
    headers:
      X-Request-Id: "831429af-1e40-4b44-8be3-06fd252578b0"
    message:
      body: |
        {
          "orderId": "order-1",
          "customerId": "customer-1",
          "status": "CREATED"
        }
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route used to index and switch the scenario. |
| `method` | yes | HTTP method used to index and switch the scenario. |
| `type` | yes | Must be `nats`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of NATS scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `host` | yes | NATS server URL, for example `nats://localhost:4222`. |
| `subject` | yes | NATS subject. |
| `dryRun` | no | When `true`, skips NATS broker calls and returns a successful response. |
| `headers` | no | Message headers. |
| `message.body` | yes | Message body. |

Runtime responses:

Successful publish response:

```json
{
  "subject": "orders.created",
  "published": true,
  "dryRun": false
}
```

Dry-run publish response:

```json
{
  "subject": "orders.created.dry-run",
  "published": true,
  "dryRun": true
}
```

Runtime metrics:

```bash
curl http://localhost:7777/internal/v1/nats/metrics
```

Example files:

```text
examples/nats/dry-run-order-created.yaml
examples/nats/order-created.yaml
```
