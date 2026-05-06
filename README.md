# gigamock

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/NickTaporuk/gigamock)](https://goreportcard.com/report/github.com/NickTaporuk/gigamock)

Gigamock is a production-ready mock server foundation for describing predictable
responses from configuration files. It can be used today for HTTP REST mocks,
GraphQL-over-HTTP mocks, dynamic gRPC mocks, Kafka producer/consumer scenarios,
and as a control plane for broker-based mocks.

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
- `kafka`: runtime producer/consumer scenario support
- `grpc`: production-ready dynamic gRPC runtime for unary and scripted streaming
  mocks loaded from `.proto` files
- `nats`: production-ready YAML contract and UI indexing; native NATS runtime is
  planned
- `rabbitmq`: production-ready YAML contract and UI indexing; native RabbitMQ
  runtime is planned
- `mqtt`: production-ready YAML contract and UI indexing; native MQTT runtime is
  planned

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
  --dir-path ./examples/mqtt
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
- `examples/kafka/test-topic.yaml`: Kafka producer/consumer scenario.
- `examples/nats/order-created.yaml`: planned NATS publish scenario format.
- `examples/rabbitmq/payment-events.yaml`: planned RabbitMQ publish scenario
  format.
- `examples/mqtt/device-telemetry.yaml`: planned MQTT publish scenario format.

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

## Internal API

List raw in-memory route state:

```bash
curl http://localhost:7777/internal/v1/in-memory
```

List UI-friendly scenario details:

```bash
curl http://localhost:7777/internal/v1/scenarios
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

- NATS runtime publisher/subscriber provider
- RabbitMQ runtime publisher/consumer provider
- MQTT runtime publisher/subscriber provider
- Swagger/OpenAPI parser for generating mock scenarios

## License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FNickTaporuk%2Fgigamock?ref=badge_large)
