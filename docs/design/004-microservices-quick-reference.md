# Design-004: Microservices Quick Reference

## Summary

The Weather Forecast API should be decomposed into **5 microservices** with clear bounded contexts and communication patterns.

## Microservices Overview

| Service | Port | Database | Primary Responsibility |
|---------|------|----------|----------------------|
| **API Gateway** | 8080 | None | Request routing & aggregation |
| **Weather Service** | 8081 | Redis | Weather data aggregation |
| **Subscription Service** | 8082 | PostgreSQL | Subscription management |
| **Notification Service** | 8083 | PostgreSQL | Multi-channel notifications |
| **Scheduler Service** | 8084 | PostgreSQL | Background job orchestration |

## Communication Patterns

### Synchronous Communication
- **REST API**: Client-facing endpoints, external API calls
- **gRPC**: High-performance internal service communication

### Asynchronous Communication
- **Message Queue**: Event-driven communication, reliability

### Communication Matrix

| From \ To | Weather | Subscription | Notification | Scheduler | API Gateway |
|-----------|---------|--------------|--------------|-----------|-------------|
| **Weather** | - | gRPC | - | gRPC | REST |
| **Subscription** | - | - | Message Queue | gRPC | REST |
| **Notification** | - | - | - | Message Queue | REST |
| **Scheduler** | gRPC | gRPC | Message Queue | - | - |
| **API Gateway** | REST | REST | REST | - | - |

## Key Benefits

### ✅ Scalability
- Independent scaling per service
- Resource optimization
- Horizontal scaling

### ✅ Maintainability
- Focused codebases
- Independent deployments
- Technology diversity

### ✅ Resilience
- Fault isolation
- Circuit breakers
- Graceful degradation

### ✅ Team Organization
- Service ownership
- Parallel development
- Reduced coordination

## Migration Phases

### Phase 1: Weather Service
- Extract weather provider logic
- Implement caching layer
- Set up service communication

### Phase 2: Subscription Service
- Extract subscription management
- Implement message queue
- Migrate data

### Phase 3: Notification Service
- Extract notification logic
- Set up template management
- Configure message consumers

### Phase 4: Scheduler Service
- Extract scheduling logic
- Implement orchestration
- Complete migration

## Infrastructure Requirements

### Shared Services
- **PostgreSQL**: Primary data storage
- **Redis**: Caching and session data
- **Message Queue**: RabbitMQ or Apache Kafka
- **Monitoring**: Prometheus + Grafana + Jaeger

### Security
- **mTLS**: Service-to-service communication
- **JWT**: API Gateway authentication
- **Secrets Management**: Kubernetes secrets or Vault

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| **Complexity** | Comprehensive monitoring & automation |
| **Network Latency** | gRPC, caching, connection pooling |
| **Data Consistency** | Event-driven architecture, eventual consistency |
| **Testing Complexity** | Contract testing, CI/CD pipelines |

## Quick Start Commands

```bash
# Start all services (future)
docker-compose -f docker-compose.microservices.yml up

# Individual service commands
docker-compose up weather-service
docker-compose up subscription-service
docker-compose up notification-service
docker-compose up scheduler-service
docker-compose up api-gateway
```

## Monitoring Endpoints

Each service exposes:
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics
- `GET /ready` - Readiness probe

## Configuration Management

- Environment-based configuration
- Service-specific config files
- Centralized secrets management
- Feature flags per service

---

**Note**: This is a reference document. For detailed analysis, see [Design-003: Microservices Decomposition Analysis](003-microservices-decomposition.md). 