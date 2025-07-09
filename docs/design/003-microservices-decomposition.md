# Design-003: Microservices Decomposition & Event-Driven Architecture

## 1. Overview

The Weather Forecast API is currently a monolithic Go application providing weather subscription services. Its Clean Architecture makes it well-suited for microservices decomposition.

### Core Components:
- **Weather Service**: Aggregates weather data from external providers, with caching.
- **Subscription Service**: Manages user subscriptions and their lifecycle.
- **Notification Service**: Delivers notifications via email (SendGrid).
- **Scheduler Service**: Orchestrates background jobs for weather updates.
- **Database**: PostgreSQL for persistent data.
- **Cache**: Redis for weather data.
- **Monitoring**: Prometheus and Grafana.

---

## 2. Microservices Decomposition

### Guiding Principles
- **Domain-Driven Boundaries**: Each service encapsulates a distinct business capability.
- **Autonomous Services**: Each service owns its data and logic.
- **Event-Driven Communication**: All inter-service communication is asynchronous, via a message broker (Kafka, NATS, RabbitMQ, etc.).
- **Event-Carried State Transfer**: Events include the full relevant entity state, eliminating the need for consumers to query producers for additional data.

### Identified Microservices

#### A. Weather Service
- **Responsibilities**:
  - Aggregate weather data from multiple providers.
  - Implement provider fallback logic.
  - Cache weather data in Redis.
  - **Publish** `WeatherUpdated` events (with full weather state) to the message broker.
- **Data Ownership**: Weather cache.
- **Dependencies**: External weather APIs, Redis, Message Broker.

#### B. Subscription Service
- **Responsibilities**:
  - Manage the full subscription lifecycle (create, confirm, cancel).
  - Store subscription and template data.
  - **Publish** `SubscriptionCreated`, `SubscriptionConfirmed`, and `SubscriptionCancelled` events (with full subscription state).
- **Data Ownership**: Subscriptions, email templates.
- **Dependencies**: PostgreSQL, Message Broker.

#### C. Notification Service
- **Responsibilities**:
  - **Subscribe** to relevant events (e.g., `SubscriptionCreated`, `WeatherUpdated`).
  - Deliver notifications via SendGrid (and future channels).
  - Manage notification templates and delivery logs.
  - **Publish** `NotificationSent` events (for audit/logging).
- **Data Ownership**: Notification templates, delivery logs.
- **Dependencies**: SendGrid, PostgreSQL, Message Broker.

#### D. Scheduler Service
- **Responsibilities**:
  - **Subscribe** to events to schedule jobs (e.g., weather update notifications).
  - Manage job scheduling, retries, and failures.
- **Data Ownership**: Job scheduling data.
- **Dependencies**: PostgreSQL, Message Broker.

#### E. API Gateway
- **Responsibilities**:
  - Route and aggregate client requests.
  - Publish commands/events to the broker.
  - Handle authentication, authorization, rate limiting, and request transformation.
- **Data Ownership**: API configuration, routing rules.
- **Dependencies**: Message Broker.

---

## 3. Event-Driven Communication

### Event Streaming & State Transfer
- All services communicate asynchronously by publishing and subscribing to events via a message broker.
- Events use **event-carried state transfer**: each event contains the full, relevant entity state, so consumers never need to query the producer for additional data.
- Each service maintains its own local state, updated solely from the event stream.

### Reliability, Failure Handling, and Monitoring
- **Delivery Guarantees**: The message broker provides at-least-once delivery. Events are persisted until successfully processed by all subscribers.
- **Retry Mechanisms**: If a subscriber is temporarily unavailable or fails to process an event, the broker automatically retries delivery until acknowledgment or until a configurable retention period expires.
- **Dead Letter Queue (DLQ)**: Events that cannot be processed after multiple attempts are routed to a DLQ for further inspection and manual intervention.
- **Idempotency**: All event handlers must be idempotent to ensure that repeated deliveries do not cause inconsistent state or duplicate side effects.
- **Monitoring & Alerting**: The system continuously monitors event delivery metrics, DLQ size, and subscriber health. Alerts are triggered on delivery failures, processing delays, or abnormal DLQ growth.
- **Observability**: Distributed tracing and centralized logging are used to track event flow and diagnose issues across services.

#### Example Event Flows
- **Weather Service** publishes `WeatherUpdated` events (full weather data).
- **Subscription Service** publishes `SubscriptionCreated`, `SubscriptionConfirmed`, `SubscriptionCancelled` events (full subscription state).
- **Notification Service** subscribes to these events and sends notifications, then publishes `NotificationSent` events.
- **Scheduler Service** subscribes to events to schedule jobs.

### Event Topics Matrix

