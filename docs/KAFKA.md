# Kafka Documentation

The canonical Kafka scenario documentation now lives here:

- [Kafka scenario fields](scenario-types/kafka.md)
- [Scenario type overview](scenario-types/README.md)
- [Kafka examples](../examples/kafka)
- [Manual broker requests](../testsing/requests/brokers.http)

Quick run:

```bash
go run ./cmd --dir-path ./examples/kafka
```

Quick request:

```bash
curl http://localhost:7777/internal/queue/message-1
```
