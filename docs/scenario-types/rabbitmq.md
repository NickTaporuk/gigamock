# RabbitMQ Scenario Fields

RabbitMQ scenarios publish configured messages to an exchange with a routing
key. For local route/UI testing without a RabbitMQ broker, use `dryRun: true`;
the mock will validate the scenario, skip network calls, and return a successful
JSON response.

Example:

```yaml
path: "/internal/rabbitmq/payments/:paymentID"
method: POST
type: rabbitmq
description: "RabbitMQ mock scenario for publishing payment events"
scenarios:
  - name: "publish payment captured"
    url: "amqp://guest:guest@localhost:5672/"
    exchange: "payments"
    routingKey: "payments.captured"
    dryRun: false
    headers:
      X-Request-Id: "831429af-1e40-4b44-8be3-06fd252578b0"
    message:
      contentType: "application/json"
      body: |
        {
          "paymentId": "payment-1",
          "orderId": "order-1",
          "status": "CAPTURED"
        }
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route used to index and switch the scenario. |
| `method` | yes | HTTP method used to index and switch the scenario. |
| `type` | yes | Must be `rabbitmq`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of RabbitMQ scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `url` | yes | RabbitMQ AMQP connection URL. |
| `exchange` | yes | Exchange name. |
| `routingKey` | yes | Routing key. |
| `dryRun` | no | When `true`, skips RabbitMQ broker calls and returns a successful response. |
| `headers` | no | Message headers. |
| `message.contentType` | no | Message content type. |
| `message.body` | yes | Message body. |

Runtime responses:

Successful publish response:

```json
{
  "exchange": "payments",
  "routingKey": "payments.captured",
  "published": true,
  "dryRun": false
}
```

Dry-run publish response:

```json
{
  "exchange": "payments",
  "routingKey": "payments.captured.dry-run",
  "published": true,
  "dryRun": true
}
```

Runtime metrics:

```bash
curl http://localhost:7777/internal/v1/rabbitmq/metrics
```

Example files:

```text
examples/rabbitmq/dry-run-payment-events.yaml
examples/rabbitmq/payment-events.yaml
```
