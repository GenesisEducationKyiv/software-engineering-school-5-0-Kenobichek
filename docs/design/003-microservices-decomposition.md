# Design-003: Microservices Decomposition & Event-Driven Architecture

## Overview

Refactoring the monolithic Weather Forecast application into a scalable, resilient, and asynchronous microservices architecture using Kafka for event streaming and communication. The Weather Service remains synchronous for real-time requests, while other interactions use asynchronous event streams.

## Microservices and Responsibilities

### 1. Weather Service (Synchronous)

* Retrieves weather data from external APIs.
* Provides synchronous REST endpoints for real-time client queries.
* Caches latest weather data locally (Redis).
* Publishes `WeatherUpdated` events with full weather data state to Kafka.

### 2. Subscription Service

* Manages user subscriptions (create, confirm, cancel).
* Stores subscription data.
* Publishes events: `SubscriptionCreated`, `SubscriptionConfirmed`, `SubscriptionCancelled` to Kafka with full subscription state.

### 3. Notification Service

* Subscribes to `WeatherUpdated` and subscription-related events from Kafka.
* Sends notifications asynchronously via external services (e.g., SendGrid).
* Publishes `NotificationSent` events to Kafka for logging/auditing.

### 4. Scheduler Service

* Subscribes to `WeatherUpdated` and subscription-related events from Kafka.
* Manages scheduling of notifications and background tasks.
* Handles retries, error handling, and task management internally.

### 5. API Gateway

* Provides RESTful API endpoints for external client interactions.
* Serves GET requests (e.g., cached weather data) from a local read-model updated asynchronously via Kafka.
* Publishes commands (e.g., `CreateSubscriptionCommand`) to Kafka for write operations.
* Handles basic validation, authentication, and rate limiting.

## Inter-Service Communication via Kafka

All services, except direct client requests to Weather Service, communicate through Kafka:

| Kafka Topic              | Publisher            | Consumers                            |
| ------------------------ | -------------------- | ------------------------------------ |
| `weather.updated`        | Weather Service      | Notification, Scheduler, API Gateway |
| `subscription.created`   | Subscription Service | Notification, Scheduler              |
| `subscription.confirmed` | Subscription Service | Notification, Scheduler              |
| `subscription.cancelled` | Subscription Service | Notification, Scheduler              |
| `notification.sent`      | Notification Service | Audit/Logging                        |
| `commands.subscription`  | API Gateway          | Subscription Service                 |

## Example Interaction Flow

### Subscription Flow

* Client sends subscription request (`POST`) to API Gateway.
* API Gateway validates request and publishes `CreateSubscriptionCommand` to Kafka.
* Subscription Service processes command, stores data, and publishes `SubscriptionCreated` event.
* Notification Service consumes event and sends a confirmation email.
* Scheduler Service schedules follow-up notifications.

### Weather Update Flow

* Weather Service fetches new weather data, updates local cache, and publishes `WeatherUpdated` event.
* Notification Service receives event and sends notifications if required.
* Scheduler Service schedules tasks based on updated weather.
* API Gateway updates its local cache to serve future GET requests quickly.

## Error Handling & Reliability

* Kafka ensures at-least-once delivery of events.
* Dead-letter queues (DLQ) handle persistent failures.
* Services implement idempotent event processing to avoid duplication.
* Monitoring and alerts for failures, delays, and abnormal DLQ growth.


## Benefits

* Scalability: Independent service scaling.
* Resilience: Fault isolation and asynchronous processing.
* Maintainability: Clear separation of concerns and easy service evolution.
