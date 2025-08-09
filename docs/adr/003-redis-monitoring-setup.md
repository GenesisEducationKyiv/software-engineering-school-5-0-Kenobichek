# 003: Redis Monitoring Setup with Prometheus and Grafana

## Context

The Weather Forecast API uses Redis for caching weather data to improve performance and reduce external API calls. We need to monitor Redis performance and health to ensure optimal cache performance and identify potential issues before they impact the application.

## Decision

We will implement a comprehensive Redis monitoring solution using:

1. **Redis Exporter** - Converts Redis metrics to Prometheus format
2. **Prometheus** - Time-series database for metrics storage
3. **Grafana** - Visualization and alerting platform
4. **Environment Variables** - Configuration management for flexibility

## Consequences

### Positive

- **Comprehensive Monitoring**: Full visibility into Redis performance, memory usage, and cache hit rates
- **Professional Dashboards**: Beautiful, pre-built Grafana dashboards with real-time metrics
- **Flexible Configuration**: Environment variables allow easy customization for different environments
- **Industry Standard**: Uses proven, battle-tested monitoring tools
- **No Code Changes**: Redis Exporter works out of the box without modifying application code
- **Alerting Capabilities**: Built-in alerting for proactive issue detection
- **Scalable**: Can easily add more services to the monitoring stack

### Negative

- **Additional Complexity**: Introduces 3 new containers to manage
- **Resource Usage**: Slight increase in system resources
- **Learning Curve**: Team needs to understand Prometheus and Grafana
- **Port Management**: Requires careful port configuration to avoid conflicts

## Implementation Details

### Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐    ┌─────────────┐
│   Redis     │    │ Redis        │    │ Prometheus  │    │   Grafana   │
│   Server    │───▶│ Exporter     │───▶│             │───▶│             │
│             │    │ (Port 9121)  │    │ (Port 9090) │    │ (Port 3000) │
└─────────────┘    └──────────────┘    └─────────────┘    └─────────────┘
```

### Key Metrics Monitored

- **Performance**: Commands per second, response times
- **Cache Performance**: Hit/miss rates, cache efficiency
- **Resource Usage**: Memory usage, connected clients
- **Network**: I/O throughput

### Configuration Management

All configuration is externalized using environment variables:

```bash
# Monitoring Configuration
REDIS_EXPORTER_PORT=9121
PROMETHEUS_PORT=9090
GRAFANA_PORT=3000
GRAFANA_ADMIN_PASSWORD=admin
PROMETHEUS_RETENTION_TIME=200h

# Health Check Configuration
HEALTH_CHECK_INTERVAL=10s
HEALTH_CHECK_RETRIES=5
HEALTH_CHECK_START_PERIOD=30s
HEALTH_CHECK_TIMEOUT=10s

# Container Configuration
CONTAINER_RESTART_POLICY=unless-stopped
```

### Files Added

- `docker-compose.yml` - Updated with monitoring services
- `monitoring/prometheus.yml` - Prometheus configuration
- `monitoring/grafana/dashboards/redis-dashboard.json` - Redis dashboard
- `monitoring/grafana/datasources/prometheus.yml` - Grafana data source
- `monitoring/grafana/dashboards/dashboards.yml` - Dashboard provisioning
- `monitoring/test-redis-metrics.sh` - Test script
- `env.example` - Environment variables template
- `setup-env.sh` - Environment setup script

## Usage

### Quick Start
```bash
# Setup environment
./setup-env.sh

# Start monitoring stack
docker-compose up -d redis redis-exporter prometheus grafana

# Test setup
./monitoring/test-redis-metrics.sh
```

### Access URLs
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)
- Redis Exporter: http://localhost:9121/metrics

## Future Considerations

- **Alerting Rules**: Add Prometheus alerting rules for critical metrics
- **Custom Metrics**: Consider adding application-specific Redis metrics
- **Authentication**: Add authentication to Prometheus and Grafana for production
- **Backup**: Implement Prometheus data backup strategy
- **Scaling**: Consider Prometheus federation for multiple environments
