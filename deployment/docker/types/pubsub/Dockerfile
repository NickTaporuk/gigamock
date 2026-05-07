FROM golang:1.20-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/gigamock ./cmd

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /out/gigamock /usr/local/bin/gigamock
COPY examples/pubsub /app/examples/pubsub
EXPOSE 7777

CMD ["gigamock", "--server-ip", "0.0.0.0", "--server-port", ":7777", "--dir-path", "/app/examples/pubsub"]
