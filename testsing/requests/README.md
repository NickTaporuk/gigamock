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
- `brokers.http`: Kafka plus broker HTTP-facing checks.
- `all-examples.http`: one file with a small request from each scenario type.

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
