FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /go/bin && go build -o /go/bin/app ./cmd/main.go

CMD ["/go/bin/app"]