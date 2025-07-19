FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o weather-service ./cmd/weather-service

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/weather-service .
CMD ["./weather-service"] 