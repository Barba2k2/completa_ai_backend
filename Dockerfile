# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.19

WORKDIR /app

RUN adduser -D -g '' appuser

COPY --from=builder /app/server /app/server
COPY --from=builder /app/migrations /app/migrations

USER appuser

ENV PORT=8080
EXPOSE 8080

CMD ["/app/server"]
