# NATS Scenario Fields

NATS scenarios define the planned NATS publish/subscribe mock contract.
Gigamock currently indexes these files and displays them in the control UI.
Native NATS runtime support is planned.

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
| `host` | planned | NATS server URL. |
| `subject` | planned | NATS subject. |
| `headers` | planned | Message headers. |
| `message.body` | planned | Message body. |

Example file:

```text
examples/nats/order-created.yaml
```
