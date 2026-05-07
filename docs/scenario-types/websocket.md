# WebSocket Scenario Fields

WebSocket scenarios define scripted bidirectional communication over an HTTP
`GET` route. Gigamock upgrades the request to WebSocket, sends optional
messages on connect, then executes ordered receive/send/close steps.

For local route/UI testing without a WebSocket client, use `dryRun: true`; the
mock will validate the script, skip the protocol upgrade, and return a JSON
response.

Example:

```yaml
path: "/ws/chat"
method: GET
type: websocket
description: "Scripted WebSocket chat mock"
scenarios:
  - name: "ping pong chat"
    sendOnConnect:
      - type: "text"
        text: "{\"sender\":\"mock-server\",\"text\":\"connected\"}"
    steps:
      - receive:
          type: "text"
          text: "{\"sender\":\"client\",\"text\":\"ping\"}"
      - send:
          type: "text"
          text: "{\"sender\":\"mock-server\",\"text\":\"pong\"}"
      - close:
          code: 1000
          reason: "scenario completed"
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route upgraded to WebSocket. |
| `method` | yes | Must usually be `GET` for WebSocket upgrade requests. |
| `type` | yes | Must be `websocket`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of WebSocket scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `dryRun` | no | When `true`, skips WebSocket upgrade and returns a JSON response. |
| `sendOnConnect` | no | Messages sent immediately after the WebSocket connection opens. |
| `steps` | no | Ordered script steps. Required unless `dryRun` or `sendOnConnect` is present. |

Message fields:

| Field | Required | Description |
| --- | --- | --- |
| `type` | no | `text` or `binary`. Defaults to `text`. |
| `text` | yes | Message payload. |

Step fields:

| Field | Required | Description |
| --- | --- | --- |
| `receive` | no | Expected message from the client. |
| `send` | no | Message sent by Gigamock. |
| `close` | no | Close frame sent by Gigamock. |
| `delay` | no | Go duration string, for example `500ms` or `2s`. |

Close fields:

| Field | Required | Description |
| --- | --- | --- |
| `code` | no | WebSocket close code. Defaults to `1000` when omitted. |
| `reason` | no | Close reason text. |

Dry-run check:

```bash
curl http://localhost:7777/ws/dry-run/chat
```

Real WebSocket check with `websocat`:

```bash
printf '{"sender":"client","text":"ping"}\n' | websocat ws://localhost:7777/ws/chat
```

Runtime metrics:

```bash
curl http://localhost:7777/internal/v1/websocket/metrics
```

Example files:

```text
examples/websocket/chat.yaml
examples/websocket/dry-run-chat.yaml
```
