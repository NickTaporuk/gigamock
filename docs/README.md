# Gigamock Documentation

Start here when you need to understand how to describe mocks, run the server,
or test scenario files.

## Scenario Types

Each scenario type has its own reference page:

| Type | Documentation | Example Directory | Runtime Status |
| --- | --- | --- | --- |
| HTTP | [HTTP scenarios](scenario-types/http.md) | [`examples/rest`](../examples/rest) | Runtime implemented |
| GraphQL | [GraphQL scenarios](scenario-types/graphql.md) | [`examples/graphql`](../examples/graphql) | Runtime implemented |
| gRPC | [gRPC scenarios](scenario-types/grpc.md) | [`examples/grpc`](../examples/grpc) | Runtime implemented |
| Kafka | [Kafka scenarios](scenario-types/kafka.md) | [`examples/kafka`](../examples/kafka) | Runtime implemented |
| NATS | [NATS scenarios](scenario-types/nats.md) | [`examples/nats`](../examples/nats) | Runtime implemented |
| RabbitMQ | [RabbitMQ scenarios](scenario-types/rabbitmq.md) | [`examples/rabbitmq`](../examples/rabbitmq) | Runtime implemented |
| MQTT | [MQTT scenarios](scenario-types/mqtt.md) | [`examples/mqtt`](../examples/mqtt) | Runtime implemented |
| WebSocket | [WebSocket scenarios](scenario-types/websocket.md) | [`examples/websocket`](../examples/websocket) | Runtime implemented |

## Common References

- [Scenario type overview](scenario-types/README.md)
- [Example YAML files](../examples/README.md)
- [Manual testing requests](../testsing/requests/README.md)
- [Feature specs](../features/README.md)
- [Per-type Dockerfiles](../deployment/docker/types/README.md)

## Runtime Control

The HTTP control plane is available when Gigamock is running:

```text
http://localhost:7777/internal/v1/mock-ui
```

Useful control API endpoints:

```bash
curl http://localhost:7777/internal/v1/scenarios
curl http://localhost:7777/internal/v1/in-memory
```
