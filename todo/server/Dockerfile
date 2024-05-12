# Стадия сборки.
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копирование go.mod и go.sum и загрузка зависимостей.
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копирование всех файлов проекта.
COPY . .

# Сборка приложения.
RUN CGO_ENABLED=0 go build -o myserver ./server/main.go

# Стадия запуска.
FROM alpine:latest

WORKDIR /app

# Копирование собранного приложения из стадии сборки.
COPY --from=builder /app/myserver .

# Копирование статических файлов.
COPY ./static /app/static

# Команда запуска приложения.
CMD ["/app/myserver"]