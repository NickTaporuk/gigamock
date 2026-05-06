# Gherkin Feature Specs

This directory contains Gherkin `.feature` files for acceptance testing
Gigamock scenario types.

The files are intentionally runner-neutral. They can be wired to `godog`,
Cucumber, or another Gherkin runner.

Feature files:

- `http.feature`
- `graphql.feature`
- `grpc.feature`
- `kafka.feature`
- `nats.feature`
- `rabbitmq.feature`
- `mqtt.feature`

Suggested local command shape once a runner is added:

```bash
godog features
```

The mock server can be started with all examples:

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
