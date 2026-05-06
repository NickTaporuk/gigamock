# Start Gigamock

The main documentation index is available here:

- [Documentation index](README.md)
- [Scenario type overview](scenario-types/README.md)
- [Example YAML files](../examples/README.md)
- [Manual testing requests](../testsing/requests/README.md)

Run all examples:

```bash
go run ./cmd --dir-path ./examples
```

Run examples split by type:

```bash
go run ./cmd \
  --dir-path ./examples/rest \
  --dir-path ./examples/graphql \
  --dir-path ./examples/grpc \
  --dir-path ./examples/kafka \
  --dir-path ./examples/nats \
  --dir-path ./examples/rabbitmq \
  --dir-path ./examples/mqtt
```

Useful URLs:

```text
http://localhost:7777/internal/v1/mock-ui
http://localhost:7777/internal/v1/scenarios
```

Real gRPC endpoint:

```bash
grpcurl -plaintext localhost:7778 list
```

gRPC runtime metrics:

```bash
curl http://localhost:7777/internal/v1/grpc/metrics
```
