# SOAP Scenario Fields

SOAP scenarios provide XML-over-HTTP mocks for local and CI contract tests.
Request matching can use `SOAPAction`, headers, and a body substring.

Supported route shape:

```text
POST /soap/customers
```

The route is fully configurable through the top-level `path` and `method`
fields, like HTTP and GraphQL mocks.

Example:

```yaml
path: "/soap/customers"
method: POST
type: soap
description: "SOAP customer service mock"
scenarios:
  - name: "get customer"
    request:
      soapAction: "GetCustomer"
      bodyContains: "<customerId>customer-1</customerId>"
    response:
      statusCode: 200
      headers:
        Content-Type: "text/xml; charset=utf-8"
      body: |
        <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
          <soap:Body>
            <GetCustomerResponse xmlns="urn:gigamock:customers">
              <customer>
                <id>customer-1</id>
              </customer>
            </GetCustomerResponse>
          </soap:Body>
        </soap:Envelope>
```

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `request.headers` | no | Exact HTTP header matches. |
| `request.soapAction` | no | Expected `SOAPAction` header value. Quotes in the incoming header are ignored. |
| `request.bodyContains` | no | XML substring that must be present in the request body. |
| `response.statusCode` | yes | HTTP status code. |
| `response.headers` | no | Response headers. Defaults to `Content-Type: text/xml; charset=utf-8`. |
| `response.body` | no | SOAP envelope or fault XML. |

Smoke test:

```bash
curl -X POST http://localhost:7777/soap/customers \
  -H 'Content-Type: text/xml; charset=utf-8' \
  -H 'SOAPAction: "GetCustomer"' \
  --data '<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><GetCustomerRequest xmlns="urn:gigamock:customers"><customerId>customer-1</customerId></GetCustomerRequest></soap:Body></soap:Envelope>'

curl http://localhost:7777/internal/v1/soap/metrics
```

Example files:

```text
examples/soap/customer-service.yaml
```
