# SNS Scenario Fields

SNS scenarios provide a small SNS-compatible in-memory topic API for local and
CI integration tests.

Supported routes:

```text
POST /aws/sns/:topic
GET  /aws/sns/:topic
POST /aws/sns-dry-run/:topic
```

Example:

```yaml
path: "/aws/sns/:topic"
method: POST
type: sns
description: "SNS-compatible Publish mock"
scenarios:
  - name: "publish message"
    subject: "order event"
    attributes:
      source: "gigamock"
```

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `topic` | no | Topic override. Defaults to the topic in the request path. |
| `dryRun` | no | When `true`, returns a JSON response without changing topic state. |
| `message` | no | Default publish body when the request body is empty. |
| `subject` | no | Message subject. |
| `attributes` | no | Message attributes returned with the published message. |

Smoke test:

```bash
curl -X POST http://localhost:7777/aws/sns/order-events \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl http://localhost:7777/aws/sns/order-events
curl http://localhost:7777/internal/v1/sns/metrics
```

Dry-run check:

```bash
curl -X POST http://localhost:7777/aws/sns-dry-run/order-events
```

Example files:

```text
examples/sns/publish-message.yaml
examples/sns/list-messages.yaml
examples/sns/dry-run.yaml
```
