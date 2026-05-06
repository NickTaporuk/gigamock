# Gigamock Examples

This directory contains YAML examples for each mock scenario type currently
accepted by Gigamock.

Run all examples together:

```bash
go run ./cmd \
  --dir-path ./examples/rest \
  --dir-path ./examples/graphql \
  --dir-path ./examples/grpc \
  --dir-path ./examples/kafka \
  --dir-path ./examples/nats \
  --dir-path ./examples/rabbitmq \
  --dir-path ./examples/mqtt
```

Open the control UI:

```text
http://localhost:7777/internal/v1/mock-ui
```

## HTTP

File:

```text
examples/rest/control-ui-users.yaml
```

Use it to test normal HTTP response switching from the control UI.

Request:

```bash
curl http://localhost:7777/control-ui/users
```

## GraphQL

File:

```text
examples/graphql/starwars-operations.yaml
```

Use it to mock multiple GraphQL operations on one HTTP endpoint. Matching can
use `operationName`, `query`, and `variables`.

Request:

```bash
curl -X POST http://localhost:7777/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "operationName": "GetHero",
    "query": "query GetHero($episode: String!) { hero(episode: $episode) { id name } }",
    "variables": {
      "episode": "NEWHOPE"
    }
  }'
```

## gRPC

Files:

```text
examples/grpc/customer-service-unary.yaml
examples/grpc/chat-service-bidi-stream.yaml
```

These files define real gRPC mocks. They are indexed, shown in the control UI,
and served from the gRPC listener.

```bash
grpcurl -plaintext \
  -d '{"customerId":"customer-1"}' \
  localhost:7778 \
  customers.CustomersService/GetCustomer
```

## Kafka

File:

```text
examples/kafka/test-topic.yaml
```

This is the current Kafka producer/consumer scenario format.

Request that triggers the configured Kafka scenario:

```bash
curl http://localhost:7777/internal/queue/message-1
```

## NATS

File:

```text
examples/nats/order-created.yaml
```

This file defines the planned NATS publish scenario format. It is indexed and
shown in the control UI. The real NATS runtime is still a future implementation
step.

## RabbitMQ

File:

```text
examples/rabbitmq/payment-events.yaml
```

This file defines the planned RabbitMQ publish scenario format. It is indexed
and shown in the control UI. The real RabbitMQ runtime is still a future
implementation step.

## MQTT

File:

```text
examples/mqtt/device-telemetry.yaml
```

This file defines the planned MQTT publish scenario format. It is indexed and
shown in the control UI. The real MQTT runtime is still a future implementation
step.
