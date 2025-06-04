FROM golang:latest

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]