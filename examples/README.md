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
  --dir-path ./examples/mqtt \
  --dir-path ./examples/websocket \
  --dir-path ./examples/s3 \
  --dir-path ./examples/sqs \
  --dir-path ./examples/sns \
  --dir-path ./examples/pubsub \
  --dir-path ./examples/azure-servicebus \
  --dir-path ./examples/soap
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
examples/grpc/billing-service-unary.yaml
examples/grpc/billing-service-server-stream.yaml
examples/grpc/billing-service-client-stream.yaml
examples/grpc/billing-service-bidi-stream.yaml
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

Billing unary, server-streaming, and bidirectional examples:

```bash
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  localhost:7778 \
  billing.BillingService/GetInvoice

grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  localhost:7778 \
  billing.BillingService/WatchInvoice

grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1","status":"OPEN","message":"invoice created"}{"invoiceId":"invoice-1","status":"PAID","message":"invoice paid"}' \
  localhost:7778 \
  billing.BillingService/UploadInvoiceEvents

grpcurl -plaintext \
  -d '{"text":"hello","sender":"client"}' \
  localhost:7778 \
  billing.BillingService/BillingChat
```

## Kafka

Files:

```text
examples/kafka/dry-run-topic.yaml
examples/kafka/docker-topic.yaml
examples/kafka/test-topic.yaml
```

`dry-run-topic.yaml` is useful for local smoke tests and CI because it does not
require a running Kafka broker. `test-topic.yaml` uses the real producer/consumer
runtime and requires Kafka on the configured host/port.

Request that works without Kafka:

```bash
curl http://localhost:7777/internal/kafka/dry-run/message-1
```

Request that triggers the real Kafka scenario:

```bash
curl http://localhost:7777/internal/queue/message-1
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/kafka/metrics
```

Full Docker flow with Kafka broker and Gigamock:

```bash
task docker:kafka:up
```

Then in another terminal:

```bash
curl http://localhost:7777/internal/kafka/docker/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
```

## NATS

Files:

```text
examples/nats/dry-run-order-created.yaml
examples/nats/order-created.yaml
```

`dry-run-order-created.yaml` works without a running NATS broker. `order-created.yaml`
uses the real NATS publish runtime and requires a broker on the configured URL.

Request that works without NATS:

```bash
curl -X POST http://localhost:7777/internal/nats/dry-run/orders/order-1 \
  -H "Content-Type: application/json" \
  -d '{"orderId":"order-1"}'
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/nats/metrics
```

## RabbitMQ

Files:

```text
examples/rabbitmq/dry-run-payment-events.yaml
examples/rabbitmq/payment-events.yaml
```

`dry-run-payment-events.yaml` works without a running RabbitMQ broker.
`payment-events.yaml` uses the real RabbitMQ publish runtime and requires a
broker on the configured URL.

Request that works without RabbitMQ:

```bash
curl -X POST http://localhost:7777/internal/rabbitmq/dry-run/payments/payment-1 \
  -H "Content-Type: application/json" \
  -d '{"paymentId":"payment-1"}'
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/rabbitmq/metrics
```

## MQTT

Files:

```text
examples/mqtt/dry-run-device-telemetry.yaml
examples/mqtt/device-telemetry.yaml
```

`dry-run-device-telemetry.yaml` works without a running MQTT broker.
`device-telemetry.yaml` uses the real MQTT publish runtime and requires a broker
on the configured URL.

Request that works without MQTT:

```bash
curl -X POST http://localhost:7777/internal/mqtt/dry-run/devices/device-1/telemetry \
  -H "Content-Type: application/json" \
  -d '{"deviceId":"device-1"}'
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/mqtt/metrics
```

## WebSocket

Files:

```text
examples/websocket/chat.yaml
examples/websocket/dry-run-chat.yaml
```

`dry-run-chat.yaml` works without a WebSocket client. `chat.yaml` upgrades the
HTTP route to WebSocket and executes scripted receive/send/close steps.

