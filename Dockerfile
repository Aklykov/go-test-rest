# Build Stage
FROM golang:1.25.7-alpine AS builder
WORKDIR /app
COPY . .

# Устанавливаем Goose и применяем миграции
RUN apk add --no-cache postgresql-client && \
    go install github.com/pressly/goose/v3/cmd/goose@latest

# компилируем исполнительный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final Stage
FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# для свагера
COPY --from=builder /app/docs ./docs
# для миграций
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/goose .
EXPOSE ${APP_PORT}
#CMD ["./main"]
# Запускаем миграции (опционально) и приложение
CMD [ \
    "sh", "-c", "until pg_isready -h ${DB_HOST} -U ${DB_USER} -d ${DB_NAME}; do sleep 1; done && \
    ./goose postgres \"host=$DB_HOST port=$DB_PORT user=$DB_USER dbname=$DB_NAME password=$DB_PASSWORD sslmode=$DB_SSLMODE\" -dir ./migrations up && \
    ./main" \
]
