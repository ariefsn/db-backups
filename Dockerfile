# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o db-backup cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install database client tools
# postgresql-client: for pg_dump
# mariadb-client: for mysqldump
# mongodb-tools: for mongodump
# redis: for redis-cli
RUN apk add --no-cache \
    postgresql-client \
    mariadb-client \
    mongodb-tools \
    redis \
    ca-certificates \
    tzdata

COPY --from=builder /app/db-backup .

EXPOSE 8080

ENTRYPOINT ["./db-backup"]
