# driftwatch

> Monitors config file changes and alerts via webhook

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git && cd driftwatch && go build -o driftwatch .
```

---

## Usage

Create a `driftwatch.yaml` configuration file:

```yaml
watch:
  - path: /etc/myapp/config.yaml
  - path: /etc/myapp/secrets.env

webhook:
  url: https://hooks.slack.com/services/your/webhook/url
  method: POST
```

Then run:

```bash
driftwatch --config driftwatch.yaml
```

When a monitored file changes, driftwatch sends a POST request to your webhook URL with a JSON payload containing the file path, change type, and timestamp.

```json
{
  "file": "/etc/myapp/config.yaml",
  "event": "modified",
  "timestamp": "2024-05-10T14:32:00Z"
}
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `driftwatch.yaml` | Path to config file |
| `--interval` | `5s` | Polling interval |
| `--verbose` | `false` | Enable verbose logging |

---

## License

MIT © [yourusername](https://github.com/yourusername)