# 004: Microservice Architecture

## Overview

This document outlines the microservice architecture for the Weather Forecast application. The system provides users with real-time weather data and scheduled notifications.

A hybrid architecture is used, combining synchronous (gRPC) and asynchronous (Kafka) communication. This model supports both immediate, real-time data lookups and a scalable, event-driven backbone for background tasks like notifications.

## Core Components

The system is composed of four main services:

-   **API Gateway**: The single entry point for all client HTTP requests. It routes requests, and for write operations, publishes commands to Kafka.
-   **Weather Service**: Fetches weather data from external APIs and serves it to the API Gateway via a synchronous gRPC endpoint.
-   **Subscription Service**: Manages user subscriptions. It consumes commands from Kafka, updates its database, and publishes the results as events back to Kafka.
-   **Notification Service**: Listens for events on Kafka (e.g., `subscription.confirmed`) and sends notifications to users.

## Communication Patterns

-   **Synchronous (gRPC)**: Used for real-time, request-response queries. The `API Gateway` calls the `Weather Service` directly via gRPC to get current weather data instantly.
-   **Asynchronous (Kafka)**: Used for all other inter-service communication. Services are decoupled and communicate through events and commands on Kafka topics, making the system resilient and scalable.

## Event Processing Example: New Subscription

1.  **Request**: A client sends an HTTP request to the `API Gateway` to create a subscription.
2.  **Command**: The `API Gateway` publishes a `create-subscription` command to Kafka.
3.  **Processing**: The `Subscription Service` consumes the command, creates the subscription in its database, and publishes a `subscription-confirmed` event.
4.  **Notification**: The `Notification Service` consumes the event and sends a confirmation email to the user. 