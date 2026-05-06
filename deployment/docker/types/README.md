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
```

Run example:

```bash
docker run --rm -p 7777:7777 gigamock-http
```

Open the control UI:

```text
http://localhost:7777/internal/v1/mock-ui
```
