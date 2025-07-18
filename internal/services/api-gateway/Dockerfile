FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-gateway ./cmd/api-gateway

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api-gateway .
CMD ["./api-gateway"] 