# 1. Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN go build -o calc_server

# 2. Runtime stage
FROM alpine:latest

COPY --from=builder /app/calc_server /

ENTRYPOINT [ "/calc_server" ]
