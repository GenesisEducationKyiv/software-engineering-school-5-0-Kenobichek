# Design-003: Microservices Decomposition Analysis

## 1. Current Application Overview

The Weather Forecast API is a monolithic application built with Go that provides weather subscription services. The application follows Clean Architecture principles and is well-structured for microservices extraction.

### Current Architecture Components:
- **Weather Service**: Fetches weather data from multiple providers (OpenWeather, WeatherAPI) with caching
- **Subscription Service**: Manages user subscriptions, confirmations, and unsubscriptions
- **Notification Service**: Handles email notifications via SendGrid
- **Scheduler Service**: Runs background jobs to send weather updates
- **Database**: PostgreSQL for subscriptions and email templates
- **Cache**: Redis for weather data caching
- **Monitoring**: Prometheus and Grafana for observability

## 2. Microservices Decomposition Strategy

### 2.1 Identified Microservices

Based on domain boundaries and business capabilities, the following microservices should be extracted:

#### **A. Weather Service**
- **Responsibility**: Weather data aggregation and caching
- **Bounded Context**: Weather information domain
- **Core Functions**:
  - Fetch weather data from multiple providers (OpenWeather, WeatherAPI)
  - Implement provider fallback chain
  - Cache weather data in Redis
  - Provide weather data via REST API
- **Data Ownership**: Weather cache data
- **Dependencies**: External weather APIs, Redis cache

#### **B. Subscription Service**
- **Responsibility**: User subscription management
- **Bounded Context**: Subscription domain
- **Core Functions**:
  - Create, confirm, and cancel subscriptions
  - Manage subscription lifecycle
  - Store subscription data
  - Provide subscription status
- **Data Ownership**: Subscription data, email templates
- **Dependencies**: PostgreSQL database

#### **C. Notification Service**
- **Responsibility**: Multi-channel notification delivery
- **Bounded Context**: Notification domain
- **Core Functions**:
  - Send email notifications via SendGrid
  - Template management and rendering
  - Support for multiple notification channels
  - Notification delivery tracking
- **Data Ownership**: Notification templates, delivery logs
- **Dependencies**: SendGrid API, template storage

#### **D. Scheduler Service**
- **Responsibility**: Background job orchestration
- **Bounded Context**: Scheduling domain
- **Core Functions**:
  - Schedule weather update notifications
  - Coordinate between services
  - Handle job retries and failures
  - Manage notification timing
- **Data Ownership**: Job scheduling data
- **Dependencies**: All other services

#### **E. API Gateway**
- **Responsibility**: Request routing and aggregation
- **Bounded Context**: API management domain
- **Core Functions**:
  - Route requests to appropriate services
  - Handle authentication and authorization
  - Rate limiting and throttling
  - Request/response transformation
- **Data Ownership**: API configuration, routing rules
- **Dependencies**: All microservices

## 3. Communication Methods & Rationale

### 3.1 Synchronous Communication (REST/gRPC)

#### **REST API** - Recommended for:
- **Weather Service**: External API calls to weather providers
- **API Gateway**: Client-facing endpoints
- **Service-to-Service**: Simple request/response patterns

**Rationale**: REST is widely supported, easy to debug, and suitable for stateless operations.

#### **gRPC** - Recommended for:
- **Internal Service Communication**: High-performance inter-service calls
- **Weather Service ↔ Scheduler Service**: Frequent weather data requests
- **Subscription Service ↔ Notification Service**: Real-time subscription events

**Rationale**: gRPC provides better performance, type safety, and streaming capabilities for internal communication.

### 3.2 Asynchronous Communication (Message Queue)

#### **Message Queue (RabbitMQ/Apache Kafka)** - Recommended for:
- **Scheduler Service → Notification Service**: Weather update notifications
- **Subscription Service → Notification Service**: Confirmation emails
- **Event-driven communication**: Subscription lifecycle events

**Rationale**: Message queues provide reliability, decoupling, and support for event-driven architecture.

### 3.3 Communication Matrix

| Service | Weather | Subscription | Notification | Scheduler | API Gateway |
|---------|---------|--------------|--------------|-----------|-------------|
| Weather | - | gRPC | - | gRPC | REST |
| Subscription | - | - | Message Queue | gRPC | REST |
| Notification | - | - | - | Message Queue | REST |
| Scheduler | gRPC | gRPC | Message Queue | - | - |
| API Gateway | REST | REST | REST | - | - |

## 4. Service Documentation

### 4.1 Weather Service

```yaml
Service Name: weather-service
Port: 8081
Protocol: HTTP/gRPC
Database: Redis (cache only)

Endpoints:
  GET /weather?city={city} - Get current weather for city
  GET /health - Health check
  GET /metrics - Prometheus metrics

Dependencies:
  - OpenWeather API
  - WeatherAPI
  - Redis (cache)
  - Prometheus (metrics)

Configuration:
  - Weather API keys
  - Cache TTL settings
  - Provider fallback configuration
```

**Key Features**:
- Provider chain with fallback
- Redis caching with TTL
- Circuit breaker pattern
- Metrics and monitoring

### 4.2 Subscription Service

