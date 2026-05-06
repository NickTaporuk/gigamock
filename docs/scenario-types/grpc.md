# gRPC Scenario Fields

gRPC scenarios define production-ready dynamic gRPC mocks. Gigamock loads
`.proto` files, registers services on a real `grpc.Server`, exposes reflection
for `grpcurl`, and uses YAML scenarios to build protobuf responses.

Unary example:

```yaml
path: "/customers.CustomersService/GetCustomer"
method: POST
type: grpc
description: "Unary gRPC mock for retrieving a customer by id"
proto:
  file: "customers.proto"
  importPaths:
    - "./examples/grpc/proto"
  service: "customers.CustomersService"
  method: "GetCustomer"
scenarios:
  - name: "default customer"
    request:
      match:
        customerId: "customer-1"
    response:
      code: OK
      metadata:
        x-mock-source: "gigamock"
      body:
        customer:
          id: "customer-1"
          name: "John Doe"
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | Full gRPC method path, for example `/package.Service/Method`. |
| `method` | yes | HTTP-facing index method. Use `POST`. |
| `type` | yes | Must be `grpc`. |
| `description` | no | Text shown in the control UI. |
| `proto` | yes | Protobuf source configuration. |
| `scenarios` | yes | Ordered list of gRPC scenarios. |

`proto` fields:

| Field | Required | Description |
| --- | --- | --- |
| `file` | yes | `.proto` file to compile. |
| `importPaths` | yes | Import paths used to resolve `file` and its imports. |
| `service` | yes | Fully qualified service name. |
| `method` | yes | RPC method name. |
| `descriptorSet` | no | Reserved for future descriptor-set loading. |

Unary scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `request.match` | no | Expected request fields. |
| `response.code` | yes | gRPC status code such as `OK` or `NOT_FOUND`. |
| `response.message` | no | gRPC status message for error responses. |
| `response.metadata` | no | Response metadata. |
| `response.body` | no | Response message encoded as YAML/JSON-like data. |

Streaming scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `stream.sendOnConnect` | no | Messages sent as soon as the stream opens. |
| `stream.steps` | no | Ordered scripted receive/send/close flow. |
| `stream.onReceive` | no | Rule-based receive/send flow. |
| `receive` | no | Expected client stream message. |
| `send` | no | Server stream message. |
| `close.code` | no | Status code used to close the stream. |
| `close.message` | no | Status message used to close the stream. |

Example files:

```text
examples/grpc/customer-service-unary.yaml
examples/grpc/billing-service-unary.yaml
examples/grpc/billing-service-server-stream.yaml
examples/grpc/billing-service-client-stream.yaml
examples/grpc/billing-service-bidi-stream.yaml
examples/grpc/chat-service-bidi-stream.yaml
```

Real gRPC request example:

```bash
go run ./cmd --dir-path ./examples/grpc

grpcurl -plaintext \
  -d '{"customerId":"customer-1"}' \
  localhost:7778 \
  customers.CustomersService/GetCustomer
```

Billing examples cover unary, server-streaming, and bidirectional streaming:

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

Switch a gRPC scenario at runtime:

```bash
curl -X POST http://localhost:7777/internal/v1/in-memory \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/customers.CustomersService/GetCustomer",
    "method": "POST",
    "scenarioNumber": 2
  }'
```

Reflection lists every service loaded from every configured gRPC directory:

```bash
grpcurl -plaintext localhost:7778 list
grpcurl -plaintext localhost:7778 list billing.BillingService
```

## Hardening Settings

gRPC runtime supports extra flags for safer local and CI usage:

```bash
go run ./cmd \
  --dir-path ./examples/grpc \
  --grpc-stream-max-messages 100 \
  --grpc-stream-timeout-seconds 300
```

TLS can be enabled with:

```bash
go run ./cmd \
  --dir-path ./examples/grpc \
  --grpc-tls-cert-file ./certs/server.crt \
  --grpc-tls-key-file ./certs/server.key
```

mTLS can be enabled by also passing a client CA:

```bash
go run ./cmd \
  --dir-path ./examples/grpc \
  --grpc-tls-cert-file ./certs/server.crt \
  --grpc-tls-key-file ./certs/server.key \
  --grpc-tls-client-ca-file ./certs/ca.crt
```

Runtime metrics are exposed through the HTTP control plane:

```bash
curl http://localhost:7777/internal/v1/grpc/metrics
```

The response is keyed by full gRPC method name and contains call/error counts.

## Descriptor Sets

gRPC scenarios can load either `.proto` source files or descriptor sets. Use
exactly one of `proto.file` or `proto.descriptorSet`.

```yaml
proto:
  descriptorSet: "./examples/grpc/proto/billing.pb"
  service: "billing.BillingService"
  method: "GetInvoice"
```

Descriptor-set support is useful in CI or container images where generated
descriptor artifacts are easier to ship than source proto trees.
