# Scenario Type Reference

This directory documents every scenario type accepted by Gigamock.

Common top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route or logical endpoint used by Gigamock to index the scenario. |
| `method` | yes | HTTP method used for route matching and UI switching. |
| `type` | yes | Scenario type: `http`, `graphql`, `grpc`, `kafka`, `nats`, `rabbitmq`, `mqtt`, `websocket`, `s3`, `sqs`, `sns`, `pubsub`, `servicebus`, or `soap`. |
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
| S3 | [S3](s3.md) | `bucket`, `key`, `contentType`, `body`, `metadata` |
| SQS | [SQS](sqs.md) | `queue`, `message`, `attributes` |
| SNS | [SNS](sns.md) | `topic`, `message`, `subject`, `attributes` |
| Google Pub/Sub | [Pub/Sub](pubsub.md) | `topic`, `subscription`, `message`, `attributes` |
| Azure Service Bus | [Service Bus](servicebus.md) | `queue`, `topic`, `subscription`, `message`, `properties` |
| SOAP | [SOAP](soap.md) | `request.soapAction`, `request.bodyContains`, `response.body` |

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
| `s3` | Runtime in-memory S3-compatible object API implemented. |
| `sqs` | Runtime in-memory SQS-compatible queue API implemented. |
| `sns` | Runtime in-memory SNS-compatible topic API implemented. |
| `pubsub` | Runtime in-memory Google Pub/Sub-compatible topic/subscription API implemented. |
| `servicebus` | Runtime in-memory Azure Service Bus-compatible queue API implemented. |
| `soap` | Runtime SOAP XML-over-HTTP response mocking implemented. |

See the repository documentation index in [../README.md](../README.md).
