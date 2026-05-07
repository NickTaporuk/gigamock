# Per-Type Dockerfiles

This directory contains one Dockerfile per scenario type.

Build examples from the repository root:

```bash
docker build -f deployment/docker/types/http/Dockerfile -t gigamock-http .
docker build -f deployment/docker/types/graphql/Dockerfile -t gigamock-graphql .
docker build -f deployment/docker/types/grpc/Dockerfile -t gigamock-grpc .
docker build -f deployment/docker/types/kafka/Dockerfile -t gigamock-kafka .
docker build -f deployment/docker/types/nats/Dockerfile -t gigamock-nats .
docker build -f deployment/docker/types/rabbitmq/Dockerfile -t gigamock-rabbitmq .
docker build -f deployment/docker/types/mqtt/Dockerfile -t gigamock-mqtt .
docker build -f deployment/docker/types/websocket/Dockerfile -t gigamock-websocket .
docker build -f deployment/docker/types/s3/Dockerfile -t gigamock-s3 .
docker build -f deployment/docker/types/sqs/Dockerfile -t gigamock-sqs .
docker build -f deployment/docker/types/sns/Dockerfile -t gigamock-sns .
docker build -f deployment/docker/types/pubsub/Dockerfile -t gigamock-pubsub .
docker build -f deployment/docker/types/azure-servicebus/Dockerfile -t gigamock-servicebus .
docker build -f deployment/docker/types/soap/Dockerfile -t gigamock-soap .
```

Run example:

```bash
docker run --rm -p 7777:7777 gigamock-http
```

Open the control UI:

```text
http://localhost:7777/internal/v1/mock-ui
```
