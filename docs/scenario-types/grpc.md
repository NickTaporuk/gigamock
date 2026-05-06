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
