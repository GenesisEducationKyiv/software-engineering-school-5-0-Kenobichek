# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o subscription-service ./cmd/subscription-service

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/subscription-service /app/subscription-service
COPY --from=builder /app/migrations /app/migrations
EXPOSE 8083
ENTRYPOINT ["/app/subscription-service"] 