# cronitor-local

Self-hosted cron job monitoring daemon with a simple web dashboard and alerting.

---

## Installation

```bash
go install github.com/yourusername/cronitor-local@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronitor-local.git
cd cronitor-local
go build -o cronitor-local .
```

---

## Usage

Start the daemon and web dashboard:

```bash
cronitor-local --port 8080 --config /etc/cronitor-local/config.yaml
```

In your crontab, wrap any job with the monitor command:

```cron
* * * * * cronitor-local run --job "backup" -- /usr/local/bin/backup.sh
```

The dashboard will be available at `http://localhost:8080`, showing job history, last run times, and alert status.

### Example config (`config.yaml`)

```yaml
alert:
  email: ops@example.com
  on_failure: true
  on_missed: true

jobs:
  - name: backup
    schedule: "0 2 * * *"
    timeout: 300
```

---

## Features

- Lightweight daemon with no external dependencies
- Web dashboard for real-time job monitoring
- Email and webhook alerting on failure or missed runs
- Simple YAML configuration

---

## License

MIT © 2024 Your Name