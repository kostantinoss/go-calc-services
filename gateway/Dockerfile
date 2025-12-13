FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.mod

RUN go mod download

COPY . .

RUN go build -o gateway

# 2. Runtime stage
FROM alpine:latest

COPY --from=builder /app/gateway /
COPY config.yaml /app/config.yaml

ENTRYPOINT [ "/gateway" ]