# Gigamock Testing Requests

These files are ready-to-run manual requests for JetBrains HTTP Client,
VS Code REST Client, and `grpcurl`.

Start all examples:

```bash
go run ./cmd --dir-path ./examples
```

Useful files:

- `in-memory.http`: control API, scenario listing, scenario switching.
- `rest.http`: HTTP REST examples.
- `graphql.http`: GraphQL examples.
- `grpc.grpc`: IDE gRPC client examples.
- `grpcurl.sh`: CLI gRPC smoke test with scenario switching.
- `websocket.http`: WebSocket dry-run and CLI smoke commands.
- `s3.http`: S3-compatible object API smoke requests.
- `aws.http`: SQS/SNS-compatible queue and topic smoke requests.
- `pubsub.http`: Google Pub/Sub-compatible topic and subscription smoke requests.
- `azure-servicebus.http`: Azure Service Bus-compatible queue smoke requests.
- `soap.http`: SOAP XML-over-HTTP smoke requests.
- `brokers.http`: Kafka/NATS/RabbitMQ/MQTT dry-run, real broker routes, and broker
  HTTP-facing checks.
- `all-examples.http`: one file with a small request from each scenario type.

Kafka dry-run check without a broker:

```bash
curl http://localhost:7777/internal/kafka/dry-run/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
```

Kafka end-to-end check with Docker:

```bash
task docker:kafka:up
curl http://localhost:7777/internal/kafka/docker/message-1
curl http://localhost:7777/internal/v1/kafka/metrics
task docker:kafka:down
```

Use `PORT=7781 task docker:kafka:up` when `7777` is already occupied.

NATS dry-run check without a broker:

```bash
curl -X POST http://localhost:7777/internal/nats/dry-run/orders/order-1 \
  -H "Content-Type: application/json" \
  -d '{"orderId":"order-1"}'
curl http://localhost:7777/internal/v1/nats/metrics
```

RabbitMQ dry-run check without a broker:

```bash
curl -X POST http://localhost:7777/internal/rabbitmq/dry-run/payments/payment-1 \
  -H "Content-Type: application/json" \
  -d '{"paymentId":"payment-1"}'
curl http://localhost:7777/internal/v1/rabbitmq/metrics
```

MQTT dry-run check without a broker:

```bash
curl -X POST http://localhost:7777/internal/mqtt/dry-run/devices/device-1/telemetry \
  -H "Content-Type: application/json" \
  -d '{"deviceId":"device-1"}'
curl http://localhost:7777/internal/v1/mqtt/metrics
```

WebSocket dry-run check:

```bash
curl http://localhost:7777/ws/dry-run/chat
curl http://localhost:7777/internal/v1/websocket/metrics
```

Real WebSocket check:

```bash
printf '{"sender":"client","text":"ping"}\n' | websocat ws://localhost:7777/ws/chat
```

S3-compatible check:

```bash
curl -X PUT http://localhost:7777/s3/demo-bucket/readme.txt \
  -H "Content-Type: text/plain" \
  --data "hello from gigamock"
curl http://localhost:7777/s3/demo-bucket/readme.txt
curl http://localhost:7777/internal/v1/s3/metrics
```

For real gRPC checks:

```bash
bash testsing/requests/grpcurl.sh
```

Billing gRPC example:

```bash
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  localhost:7778 \
  billing.BillingService/GetInvoice
```

Billing server-streaming example:

```bash
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  localhost:7778 \
  billing.BillingService/WatchInvoice
```

Billing bidirectional example:

```bash
grpcurl -plaintext \
  -d '{"text":"hello","sender":"client"}' \
  localhost:7778 \
  billing.BillingService/BillingChat
```

Billing client-streaming example:

```bash
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1","status":"OPEN","message":"invoice created"}{"invoiceId":"invoice-1","status":"PAID","message":"invoice paid"}' \
  localhost:7778 \
  billing.BillingService/UploadInvoiceEvents
```
