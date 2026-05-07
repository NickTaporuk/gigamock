# gigamock

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/NickTaporuk/gigamock)](https://goreportcard.com/report/github.com/NickTaporuk/gigamock)

Gigamock is a production-ready mock server foundation for describing predictable
responses from configuration files. It can be used today for HTTP REST mocks,
GraphQL-over-HTTP mocks, dynamic gRPC mocks, scripted WebSocket mocks, Kafka
producer/consumer scenarios, and as a control plane for broker-based mocks.

Mock behavior is described in YAML or JSON files. At runtime Gigamock indexes
the configured files, serves mock responses, and exposes an internal control UI
where an active scenario can be changed without restarting the server.

## Production Readiness

Gigamock is designed to be production-ready as a mock control plane:

- deterministic YAML/JSON scenario files
- repeatable `--dir-path` configuration loading
- duplicate endpoint detection on startup
- runtime scenario switching through the built-in UI and internal API
- structured logging
- test coverage for core providers

Current runtime status by scenario type:

- `http`: production-ready HTTP response mocking
- `graphql`: production-ready GraphQL-over-HTTP response mocking with
  `operationName`, `query`, and `variables` matching
- `kafka`: runtime producer/consumer scenario support with local dry-run mode
- `grpc`: production-ready dynamic gRPC runtime for unary and scripted streaming
  mocks loaded from `.proto` files
- `nats`: runtime publish scenario support with local dry-run mode
- `rabbitmq`: runtime publish scenario support with local dry-run mode
- `mqtt`: runtime publish scenario support with local dry-run mode
- `websocket`: runtime scripted bidirectional WebSocket support with local
  dry-run mode

## Run

By default Gigamock reads mock files from `./config` and starts an HTTP server on
`0.0.0.0:7777`.

```bash
go run ./cmd
```

Available flags:

```bash
go run ./cmd \
  -server-ip 0.0.0.0 \
  -server-port :7777 \
  -grpc-server-port :7778 \
  -grpc-stream-max-messages 100 \
  -grpc-stream-timeout-seconds 300 \
  -dir-path ./config \
  -logger-level DEBUG \
  -logger-pretty-print=false
```

gRPC TLS/mTLS flags are also available. See
[gRPC scenario hardening settings](docs/scenario-types/grpc.md#hardening-settings).

`--dir-path` can be used multiple times when mocks are split across different
directories:

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

Print console help:

```bash
go run ./cmd --help
```

## Mock File Format

Mock files share a small common envelope:

- `path`: endpoint path used for matching requests.
- `method`: HTTP method used for matching requests.
- `type`: scenario type, for example `http` or `kafka`.
- `description`: human-readable endpoint description shown in the control UI.
- `scenarios`: ordered list of available responses. Scenario index starts at
  `0`.

Detailed field references are split by scenario type:

| Type | Documentation | Example |
| --- | --- | --- |
| HTTP | [HTTP scenarios](docs/scenario-types/http.md) | [`examples/rest`](examples/rest) |
| GraphQL | [GraphQL scenarios](docs/scenario-types/graphql.md) | [`examples/graphql`](examples/graphql) |
| gRPC | [gRPC scenarios](docs/scenario-types/grpc.md) | [`examples/grpc`](examples/grpc) |
| Kafka | [Kafka scenarios](docs/scenario-types/kafka.md) | [`examples/kafka`](examples/kafka) |
| NATS | [NATS scenarios](docs/scenario-types/nats.md) | [`examples/nats`](examples/nats) |
| RabbitMQ | [RabbitMQ scenarios](docs/scenario-types/rabbitmq.md) | [`examples/rabbitmq`](examples/rabbitmq) |
| MQTT | [MQTT scenarios](docs/scenario-types/mqtt.md) | [`examples/mqtt`](examples/mqtt) |
| WebSocket | [WebSocket scenarios](docs/scenario-types/websocket.md) | [`examples/websocket`](examples/websocket) |

The full documentation index is available in [`docs/README.md`](docs/README.md).

## Examples

Ready-to-read YAML examples are available in the `examples` directory:

- `examples/README.md`: overview and request examples for each scenario type.
- `examples/rest/control-ui-users.yaml`: HTTP mock with multiple responses for
  testing the control UI scenario switcher.
- `examples/graphql/starwars-operations.yaml`: GraphQL mock with
  `operationName`, `query`, and `variables` request matching.
- `examples/grpc/customer-service-unary.yaml`: real unary gRPC mock driven by a
  `.proto` file.
- `examples/grpc/billing-service-unary.yaml`: real billing unary gRPC mock.
- `examples/grpc/billing-service-server-stream.yaml`: real billing
  server-streaming gRPC mock.
- `examples/grpc/billing-service-bidi-stream.yaml`: real billing bidirectional
  gRPC mock.
- `examples/grpc/chat-service-bidi-stream.yaml`: real bidirectional gRPC stream
  mock with scripted receive/send steps.
- `examples/kafka/dry-run-topic.yaml`: Kafka producer dry-run scenario that
  works without a running Kafka broker.
- `examples/kafka/docker-topic.yaml`: Kafka producer scenario for the
  Docker Compose end-to-end stack.
- `examples/kafka/test-topic.yaml`: real Kafka producer/consumer scenario.
- `examples/nats/dry-run-order-created.yaml`: NATS dry-run publish scenario
  that works without a running NATS broker.
- `examples/nats/order-created.yaml`: real NATS publish scenario.
- `examples/rabbitmq/dry-run-payment-events.yaml`: RabbitMQ dry-run publish
  scenario that works without a running RabbitMQ broker.
- `examples/rabbitmq/payment-events.yaml`: real RabbitMQ publish scenario.
- `examples/mqtt/dry-run-device-telemetry.yaml`: MQTT dry-run publish scenario
  that works without a running MQTT broker.
- `examples/mqtt/device-telemetry.yaml`: real MQTT publish scenario.
- `examples/websocket/dry-run-chat.yaml`: WebSocket dry-run scenario.
- `examples/websocket/chat.yaml`: real scripted bidirectional WebSocket
  scenario.

## Acceptance Specs And Docker

Gherkin feature files are available in [`features`](features/README.md).

Per-type Dockerfiles are available in
[`deployment/docker/types`](deployment/docker/types/README.md).

Taskfile shortcuts are available in [`Taskfile.yml`](Taskfile.yml):

```bash
task run:all
task run:examples
task build
task docker:build:all
task docker:run:http
task features:list
```

`task build` writes the local binary to `./bin/gigamock`.

## Control UI

Gigamock exposes a small built-in UI for inspecting indexed mock files and
switching the active scenario for a specific endpoint.

Open:

```text
http://localhost:7777/internal/v1/mock-ui
```

The UI shows:

- endpoint path
- method
- scenario type
- description
- source YAML/JSON file
- active scenario selector
- runtime metrics for gRPC, GraphQL, Kafka, NATS, RabbitMQ, MQTT, and WebSocket
- live metrics updates without refreshing the page

Changing the selected scenario updates the in-memory route store. The mock
server starts using the selected scenario immediately.

## Scenario Documentation

Use the scenario type documentation for field-level details and request
examples:

- [HTTP scenarios](docs/scenario-types/http.md)
- [GraphQL scenarios](docs/scenario-types/graphql.md)
- [gRPC scenarios](docs/scenario-types/grpc.md)
- [Kafka scenarios](docs/scenario-types/kafka.md)
- [NATS scenarios](docs/scenario-types/nats.md)
- [RabbitMQ scenarios](docs/scenario-types/rabbitmq.md)
- [MQTT scenarios](docs/scenario-types/mqtt.md)
- [WebSocket scenarios](docs/scenario-types/websocket.md)

## Internal API

List raw in-memory route state:

```bash
curl http://localhost:7777/internal/v1/in-memory
```

List UI-friendly scenario details:

```bash
curl http://localhost:7777/internal/v1/scenarios
```

Kafka dry-run smoke test:

```bash
curl http://localhost:7777/internal/kafka/dry-run/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
```

Kafka Docker end-to-end smoke test:

```bash
task docker:kafka:up
```

Then in another terminal:

```bash
curl http://localhost:7777/internal/kafka/docker/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
```

Stop the stack:

```bash
task docker:kafka:down
```

If `7777` is busy:

```bash
PORT=7781 task docker:kafka:up
PORT=7781 task docker:kafka:test
PORT=7781 task docker:kafka:down
```

NATS dry-run smoke test:

```bash
curl -X POST http://localhost:7777/internal/nats/dry-run/orders/order-1 \
  -H "Content-Type: application/json" \
  -d '{"orderId":"order-1"}'
curl http://localhost:7777/internal/v1/nats/metrics
```

RabbitMQ dry-run smoke test:

```bash
curl -X POST http://localhost:7777/internal/rabbitmq/dry-run/payments/payment-1 \
  -H "Content-Type: application/json" \
  -d '{"paymentId":"payment-1"}'
curl http://localhost:7777/internal/v1/rabbitmq/metrics
```

MQTT dry-run smoke test:

```bash
curl -X POST http://localhost:7777/internal/mqtt/dry-run/devices/device-1/telemetry \
  -H "Content-Type: application/json" \
  -d '{"deviceId":"device-1"}'
curl http://localhost:7777/internal/v1/mqtt/metrics
```

WebSocket dry-run smoke test:

```bash
curl http://localhost:7777/ws/dry-run/chat
curl http://localhost:7777/internal/v1/websocket/metrics
```

Real WebSocket smoke test:

```bash
printf '{"sender":"client","text":"ping"}\n' | websocat ws://localhost:7777/ws/chat
```

Set active scenario for an endpoint:

```bash
curl -X POST http://localhost:7777/internal/v1/in-memory \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/users",
    "method": "GET",
    "scenarioNumber": 1
  }'
```

## Download

## Precompiled Binaries

You can download the precompiled release binary from [releases](https://github.com/NickTaporuk/gigamock/releases/) via web
or via

```bash
wget https://github.com/NickTaporuk/gigamock/releases/<version>/gigamock_<version>_<os>_<arch>
```

#### Go get

You can also use Go 1.12 or later to build the latest stable version from source:

```bash
go get github.com/NickTaporuk/gigamock
```

#### Homebrew Tap

```bash
brew install nicktaporuk/tap/gigamock
# After initial install you can upgrade the version via:
brew upgrade gigamock
```

## Roadmap

- Swagger/OpenAPI parser for generating mock scenarios

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_large)
