FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscription-service ./cmd/subscription-service/main.go

CMD ["go","run","./cmd/subscription-service/main.go"]