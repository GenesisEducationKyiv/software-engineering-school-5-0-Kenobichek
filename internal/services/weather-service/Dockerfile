FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o weather-service ./cmd/weather-service/main.go

CMD ["go","run","./cmd/weather-service/main.go"]