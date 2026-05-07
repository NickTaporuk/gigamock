# Start Gigamock

The main documentation index is available here:

- [Documentation index](README.md)
- [Scenario type overview](scenario-types/README.md)
- [Example YAML files](../examples/README.md)
- [Manual testing requests](../testsing/requests/README.md)

Run all examples:

```bash
go run ./cmd --dir-path ./examples
```

Run examples split by type:

```bash
go run ./cmd \
  --dir-path ./examples/rest \
  --dir-path ./examples/graphql \
  --dir-path ./examples/grpc \
  --dir-path ./examples/kafka \
  --dir-path ./examples/nats \
  --dir-path ./examples/rabbitmq \
  --dir-path ./examples/mqtt \
  --dir-path ./examples/websocket
```

Useful URLs:

```text
http://localhost:7777/internal/v1/mock-ui
http://localhost:7777/internal/v1/scenarios
```

The mock UI includes endpoint scenario switching and live runtime metrics for
gRPC, GraphQL, Kafka, NATS, RabbitMQ, MQTT, and WebSocket.

Real gRPC endpoint:

```bash
grpcurl -plaintext localhost:7778 list
```

gRPC runtime metrics:

```bash
curl http://localhost:7777/internal/v1/grpc/metrics
```

GraphQL runtime metrics:

```bash
curl http://localhost:7777/internal/v1/graphql/metrics
```

Broker runtime metrics:

```bash
curl http://localhost:7777/internal/v1/kafka/metrics
curl http://localhost:7777/internal/v1/nats/metrics
curl http://localhost:7777/internal/v1/rabbitmq/metrics
curl http://localhost:7777/internal/v1/mqtt/metrics
curl http://localhost:7777/internal/v1/websocket/metrics
```

Kafka Docker end-to-end stack:

```bash
task docker:kafka:up
curl http://localhost:7777/internal/kafka/docker/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
task docker:kafka:down
```
