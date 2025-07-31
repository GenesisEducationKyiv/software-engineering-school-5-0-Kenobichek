# Weather-Forecast-API
Weather Subscription API – A simple API that lets users subscribe to weather updates for their city. Built for Genesis Software Engineering School 5.0.

---

## Quick Start

To start all microservices and Kafka, simply run:

```bash
./run_all.sh
```

This script will automatically start all required containers and display logs for all services.

---

### Monitoring Setup

The project includes Redis monitoring with Prometheus and Grafana.

#### Quick Start

```shell script
# Start monitoring stack
docker-compose up -d redis redis-exporter prometheus grafana

# Test the setup
./monitoring/test-redis-metrics.sh
```

For detailed monitoring documentation, see [ADR 003: Redis Monitoring Setup](docs/adr/003-redis-monitoring-setup.md).

---

### Running Tests

Prerequisites: Go ≥ 1.21 installed.

Command | What it runs
------- | ------------
`go test -v -short ./...` | Unit tests (fast, in-memory)
`go test -v -tags=integration ./tests/integration/...` | Integration tests (needs deps)
`go test -v -tags=e2e ./tests/e2e/...` | End-to-End tests
`go test -v ./...` | Everything

CI runs the same sequence: unit → integration → e2e.