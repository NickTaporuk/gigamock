# Kafka Scenario Fields

Kafka scenarios can prepare a topic, publish messages, and optionally run a
consumer logger.

Example:

```yaml
path: "/internal/queue/:messageID"
method: GET
type: kafka
scenarios:
  - name: "default scenario send kafka message to the topic test-topic"
    host: "0.0.0.0"
    port: "9092"
    topic: "test-topic"
    delay: 100s
    producer:
      partition: 1
      headers:
        X-Request-Id: "831429af-1e40-4b44-8be3-06fd252578b0"
      message:
        value: "{\"test\":\"test\"}"
        key: test
      retry: 1
    consumer:
      cli: true
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route that triggers the Kafka scenario. |
| `method` | yes | HTTP method that triggers the Kafka scenario. |
| `type` | yes | Must be `kafka`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of Kafka scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `host` | yes | Kafka host. |
| `port` | yes | Kafka port. |
| `topic` | yes | Kafka topic. |
| `delay` | no | Planned delay between operations. |
| `producer` | no | Producer configuration. |
| `consumer` | no | Consumer configuration. |

Producer fields:

| Field | Required | Description |
| --- | --- | --- |
| `partition` | no | Kafka partition. |
| `headers` | no | Kafka message headers. |
| `message` | yes | Message key/value payload. |
| `retry` | no | Planned retry count. |

Message fields:

| Field | Required | Description |
| --- | --- | --- |
| `key` | no | Kafka message key. |
| `value` | yes | Kafka message value. |

Consumer fields:

| Field | Required | Description |
| --- | --- | --- |
| `cli` | no | Whether to log consumed messages to the CLI. |

Example files:

```text
examples/kafka/test-topic.yaml
examples/kafka/test-duplicate-consumer-topic.yaml
```
