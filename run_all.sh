#!/bin/bash

set -e

# Start all services in detached mode
echo "Starting Kafka..."
docker-compose -f docker-compose.kafka.yaml up -d

echo "Starting weather-service..."
docker-compose -f internal/services/weather-service/docker-compose.yaml up -d

echo "Starting subscription-service..."
docker-compose -f internal/services/subscription-service/docker-compose.yaml up -d

echo "Starting notification-service..."
docker-compose -f internal/services/notification-service/docker-compose.yaml up -d

echo "Starting api-gateway..."
docker-compose -f internal/services/api-gateway/docker-compose.yaml up -d

echo "All services are up and running!"
echo

# Show logs for all services in parallel
# echo "Kafka logs:"
# docker-compose -f docker-compose.kafka.yaml logs -f &

echo "weather-service logs:"
docker-compose -f internal/services/weather-service/docker-compose.yaml logs -f &

echo "subscription-service logs:"
docker-compose -f internal/services/subscription-service/docker-compose.yaml logs -f &

echo "notification-service logs:"
docker-compose -f internal/services/notification-service/docker-compose.yaml logs -f &

echo "api-gateway logs:"
docker-compose -f internal/services/api-gateway/docker-compose.yaml logs -f &

wait 
