# Build stage
# Pin the builder to the host's native platform and cross-compile to the
# target arch. This keeps the Go compiler running natively (no QEMU), avoiding
# segfaults when building linux/amd64 + linux/arm64 via buildx.
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Provided automatically by buildx for the image being produced.
ARG TARGETOS TARGETARCH

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o db-backup cmd/server/main.go

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
