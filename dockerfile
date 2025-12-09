# Этап сборрки с тегированием начала.
FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /app

# Копирование зависимостей и инициализация их скачивания.
COPY go.mod go.sum ./
RUN go mod download

# Копирвоание кода из локально ветки go-my-knnowledge-base.
COPY . .

#Компиляция бинарника. 
# Сборка CGO_ENABLED=0 - бинарник без зависимостей. GOOS=linux- кросс-компиляция для линкус, даже если на windows локалка. -o knowledge-base убирает отладочную информацию. -o knowledge-base \ - имя созданного бинарника.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w" \
    -o knowledge-base \
    ./cmd/api/main.go

#Этап запуска
FROM alpine:3.19

WORKDIR /app

# Копируем только бинарник и миграции.
COPY --from=builder /app/knowledge-base .
COPY --from=builder /app/migrations ./migrations 

EXPOSE 2709
CMD ["./knowledge-base"]