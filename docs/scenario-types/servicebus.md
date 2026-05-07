# Azure Service Bus Scenario Fields

Azure Service Bus scenarios provide a small in-memory queue API for local and CI
integration tests.

Supported routes:

```text
POST   /azure/servicebus/:queue/send
POST   /azure/servicebus/:queue/receive
DELETE /azure/servicebus/:queue/purge
POST   /azure/servicebus-dry-run/:queue/send
```

Example:

```yaml
path: "/azure/servicebus/:queue/send"
method: POST
type: servicebus
description: "Azure Service Bus-compatible SendMessage mock"
scenarios:
  - name: "send message"
    properties:
      source: "gigamock"
```

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `queue` | no | Queue override. Defaults to the queue in the request path. |
| `topic` | no | Optional topic name for topic-like scenarios. |
| `subscription` | no | Optional subscription name returned on received messages. |
| `dryRun` | no | When `true`, returns a JSON response without changing queue state. |
| `message` | no | Default message body when the request body is empty. |
| `properties` | no | Application properties returned with the message. |

Smoke test:

```bash
curl -X POST http://localhost:7777/azure/servicebus/orders/send \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl -X POST http://localhost:7777/azure/servicebus/orders/receive
curl -X DELETE http://localhost:7777/azure/servicebus/orders/purge
curl http://localhost:7777/internal/v1/servicebus/metrics
```

Dry-run check:

```bash
curl -X POST http://localhost:7777/azure/servicebus-dry-run/orders/send
```

Example files:

```text
examples/azure-servicebus/send-message.yaml
examples/azure-servicebus/receive-message.yaml
examples/azure-servicebus/purge-queue.yaml
examples/azure-servicebus/dry-run.yaml
```
