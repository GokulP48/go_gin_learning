# syntax=docker/dockerfile:1
FROM golang:1.21-alpine

WORKDIR /app

# Install git and deps
RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Set config path env (optional)
ENV CONFIG_PATH=/config/config.yaml

RUN go build -o main ./cmd

CMD ["./main"]
