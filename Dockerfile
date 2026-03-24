FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник и миграции
COPY --from=builder /app/app .
COPY --from=builder "/app/cmd/config/config.yaml" "/app/cmd/config/config.yaml"
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]