# Стадия сборки
FROM golang:1.24 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o subscription_service ./cmd/main.go

# Стадия финального образа
FROM alpine:3.21

WORKDIR /app
RUN adduser -D -g ''
COPY --from=builder /app/subscription_service .
RUN chown -R appuser:appuser /app
USER appuser

EXPOSE 8080
CMD ["./subscription_service"]
