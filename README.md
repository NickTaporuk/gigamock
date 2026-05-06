# gigamock

Gigamock is a production-ready mock server foundation for describing predictable
responses from configuration files. It can be used today for HTTP REST mocks,
GraphQL-over-HTTP mocks, Kafka producer/consumer scenarios, and as a control
plane for planned gRPC and broker-based mocks.

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
- `grpc`: production-ready YAML contract and UI indexing; native gRPC runtime
  is planned
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
  -dir-path ./config \
  -logger-level DEBUG \
  -logger-pretty-print=false
```

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

Example HTTP mock:

```yaml
path: "/users"
method: GET
type: http
description: "retrieve a list of active users"
scenarios:
  - name: "default scenario"
    request:
    response:
      headers:
        Content-Type: "application/json"
      statusCode: 200
      body: |
        {
          "users": [
            {
              "name": "Luke Skywalker"
            }
          ]
        }
  - name: "error 500"
    request:
    response:
      statusCode: 500
      body: |
        {
          "error": "internal error"
        }
```

The top-level fields are:

- `path`: endpoint path used for matching requests.
- `method`: HTTP method used for matching requests.
- `type`: scenario type, for example `http` or `kafka`.
- `description`: human-readable endpoint description shown in the control UI.
- `scenarios`: ordered list of available responses. Scenario index starts at
  `0`.

Detailed field references are available in
[`docs/scenario-types`](docs/scenario-types/README.md).

## Examples

Ready-to-read YAML examples are available in the `examples` directory:

- `examples/README.md`: overview and request examples for each scenario type.
- `examples/rest/control-ui-users.yaml`: HTTP mock with multiple responses for
  testing the control UI scenario switcher.
- `examples/graphql/starwars-operations.yaml`: GraphQL mock with
  `operationName`, `query`, and `variables` request matching.
- `examples/grpc/customer-service-unary.yaml`: planned production-ready unary
  gRPC mock format driven by protobuf descriptors.
- `examples/grpc/chat-service-bidi-stream.yaml`: planned bidirectional gRPC
  stream mock format with scripted receive/send steps.
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

## GraphQL Mocks

GraphQL mocks use `type: graphql` and are served over HTTP. They support the
same scenario switching model as HTTP mocks, plus optional request matching by
GraphQL payload fields:

- `operationName`
- `query`
- `variables`

Example request body:

```json
{
  "operationName": "GetHero",
  "query": "query GetHero($episode: String!) { hero(episode: $episode) { id name } }",
  "variables": {
    "episode": "NEWHOPE"
  }
}
```

When the currently active scenario does not match the incoming GraphQL payload,
Gigamock searches the other scenarios for a matching request definition and
returns that response. This allows several GraphQL operations to share one
endpoint such as `/graphql`.

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

## gRPC Mocking Direction

The recommended production-ready gRPC design is a dynamic mock engine driven by
`.proto` descriptors instead of generated Go mocks per service.

Target behavior:

- start a real gRPC server beside the existing HTTP server
- load service and message descriptors from `.proto` files or descriptor sets
- match calls by full method name, for example `/customers.CustomersService/Get`
- decode requests with protobuf reflection
- build responses from human-readable YAML/JSON using `protojson`
- support multiple scenarios per endpoint
- allow the control UI to switch active gRPC scenarios the same way it switches
  HTTP scenarios

Suggested future unary gRPC mock format:

```yaml
path: "/customers.CustomersService/GetCustomer"
method: POST
type: grpc
description: "mock GetCustomer gRPC method"
proto:
  descriptorSet: "./proto/customers.pb"
scenarios:
  - name: "default customer"
    request:
      match:
        customerId: "customer-1"
    response:
      code: OK
      metadata:
        x-mock: "gigamock"
      body:
        customer:
          id: "customer-1"
          name: "John Doe"
  - name: "not found"
    request:
      match:
        customerId: "missing"
    response:
      code: NOT_FOUND
      message: "customer not found"
```

For bidirectional streaming, the recommended format is a scripted state machine:

```yaml
path: "/chat.ChatService/Chat"
method: POST
type: grpc
description: "mock bidirectional chat stream"
proto:
  descriptorSet: "./proto/chat.pb"
scenarios:
  - name: "happy path"
    stream:
      steps:
        - receive:
            text: "hello"
        - send:
            text: "hi"
        - receive:
            text: "ping"
        - send:
            text: "pong"
        - close:
            code: OK
```

The current UI/control API is already shaped so these future `type: grpc`
scenario files can be listed and switched by endpoint once the gRPC runtime
engine is added.

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

- gRPC dynamic unary mock engine
- gRPC server-streaming, client-streaming, and bidirectional-streaming mocks
- NATS runtime publisher/subscriber provider
- RabbitMQ runtime publisher/consumer provider
- MQTT runtime publisher/subscriber provider
- Swagger/OpenAPI parser for generating mock scenarios
