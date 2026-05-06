# gRPC Scenario Fields

gRPC scenarios define the planned production-ready gRPC mock contract.
Gigamock currently indexes these files and displays them in the control UI.
Native gRPC runtime serving is planned.

Unary example:

```yaml
path: "/customers.CustomersService/GetCustomer"
method: POST
type: grpc
description: "Unary gRPC mock for retrieving a customer by id"
proto:
  descriptorSet: "./examples/grpc/proto/customers.pb"
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
| `proto` | planned | Protobuf descriptor configuration. |
| `scenarios` | yes | Ordered list of gRPC scenarios. |

`proto` fields:

| Field | Required | Description |
| --- | --- | --- |
| `descriptorSet` | planned | Path to compiled protobuf descriptor set. |
| `file` | planned | Path to `.proto` file if direct proto parsing is supported. |
| `importPaths` | planned | Import paths for direct `.proto` parsing. |
| `service` | planned | Fully qualified service name. |
| `method` | planned | RPC method name. |

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
