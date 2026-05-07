# Scenario Type Reference

This directory documents every scenario type accepted by Gigamock.

Common top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route or logical endpoint used by Gigamock to index the scenario. |
| `method` | yes | HTTP method used for route matching and UI switching. |
| `type` | yes | Scenario type: `http`, `graphql`, `grpc`, `kafka`, `nats`, `rabbitmq`, `mqtt`, or `websocket`. |
| `description` | no | Human-readable endpoint description shown in the control UI. |
| `scenarios` | yes | Ordered list of selectable scenarios. Scenario index starts at `0`. |
| `webhook` | no | Optional webhook executed after scenario retrieval where supported. |

Scenario type docs:

| Type | Documentation | Main Fields |
| --- | --- | --- |
| HTTP | [HTTP](http.md) | `request`, `response`, `statusCode`, `headers`, `body` |
| GraphQL | [GraphQL](graphql.md) | `operationName`, `query`, `variables`, `response` |
| gRPC | [gRPC](grpc.md) | `proto`, `request.match`, `response`, `stream` |
| Kafka | [Kafka](kafka.md) | `host`, `port`, `topic`, `producer`, `consumer` |
| NATS | [NATS](nats.md) | `host`, `subject`, `headers`, `message` |
| RabbitMQ | [RabbitMQ](rabbitmq.md) | `url`, `exchange`, `routingKey`, `message` |
| MQTT | [MQTT](mqtt.md) | `broker`, `clientID`, `topic`, `qos`, `message` |
| WebSocket | [WebSocket](websocket.md) | `sendOnConnect`, `steps`, `receive`, `send`, `close` |

Runtime status:

| Type | Status |
| --- | --- |
| `http` | Runtime provider implemented. |
| `graphql` | Runtime provider implemented over HTTP. |
| `kafka` | Runtime producer/consumer provider implemented. |
| `grpc` | Dynamic unary and scripted streaming gRPC runtime implemented from `.proto` files. |
| `nats` | Runtime publish support implemented with local dry-run mode. |
| `rabbitmq` | Runtime publish support implemented with local dry-run mode. |
| `mqtt` | Runtime publish support implemented with local dry-run mode. |
| `websocket` | Runtime scripted bidirectional support implemented with local dry-run mode. |

See the repository documentation index in [../README.md](../README.md).
