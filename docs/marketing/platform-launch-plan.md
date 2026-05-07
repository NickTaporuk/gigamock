# Gigamock Launch Plan

## Best Platforms

Use different versions of the same story rather than copy-pasting one post
everywhere.

| Platform | Best angle | Format |
| --- | --- | --- |
| Dev.to | Practical tutorial and examples | Full article |
| Hashnode | Founder/building-in-public story | Full article |
| Medium | Higher-level engineering article | Full article |
| Hacker News | "Show HN: Gigamock..." with concise technical context | Short launch post |
| Reddit | Specific problem/lesson, not direct advertising | Discussion post |
| LinkedIn | Product story and business value | Short post with image/screenshot |
| X/Twitter | Thread with protocol list and screenshots | 6-8 posts |

## Title Options

- I Built Gigamock: A Mock Server For Multi-Protocol Systems
- Mocking Is Easy Until Your System Uses REST, gRPC, Kafka, S3, SOAP, And More
- Show HN: Gigamock, a Source-Available Mock Server for Multi-Protocol Backends
- One Mock Control Plane For HTTP, GraphQL, gRPC, Brokers, Cloud Queues, And SOAP
- Testing Distributed Systems Without Running Every Dependency Locally

## Short Description

Gigamock is a source-available mock server for describing HTTP, GraphQL, gRPC,
broker, cloud-compatible, WebSocket, S3, and SOAP scenarios in YAML or JSON,
with a UI for switching scenarios and viewing runtime metrics.

## Tags

```text
mocking, testing, golang, microservices, grpc, graphql, kafka, devtools,
backend, qa, ci, contract-testing
```

## Hacker News Draft

```text
Show HN: Gigamock, a source-available mock server for multi-protocol backends

I built Gigamock because most mock setups I saw worked fine for HTTP, then got
messy once a system added gRPC, brokers, WebSocket, S3-compatible storage, cloud
queues, or legacy SOAP.

Gigamock loads YAML/JSON mock descriptions from one or more directories, indexes
them, exposes a small UI, lets you switch active scenarios at runtime, and
tracks runtime metrics.

Currently supported:
HTTP, GraphQL, real gRPC from .proto, Kafka, NATS, RabbitMQ, MQTT, WebSocket,
S3-compatible API, AWS SQS/SNS-compatible APIs, Google Pub/Sub-compatible API,
Azure Service Bus-compatible API, and SOAP.

Repo: https://github.com/NickTaporuk/gigamock

I would appreciate feedback from people doing integration testing for
multi-protocol services.
```

## Reddit Draft

```text
Title: How do you keep mocks manageable when a backend uses more than HTTP?

I have been working on a Go tool called Gigamock because I kept running into a
testing problem: HTTP mocks are easy, but modern services often also use gRPC,
Kafka/NATS/RabbitMQ/MQTT, WebSocket, S3-compatible APIs, cloud queues, and
sometimes SOAP.

The idea is to describe scenarios in YAML/JSON, load multiple directories, show
the indexed mocks in a UI, switch active scenarios at runtime, and expose
metrics for what was actually called.

I am curious how other teams structure this. Do you prefer one mock control
plane, protocol-specific emulators, generated mocks from contracts, or a mix?

Repo if useful: https://github.com/NickTaporuk/gigamock
```

## LinkedIn Draft

```text
I built Gigamock, a source-available mock server for teams testing
multi-protocol backend systems.

Most products no longer depend only on REST. A single workflow can involve
HTTP, GraphQL, gRPC, Kafka, cloud queues, WebSocket, S3-compatible storage, and
even SOAP.

Gigamock lets teams describe scenarios in YAML/JSON, load multiple directories,
switch active responses from a UI, and inspect runtime metrics.

Supported today:
HTTP, GraphQL, gRPC, Kafka, NATS, RabbitMQ, MQTT, WebSocket, S3, AWS SQS/SNS,
Google Pub/Sub, Azure Service Bus, and SOAP.

I would love feedback from backend, QA, platform, and developer-tooling teams.

https://github.com/NickTaporuk/gigamock
```

## Launch Checklist

- Add screenshots of the control UI to the README.
- Add a short GIF showing scenario switching.
- Pin a "Quick Start" issue or discussion on GitHub.
- Add GitHub topics: `mock-server`, `testing`, `golang`, `grpc`, `graphql`,
  `kafka`, `microservices`, `contract-testing`, `qa`, `devtools`.
- Publish the full article on Dev.to or Hashnode first.
- Post the Hacker News version as `Show HN`.
- Share the Reddit version as a discussion, not a pure promo link.
- Ask 3-5 developers to star, comment, or give feedback on launch day.
- Respond to every comment in the first 24 hours.
- After feedback, create issues for the top requested features.

## Important Positioning

Do not call Gigamock "open source" unless the license changes to an OSI
approved license. Use "source-available" instead.

Avoid saying it replaces LocalStack, WireMock, grpcmock, broker emulators, or
contract-testing tools. Better positioning:

```text
Gigamock is one practical mock control plane for teams that need predictable
multi-protocol test scenarios.
```
