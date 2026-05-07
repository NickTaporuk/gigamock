# Google Pub/Sub Scenario Fields

Pub/Sub scenarios provide a small Google Pub/Sub-compatible in-memory API for
local and CI integration tests.

Supported routes:

```text
POST   /gcp/pubsub/:topic/publish
POST   /gcp/pubsub/:subscription/pull
DELETE /gcp/pubsub/:subscription/purge
POST   /gcp/pubsub-dry-run/:topic/publish
```

Example:

```yaml
path: "/gcp/pubsub/:topic/publish"
method: POST
type: pubsub
description: "Google Pub/Sub-compatible Publish mock"
scenarios:
  - name: "publish message"
    attributes:
      source: "gigamock"
```

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `topic` | no | Topic override. Defaults to the topic in the publish path. Pull and purge scenarios should set it when subscription name differs from topic name. |
| `subscription` | no | Subscription override. Defaults to the subscription in the pull or purge path. |
| `dryRun` | no | When `true`, returns a JSON response without changing topic state. |
| `message` | no | Default publish body when the request body is empty. |
| `attributes` | no | Message attributes returned with the published or pulled message. |

Smoke test:

```bash
curl -X POST http://localhost:7777/gcp/pubsub/order-events/publish \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl -X POST http://localhost:7777/gcp/pubsub/orders-sub/pull
curl -X DELETE http://localhost:7777/gcp/pubsub/orders-sub/purge
curl http://localhost:7777/internal/v1/pubsub/metrics
```

Dry-run check:

```bash
curl -X POST http://localhost:7777/gcp/pubsub-dry-run/order-events/publish
```

Example files:

```text
examples/pubsub/publish-message.yaml
examples/pubsub/pull-message.yaml
examples/pubsub/purge-subscription.yaml
examples/pubsub/dry-run.yaml
```
