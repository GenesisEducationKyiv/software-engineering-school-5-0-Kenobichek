# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные исходники
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-service ./cmd/notification-service

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/notification-service /app/notification-service
COPY --from=builder /app/config /app/config
COPY --from=builder /app/domain /app/domain
COPY --from=builder /app/infrastructure /app/infrastructure
COPY --from=builder /app/notifier /app/notifier

EXPOSE ${PORT}

ENV PORT=${PORT}

ENTRYPOINT ["/app/notification-service"]