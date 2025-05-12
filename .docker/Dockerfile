FROM golang:1.23-alpine AS builder

WORKDIR /app

# depen
COPY go.mod go.sum ./
RUN go mod download

# build
COPY . .
RUN go build -o app cmd/schedule-app/main.go

EXPOSE 8080 9000

CMD ["./app"]