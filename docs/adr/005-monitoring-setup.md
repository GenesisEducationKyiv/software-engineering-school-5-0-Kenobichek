# 005: Observability (VictoriaMetrics + Grafana)

## Overview

To ensure effective monitoring and operational insight, the Subscription Service will expose application metrics using Prometheus instrumentation. These metrics will be scraped and stored by VictoriaMetrics and visualized in Grafana dashboards. This setup facilitates tracking of subscription activity and error occurrences in near real-time, enabling data-driven operational decisions.

## Decision

1. Integrate the [prometheus/client_golang] library to expose metrics from the Go service.
2. Expose metrics at the `/metrics` HTTP endpoint, listening on the same port that VictoriaMetrics uses for its HTTP API — configured by the `VICTORIA_METRICS_PORT` environment variable (default: **8428**). This removes the need for a separate `METRICS_PORT`. 
3. Run the metrics HTTP server inside the internal/app layer instead of main to keep unit tests focused and lightweight.
4. Utilize Docker Compose to orchestrate the following services:
   - subscription-service — the application itself.
   - victoria-metrics — scrapes and stores metrics from subscription-service.
   - grafana — visualizes Prometheus data.
5. Maintain a single docker-compose.yaml.
6. Control all ports and scrape intervals through environment variables, validated within the config package.

## How it works

1. At startup `subscription-service` registers counters/gauges in `internal/observability/metrics` and launches an HTTP server on `/metrics` (port `VICTORIA_METRICS_PORT`).
2. VictoriaMetrics scrapes that endpoint every **15 seconds** (see `prometheus.yml`). The scrape config is mounted into the container, and unknown fields are ignored via `--promscrape.config.strictParse=false`. 
3. Grafana connects to VictoriaMetrics and displays the ready-made dashboard from `grafana-dashboard.json`.

## Access URLs

* VictoriaMetrics UI: `http://localhost:$VICTORIA_PORT`
* Grafana: `http://localhost:$GRAFANA_PORT` (login/password: `admin/admin`)
* Service metrics: `http://localhost:$VICTORIA_METRICS_PORT/metrics`

## Adding a new metric
1. Declare a variable in `internal/observability/metrics/metrics.go`.
2. Register it in `Register()`.
3. Increment/decrement it in the relevant part of the code.
4. Rebuild `subscription-service`.
5. Add a panel for it in Grafana.