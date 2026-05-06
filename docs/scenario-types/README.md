# Scenario Type Reference

This directory documents every scenario type accepted by Gigamock.

Common top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route or logical endpoint used by Gigamock to index the scenario. |
| `method` | yes | HTTP method used for route matching and UI switching. |
| `type` | yes | Scenario type: `http`, `graphql`, `grpc`, `kafka`, `nats`, `rabbitmq`, or `mqtt`. |
| `description` | no | Human-readable endpoint description shown in the control UI. |
| `scenarios` | yes | Ordered list of selectable scenarios. Scenario index starts at `0`. |
| `webhook` | no | Optional webhook executed after scenario retrieval where supported. |

Scenario type docs:

- [HTTP](http.md)
- [GraphQL](graphql.md)
- [gRPC](grpc.md)
- [Kafka](kafka.md)
- [NATS](nats.md)
- [RabbitMQ](rabbitmq.md)
- [MQTT](mqtt.md)

Runtime status:

| Type | Status |
| --- | --- |
| `http` | Runtime provider implemented. |
| `graphql` | Runtime provider implemented over HTTP. |
| `kafka` | Runtime producer/consumer provider implemented. |
| `grpc` | YAML contract and UI indexing implemented; native gRPC runtime planned. |
| `nats` | YAML contract and UI indexing implemented; native NATS runtime planned. |
| `rabbitmq` | YAML contract and UI indexing implemented; native RabbitMQ runtime planned. |
| `mqtt` | YAML contract and UI indexing implemented; native MQTT runtime planned. |