Dry-run request:

```bash
curl http://localhost:7777/ws/dry-run/chat
```

Real WebSocket request:

```bash
printf '{"sender":"client","text":"ping"}\n' | websocat ws://localhost:7777/ws/chat
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/websocket/metrics
```

## S3

Files:

```text
examples/s3/object-put.yaml
examples/s3/object-get.yaml
examples/s3/object-delete.yaml
examples/s3/bucket-list.yaml
examples/s3/dry-run.yaml
```

These files provide a small S3-compatible path-style API backed by in-memory
storage.

Request flow:

```bash
curl -X PUT http://localhost:7777/s3/demo-bucket/readme.txt \
  -H "Content-Type: text/plain" \
  --data "hello from gigamock"

curl http://localhost:7777/s3/demo-bucket/readme.txt
curl http://localhost:7777/s3/demo-bucket
curl -X DELETE http://localhost:7777/s3/demo-bucket/readme.txt
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/s3/metrics
```

## SQS

Files:

```text
examples/sqs/send-message.yaml
examples/sqs/receive-message.yaml
examples/sqs/purge-queue.yaml
examples/sqs/dry-run.yaml
```

These files provide a small SQS-compatible path-style API backed by in-memory
queue storage.

Request flow:

```bash
curl -X POST http://localhost:7777/aws/sqs/orders \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl http://localhost:7777/aws/sqs/orders
curl -X DELETE http://localhost:7777/aws/sqs/orders
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/sqs/metrics
```

## SNS

Files:

```text
examples/sns/publish-message.yaml
examples/sns/list-messages.yaml
examples/sns/dry-run.yaml
```

These files provide a small SNS-compatible path-style API backed by in-memory
topic storage.

Request flow:

```bash
curl -X POST http://localhost:7777/aws/sns/order-events \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl http://localhost:7777/aws/sns/order-events
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/sns/metrics
```

## Google Pub/Sub

Files:

```text
examples/pubsub/publish-message.yaml
examples/pubsub/pull-message.yaml
examples/pubsub/purge-subscription.yaml
examples/pubsub/dry-run.yaml
```

These files provide a small Google Pub/Sub-compatible path-style API backed by
in-memory topic storage.

Request flow:

```bash
curl -X POST http://localhost:7777/gcp/pubsub/order-events/publish \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl -X POST http://localhost:7777/gcp/pubsub/orders-sub/pull
curl -X DELETE http://localhost:7777/gcp/pubsub/orders-sub/purge
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/pubsub/metrics
```

## Azure Service Bus

Files:

```text
examples/azure-servicebus/send-message.yaml
examples/azure-servicebus/receive-message.yaml
examples/azure-servicebus/purge-queue.yaml
examples/azure-servicebus/dry-run.yaml
```

These files provide a small Azure Service Bus-compatible path-style API backed
by in-memory queue storage.

Request flow:

```bash
curl -X POST http://localhost:7777/azure/servicebus/orders/send \
  -H "Content-Type: application/json" \
  --data '{"orderId":"order-1"}'

curl -X POST http://localhost:7777/azure/servicebus/orders/receive
curl -X DELETE http://localhost:7777/azure/servicebus/orders/purge
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/servicebus/metrics
```

## SOAP

File:

```text
examples/soap/customer-service.yaml
```

Use it to mock SOAP XML-over-HTTP services with `SOAPAction` and request body
matching.

Request:

```bash
curl -X POST http://localhost:7777/soap/customers \
  -H 'Content-Type: text/xml; charset=utf-8' \
  -H 'SOAPAction: "GetCustomer"' \
  --data '<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><GetCustomerRequest xmlns="urn:gigamock:customers"><customerId>customer-1</customerId></GetCustomerRequest></soap:Body></soap:Envelope>'
```

Metrics:

```bash
curl http://localhost:7777/internal/v1/soap/metrics
```
