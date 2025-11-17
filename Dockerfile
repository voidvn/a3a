FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# ← ЭТА СТРОКА ОТВЕЧАЕТ ЗА СБОРКУ
RUN go build -o s4s-backend ./cmd/main.go

# ← Финальный образ
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/s4s-backend .

# ВОТ ЭТО САМОЕ ВАЖНОЕ — копируем .env внутрь образа!
COPY .env .env

EXPOSE 8080
CMD ["./s4s-backend"]
