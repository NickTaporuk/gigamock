FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/gigamock ./cmd

FROM alpine:3.19

WORKDIR /app
COPY --from=build /out/gigamock /usr/local/bin/gigamock
COPY examples/websocket /app/examples/websocket

EXPOSE 7777
CMD ["gigamock", "--server-ip", "0.0.0.0", "--server-port", ":7777", "--dir-path", "/app/examples/websocket"]
