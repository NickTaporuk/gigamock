# MQTT Scenario Fields

MQTT scenarios publish configured messages to an MQTT topic. For local route/UI
testing without an MQTT broker, use `dryRun: true`; the mock will validate the
scenario, skip network calls, and return a successful JSON response.

Example:

```yaml
path: "/internal/mqtt/devices/:deviceID/telemetry"
method: POST
type: mqtt
description: "MQTT mock scenario for publishing device telemetry"
scenarios:
  - name: "publish temperature telemetry"
    broker: "tcp://localhost:1883"
    clientID: "gigamock-device-telemetry"
    topic: "devices/device-1/telemetry"
    qos: 1
    retained: false
    dryRun: false
    message:
      contentType: "application/json"
      body: |
        {
          "deviceId": "device-1",
          "temperature": 21.7,
          "humidity": 42.5
        }
```

Top-level fields:

| Field | Required | Description |
| --- | --- | --- |
| `path` | yes | HTTP route used to index and switch the scenario. |
| `method` | yes | HTTP method used to index and switch the scenario. |
| `type` | yes | Must be `mqtt`. |
| `description` | no | Text shown in the control UI. |
| `scenarios` | yes | Ordered list of MQTT scenarios. |

Scenario fields:

| Field | Required | Description |
| --- | --- | --- |
| `name` | no | Human-readable scenario name shown in the UI. |
| `broker` | yes | MQTT broker URL, for example `tcp://localhost:1883`. |
| `clientID` | yes | MQTT client identifier. |
| `topic` | yes | MQTT topic. |
| `qos` | no | MQTT quality of service level. |
| `retained` | no | Whether the message should be retained. |
| `dryRun` | no | When `true`, skips MQTT broker calls and returns a successful response. |
| `message.contentType` | no | Message content type. |
| `message.body` | yes | Message body. |

Runtime responses:

Successful publish response:

```json
{
  "topic": "devices/device-1/telemetry",
  "published": true,
  "dryRun": false
}
```

Dry-run publish response:

```json
{
  "topic": "devices/device-1/telemetry/dry-run",
  "published": true,
  "dryRun": true
}
```

Runtime metrics:

```bash
curl http://localhost:7777/internal/v1/mqtt/metrics
```

Example files:

```text
examples/mqtt/dry-run-device-telemetry.yaml
examples/mqtt/device-telemetry.yaml
```
