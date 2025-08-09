# Weather Subscription Platform

Full-featured weather subscription platform composed of four Go microservices with Kafka as the message bus.

## 🚀 Quick Start

### Requirements

- Docker & Docker Compose
- Go 1.21+ (for development)

### Start all services with one script

1. Create a `.env` file based on `.env.example`
2. Run the services:

```bash
./run_all.sh
```

After the script finishes, services are available at:
- API-Gateway: http://localhost:8080
- weather-service: http://localhost:8081
- subscription-service: http://localhost:8082
- notification-service: http://localhost:8083
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Database (PostgreSQL): localhost:5433

## 🏗️ Architecture

- **api-gateway** – single REST entry point
- **weather-service** – fetches & caches weather data
- **subscription-service** – manages subscriptions (PostgreSQL)
- **notification-service** – sends email notifications

## 📜 Helper Scripts

| Script | Purpose |
|--------|---------|
| `run_all.sh` | Start **all** services (Kafka + 4 microservices) via multiple Docker-Compose files |
| `lint_all.sh` | Run Go linters across the entire workspace |
| `test_all.sh` | Execute unit tests for every service |

---

## 📝 License

This project is part of the Genesis Software Engineering School 5.0 curriculum.