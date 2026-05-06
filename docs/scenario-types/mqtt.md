# MQTT Scenario Fields

MQTT scenarios define the planned MQTT publish/subscribe mock contract.
Gigamock currently indexes these files and displays them in the control UI.
Native MQTT runtime support is planned.

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
| `broker` | planned | MQTT broker URL. |
| `clientID` | planned | MQTT client identifier. |
| `topic` | planned | MQTT topic. |
| `qos` | planned | MQTT quality of service level. |
| `retained` | planned | Whether the message should be retained. |
| `message.contentType` | planned | Message content type. |
| `message.body` | planned | Message body. |

Example file:

```text
examples/mqtt/device-telemetry.yaml
```
