# RabbitMQ Scenario Fields

RabbitMQ scenarios define the planned RabbitMQ publish/consume mock contract.
Gigamock currently indexes these files and displays them in the control UI.
Native RabbitMQ runtime support is planned.

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
| `url` | planned | RabbitMQ AMQP connection URL. |
| `exchange` | planned | Exchange name. |
| `routingKey` | planned | Routing key. |
| `headers` | planned | Message headers. |
| `message.contentType` | planned | Message content type. |
| `message.body` | planned | Message body. |

Example file:

```text
examples/rabbitmq/payment-events.yaml
```
