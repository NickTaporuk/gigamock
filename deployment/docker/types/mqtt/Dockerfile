FROM golang:1.20-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/gigamock ./cmd

FROM alpine:3.19

WORKDIR /app
COPY --from=build /out/gigamock /usr/local/bin/gigamock
COPY examples/mqtt /app/examples/mqtt

EXPOSE 7777
CMD ["gigamock", "--server-ip", "0.0.0.0", "--server-port", ":7777", "--dir-path", "/app/examples/mqtt"]
