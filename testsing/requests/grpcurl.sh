#!/usr/bin/env bash
set -euo pipefail

GRPC_HOST="${GRPC_HOST:-localhost:7778}"
HTTP_HOST="${HTTP_HOST:-http://localhost:7777}"

echo "== gRPC services =="
grpcurl -plaintext "${GRPC_HOST}" list

echo
echo "== Unary default customer =="
grpcurl -plaintext \
  -d '{"customerId":"customer-1"}' \
  "${GRPC_HOST}" \
  customers.CustomersService/GetCustomer

echo
echo "== Unary request-matched disabled customer =="
grpcurl -plaintext \
  -d '{"customerId":"customer-disabled"}' \
  "${GRPC_HOST}" \
  customers.CustomersService/GetCustomer

echo
echo "== Billing paid invoice =="
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  "${GRPC_HOST}" \
  billing.BillingService/GetInvoice

echo
echo "== Billing request-matched open invoice =="
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-open"}' \
  "${GRPC_HOST}" \
  billing.BillingService/GetInvoice

echo
echo "== Billing server-streaming invoice lifecycle =="
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1"}' \
  "${GRPC_HOST}" \
  billing.BillingService/WatchInvoice

echo
echo "== Billing client-streaming upload events =="
grpcurl -plaintext \
  -d '{"invoiceId":"invoice-1","status":"OPEN","message":"invoice created"}{"invoiceId":"invoice-1","status":"PAID","message":"invoice paid"}' \
  "${GRPC_HOST}" \
  billing.BillingService/UploadInvoiceEvents

echo
echo "== Billing bidirectional stream happy path =="
grpcurl -plaintext \
  -d '{"text":"hello","sender":"client"}' \
  "${GRPC_HOST}" \
  billing.BillingService/BillingChat

echo
echo "== Switch billing bidirectional stream to rule-based scenario =="
curl -s -X POST "${HTTP_HOST}/internal/v1/in-memory" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/billing.BillingService/BillingChat",
    "method": "POST",
    "scenarioNumber": 2
  }'

echo
echo "== Billing bidirectional stream rule-based response =="
grpcurl -plaintext \
  -d '{"text":"status invoice-1","sender":"client"}' \
  "${GRPC_HOST}" \
  billing.BillingService/BillingChat

echo
echo "== Switch unary to NOT_FOUND scenario =="
curl -s -X POST "${HTTP_HOST}/internal/v1/in-memory" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/customers.CustomersService/GetCustomer",
    "method": "POST",
    "scenarioNumber": 2
  }'

echo
echo "== Unary active NOT_FOUND scenario =="
grpcurl -plaintext \
  -d '{"customerId":"anything"}' \
  "${GRPC_HOST}" \
  customers.CustomersService/GetCustomer || true

echo
echo "== Switch stream to onReceive scenario =="
curl -s -X POST "${HTTP_HOST}/internal/v1/in-memory" \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/chat.ChatService/Chat",
    "method": "POST",
    "scenarioNumber": 2
  }'

echo
echo "== Bidirectional stream onReceive =="
grpcurl -plaintext \
  -d '{"text":"status","sender":"client"}' \
  "${GRPC_HOST}" \
  chat.ChatService/Chat