| Event Topic            | Published By         | Consumed By                |
|------------------------|----------------------|----------------------------|
| weather.updated        | Weather Service      | Notification, Scheduler    |
| subscription.created   | Subscription Service | Notification, Scheduler    |
| subscription.cancelled | Subscription Service | Notification, Scheduler    |
| notification.sent      | Notification Service | (Audit/Logging)            |

---

## 4. Service Specifications

### 4.1 Weather Service
- **Service Name**: `weather-service`
- **Port**: 8081
- **Protocol**: HTTP/gRPC (for health/metrics only)
- **Database**: Redis (cache)
- **Endpoints**:
  - `GET /weather?city={city}` – Get current weather
  - `GET /health` – Health check
  - `GET /metrics` – Prometheus metrics
- **Publishes**: `WeatherUpdated` events (full weather state)
- **Dependencies**: External weather APIs, Redis, Message Broker

### 4.2 Subscription Service
- **Service Name**: `subscription-service`
- **Port**: 8082
- **Protocol**: HTTP/gRPC (for health/metrics only)
- **Database**: PostgreSQL
- **Endpoints**:
  - `POST /subscriptions` – Create subscription
  - `PUT /subscriptions/{token}/confirm` – Confirm subscription
  - `DELETE /subscriptions/{token}` – Unsubscribe
  - `GET /subscriptions/due` – Get due subscriptions
  - `GET /health` – Health check
  - `GET /metrics` – Prometheus metrics
- **Publishes**: `SubscriptionCreated`, `SubscriptionConfirmed`, `SubscriptionCancelled` events (full subscription state)
- **Dependencies**: PostgreSQL, Message Broker

### 4.3 Notification Service
- **Service Name**: `notification-service`
- **Port**: 8083
- **Protocol**: HTTP/gRPC (for health/metrics only)
- **Database**: PostgreSQL (templates, logs)
- **Endpoints**:
  - `POST /notifications/weather` – Send weather update
  - `POST /notifications/confirm` – Send confirmation
  - `POST /notifications/unsubscribe` – Send unsubscribe
  - `GET /templates` – Get notification templates
  - `GET /health` – Health check
  - `GET /metrics` – Prometheus metrics
- **Subscribes**: `SubscriptionCreated`, `SubscriptionConfirmed`, `SubscriptionCancelled`, `WeatherUpdated`
- **Publishes**: `NotificationSent` events (full notification state)
- **Dependencies**: SendGrid, PostgreSQL, Message Broker

### 4.4 Scheduler Service
- **Service Name**: `scheduler-service`
- **Port**: 8084
- **Protocol**: HTTP/gRPC (for health/metrics only)
- **Database**: PostgreSQL (job tracking)
- **Endpoints**:
  - `POST /scheduler/start` – Start scheduler
  - `POST /scheduler/stop` – Stop scheduler
  - `GET /scheduler/status` – Scheduler status
  - `GET /jobs` – List scheduled jobs
  - `GET /health` – Health check
  - `GET /metrics` – Prometheus metrics
- **Subscribes**: All relevant events for scheduling
- **Dependencies**: PostgreSQL, Message Broker

### 4.5 API Gateway
- **Service Name**: `api-gateway`
- **Port**: 8080
- **Protocol**: HTTP
- **Database**: None (stateless)
- **Endpoints**:
  - `GET /weather?city={city}` – Weather endpoint
  - `POST /subscribe` – Subscription endpoint
  - `GET /confirm/{token}` – Confirmation endpoint
  - `GET /unsubscribe/{token}` – Unsubscribe endpoint
  - `GET /health` – Health check
  - `GET /metrics` – Prometheus metrics
- **Publishes**: Commands/events to the broker
- **Dependencies**: All microservices (via broker), Rate limiting, Authentication

---

## 5. Data Management
- **Database per Service**: Each service owns its database (no cross-service DB access).
- **Eventual Consistency**: Data is synchronized across services via events.
- **Event-Carried State Transfer**: All events include the full entity state.
- **No Direct Data Calls**: Services never query each other for data; all state is received via events.

---

## 6. Deployment & Infrastructure
- **Containerization**: Each service runs in its own Docker container.
- **Shared Infrastructure**: Redis, PostgreSQL, Message Broker.
- **Service Discovery & Load Balancing**: Managed by orchestration platform (e.g., Kubernetes).
- **Monitoring & Observability**: Distributed tracing, centralized logging, metrics, and health checks.
- **Security**: mTLS for internal communication, JWT for external access, secrets management.

---

## 7. Benefits
- **Scalability**: Independent scaling of services.
- **Maintainability**: Smaller, focused codebases; independent deployments.
- **Resilience**: Fault isolation, circuit breakers, graceful degradation.
- **Team Autonomy**: Teams own services, enabling parallel development.

---

## 8. Conclusion

This microservices decomposition and event-driven architecture enable scalable, maintainable, and resilient evolution of the current system. By leveraging asynchronous event streaming and event-carried state transfer, services remain autonomous, loosely coupled, and robust to change. 