FROM golang:1.23.0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api-gateway ./cmd/api-gateway/main.go

CMD ["go","run","./cmd/api-gateway/main.go"]