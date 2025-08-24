#!/bin/bash

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set defaults if not defined
REDIS_EXPORTER_PORT=${REDIS_EXPORTER_PORT:-9121}
PROMETHEUS_PORT=${PROMETHEUS_PORT:-9090}
GRAFANA_PORT=${GRAFANA_PORT:-3000}

echo "Testing Redis Metrics Collection..."
echo "=================================="

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 30

# Test Redis Exporter metrics endpoint
echo "Testing Redis Exporter metrics endpoint..."
curl -s http://localhost:${REDIS_EXPORTER_PORT}/metrics | head -20

echo ""
echo "Testing Prometheus targets..."
curl -s http://localhost:${PROMETHEUS_PORT}/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health, lastError: .lastError}'

echo ""
echo "Testing Prometheus Redis metrics..."
curl -s "http://localhost:${PROMETHEUS_PORT}/api/v1/query?query=redis_connected_clients" | jq '.'

echo ""
echo "Access URLs:"
echo "- Prometheus: http://localhost:${PROMETHEUS_PORT}"
echo "- Grafana: http://localhost:${GRAFANA_PORT} (admin/${GRAFANA_ADMIN_PASSWORD:-admin})"
echo "- Redis Exporter: http://localhost:${REDIS_EXPORTER_PORT}/metrics"

echo ""
echo "To test Redis operations, run:"
echo "docker exec -it weather-redis redis-cli"
echo "Then try: SET test:key 'hello world'"
echo "And: GET test:key" 