```yaml
Service Name: subscription-service
Port: 8082
Protocol: HTTP/gRPC
Database: PostgreSQL

Endpoints:
  POST /subscriptions - Create subscription
  PUT /subscriptions/{token}/confirm - Confirm subscription
  DELETE /subscriptions/{token} - Unsubscribe
  GET /subscriptions/due - Get due subscriptions
  GET /health - Health check
  GET /metrics - Prometheus metrics

Dependencies:
  - PostgreSQL
  - Message Queue (for notifications)
  - Prometheus (metrics)

Configuration:
  - Database connection
  - Message queue settings
  - Subscription validation rules
```

**Key Features**:
- Subscription lifecycle management
- Token-based confirmation
- Due subscription queries
- Data validation and business rules

### 4.3 Notification Service

```yaml
Service Name: notification-service
Port: 8083
Protocol: HTTP/gRPC
Database: PostgreSQL (templates)

Endpoints:
  POST /notifications/weather - Send weather update
  POST /notifications/confirm - Send confirmation
  POST /notifications/unsubscribe - Send unsubscribe
  GET /templates - Get notification templates
  GET /health - Health check
  GET /metrics - Prometheus metrics

Dependencies:
  - SendGrid API
  - PostgreSQL (templates)
  - Message Queue (receiving)
  - Prometheus (metrics)

Configuration:
  - SendGrid API key
  - Template storage settings
  - Message queue settings
```

**Key Features**:
- Multi-channel notification support
- Template management
- Delivery tracking
- Retry mechanisms

### 4.4 Scheduler Service

```yaml
Service Name: scheduler-service
Port: 8084
Protocol: HTTP/gRPC
Database: PostgreSQL (job tracking)

Endpoints:
  POST /scheduler/start - Start scheduler
  POST /scheduler/stop - Stop scheduler
  GET /scheduler/status - Get scheduler status
  GET /jobs - List scheduled jobs
  GET /health - Health check
  GET /metrics - Prometheus metrics

Dependencies:
  - Weather Service (gRPC)
  - Subscription Service (gRPC)
  - Notification Service (Message Queue)
  - PostgreSQL (job tracking)
  - Prometheus (metrics)

Configuration:
  - Cron schedule settings
  - Service endpoints
  - Job retry configuration
```

**Key Features**:
- Cron-based job scheduling
- Service orchestration
- Job retry and failure handling
- Monitoring and alerting

### 4.5 API Gateway

```yaml
Service Name: api-gateway
Port: 8080
Protocol: HTTP
Database: None (stateless)

Endpoints:
  GET /weather?city={city} - Weather endpoint
  POST /subscribe - Subscription endpoint
  GET /confirm/{token} - Confirmation endpoint
  GET /unsubscribe/{token} - Unsubscribe endpoint
  GET /health - Health check
  GET /metrics - Prometheus metrics

Dependencies:
  - All microservices
  - Rate limiting service
  - Authentication service (future)

Configuration:
  - Service routing rules
  - Rate limiting settings
  - CORS configuration
```

**Key Features**:
- Request routing and aggregation
- Rate limiting and throttling
- CORS handling
- Request/response transformation

## 5. Data Management Strategy

### 5.1 Database Per Service
- **Weather Service**: Redis (cache only)
- **Subscription Service**: PostgreSQL (subscriptions table)
- **Notification Service**: PostgreSQL (templates table)
- **Scheduler Service**: PostgreSQL (job tracking table)

### 5.2 Data Consistency
- **Eventual Consistency**: Use message queues for cross-service data updates
- **Saga Pattern**: For complex multi-service transactions
- **CQRS**: Consider for read-heavy operations (weather data)

## 6. Deployment & Infrastructure

### 6.1 Containerization
- Each service as separate Docker container
- Shared infrastructure services (Redis, PostgreSQL, Message Queue)
- Service discovery and load balancing

### 6.2 Monitoring & Observability
- **Distributed Tracing**: Jaeger or Zipkin
- **Centralized Logging**: ELK Stack or Fluentd
- **Metrics**: Prometheus + Grafana
- **Health Checks**: Each service exposes /health endpoint

### 6.3 Security
- **Service-to-Service**: mTLS for internal communication
- **API Gateway**: JWT tokens for external access
- **Secrets Management**: Kubernetes secrets or HashiCorp Vault

## 7. Benefits of Microservices Architecture

### 7.1 Scalability
- Independent scaling of services
- Resource optimization based on load
- Horizontal scaling capabilities

### 7.2 Maintainability
- Smaller, focused codebases
- Independent deployment cycles
- Technology diversity per service

### 7.3 Resilience
- Fault isolation
- Circuit breaker patterns
- Graceful degradation

### 7.4 Team Organization
- Team ownership per service
- Parallel development
- Reduced coordination overhead

## 8. Conclusion

The proposed microservices decomposition provides a clear path for evolving the current monolithic application into a scalable, maintainable, and resilient system. The architecture maintains the benefits of the current Clean Architecture while enabling independent development, deployment, and scaling of individual services.

The communication strategy balances performance (gRPC for internal calls) with simplicity (REST for external APIs) and reliability (message queues for asynchronous operations). The phased migration approach minimizes risk while providing incremental value delivery. 