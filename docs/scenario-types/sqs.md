# SQS Scenario Fields

SQS scenarios provide a small SQS-compatible in-memory queue API for local and
CI integration tests.

Supported routes:

```text
POST   /aws/sqs/:queue
GET    /aws/sqs/:queue
DELETE /aws/sqs/:queue
POST   /aws/sqs-dry-run/:queue
```

Example:

```yaml
path: "/aws/sqs/:queue"
method: POST
type: sqs
description: "SQS-compatible SendMessage mock"
scenarios:
  - name: "send message"
    attributes:
      source: "gigamock"
```

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `queue` | no | Queue override. Defaults to the queue in the request path. |
| `dryRun` | no | When `true`, returns a JSON response without changing queue state. |
| `message` | no | Default message body when the request body is empty. |
| `attributes` | no | Message attributes returned with the message. |

Smoke test:

```bash
curl -X POST http://localhost:7777/aws/sqs/orders \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl http://localhost:7777/aws/sqs/orders
curl -X DELETE http://localhost:7777/aws/sqs/orders
curl http://localhost:7777/internal/v1/sqs/metrics
```

Dry-run check:

```bash
curl -X POST http://localhost:7777/aws/sqs-dry-run/orders
```

Example files:

```text
examples/sqs/send-message.yaml
examples/sqs/receive-message.yaml
examples/sqs/purge-queue.yaml
examples/sqs/dry-run.yaml
```
