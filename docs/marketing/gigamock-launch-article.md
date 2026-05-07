# I Built Gigamock: A Production-Ready Mock Server For Multi-Protocol Systems

Modern backend testing used to be mostly HTTP stubs.

That world is gone.

Today one feature can touch REST, GraphQL, gRPC, Kafka, NATS, RabbitMQ, MQTT,
WebSocket, S3-compatible storage, AWS SQS/SNS, Google Pub/Sub, Azure Service
Bus, and sometimes a SOAP endpoint that has been quietly paying the bills since
2011.

The problem is not that mocking is impossible. The problem is that teams often
mock each protocol with a different tool, a different config style, and a
different mental model. After a few services, the test environment becomes a
small distributed system of its own.

Gigamock is my attempt to make that simpler.

## What Gigamock Is

Gigamock is a source-available mock server for describing predictable service
behavior in YAML or JSON files.

You run it, point it at one or more directories, and it indexes your mock
descriptions:

```bash
go run ./cmd \
  --dir-path ./examples/rest \
  --dir-path ./examples/graphql \
  --dir-path ./examples/grpc \
  --dir-path ./examples/kafka \
  --dir-path ./examples/s3 \
  --dir-path ./examples/sqs \
  --dir-path ./examples/pubsub \
  --dir-path ./examples/azure-servicebus \
  --dir-path ./examples/soap
```

Each mock endpoint can have multiple scenarios. You can switch the active
scenario from the built-in UI or from the internal API without restarting the
server.

That means one test environment can simulate:

- a normal user response;
- a validation error;
- a downstream timeout;
- an empty queue;
- a SOAP fault;
- a gRPC stream with scripted messages;
- an object uploaded to S3-compatible storage;
- a broker message being produced or consumed.

## Why I Built It

I wanted a mock server that behaves like infrastructure, not like a pile of
one-off fixtures.

The goals were:

- keep mocks as reviewable config files;
- support many protocols with one control plane;
- make scenarios visible and switchable in a UI;
- support local and CI testing;
- provide examples that a developer can copy and adapt;
- expose runtime metrics so you can see what was actually called.

The big idea is simple: the mock description should live next to the system
contract, and the running mock server should tell you exactly what it loaded and
what happened at runtime.

## Example: HTTP Scenario

```yaml
path: "/control-ui/users"
method: GET
type: http
description: "HTTP example with multiple scenarios"
scenarios:
  - name: "active users"
    response:
      statusCode: 200
      body: |
        [{"id":"user-1","name":"Ada Lovelace"}]
  - name: "temporary error"
    response:
      statusCode: 503
      body: |
        {"error":"temporary unavailable"}
```

## Example: gRPC Scenario

Gigamock can load real gRPC mocks driven by `.proto` files.

```bash
grpcurl -plaintext \
  -d '{"customerId":"customer-1"}' \
  localhost:7778 \
  customers.CustomersService/GetCustomer
```

It supports unary and scripted streaming flows, including bidirectional streams.

## Example: SOAP Scenario

SOAP is still alive in many enterprise systems, so Gigamock can also match
`SOAPAction` and request XML fragments:

```yaml
path: "/soap/customers"
method: POST
type: soap
description: "SOAP customer service mock"
scenarios:
  - name: "get customer"
    request:
      soapAction: "GetCustomer"
      bodyContains: "<customerId>customer-1</customerId>"
    response:
      statusCode: 200
      headers:
        Content-Type: "text/xml; charset=utf-8"
      body: |
        <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
          <soap:Body>
            <GetCustomerResponse xmlns="urn:gigamock:customers">
              <customer>
                <id>customer-1</id>
                <name>Ada Lovelace</name>
              </customer>
            </GetCustomerResponse>
          </soap:Body>
        </soap:Envelope>
```

## The Control UI

The UI shows loaded endpoints, source files, scenario names, active scenarios,
service directories, and runtime metrics.

This matters because mocks are only useful if the team can understand them.
When a QA engineer, backend developer, or frontend developer asks "which
response am I getting right now?", the answer should not require reading a
stack trace or restarting a container.

## Supported Scenario Types

Gigamock currently supports:

- HTTP REST
- GraphQL-over-HTTP
- real gRPC from `.proto`
- Kafka
- NATS
- RabbitMQ
- MQTT
- WebSocket
- S3-compatible API
- AWS SQS/SNS-compatible APIs
- Google Pub/Sub-compatible API
- Azure Service Bus-compatible API
- SOAP

Some broker/cloud-compatible mocks are intentionally in-memory or dry-run
friendly, because local and CI tests should not always require a full external
broker stack.

## Where This Helps

Gigamock is useful when:

- frontend teams need stable backend responses;
- backend teams need to test failure paths;
- QA needs repeatable scenario switching;
- microservices depend on multiple protocols;
- CI pipelines need deterministic test infrastructure;
- legacy SOAP and modern event-driven services live in the same company.

It is not trying to replace every specialized emulator. It is trying to provide
one practical control plane for the mocks a product team actually needs day to
day.

## What I Learned

The hardest part of mocking is not returning JSON.

The hard part is keeping mocks understandable when the system grows.

That is why I focused on:

- file provenance in the UI;
- multiple directories via repeated `--dir-path`;
- live metrics;
- clear scenario descriptions;
- examples for every supported type;
- request files for manual smoke testing;
- Dockerfiles per protocol type.

## Try It

```bash
git clone https://github.com/NickTaporuk/gigamock
cd gigamock
go run ./cmd --dir-path ./examples
```

Then open:

```text
http://localhost:7777/internal/v1/mock-ui
```

If you are building or testing a distributed system with too many protocols and
too many fragile test dependencies, Gigamock may save you some of that pain.

Feedback, issues, and real-world use cases are very welcome.

Repository:

```text
https://github.com/NickTaporuk/gigamock
```

License note: Gigamock is source-available. Non-commercial use is free, and
commercial use requires a revenue share under the repository license.
