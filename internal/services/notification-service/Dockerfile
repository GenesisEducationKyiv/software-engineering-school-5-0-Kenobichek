FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-service ./cmd/notification-service/main.go

CMD ["go","run","./cmd/notification-service/main.go"]