# Db Backup

A generic database backup application capable of backing up PostgreSQL, MySQL, MongoDB, and Redis databases. Features background processing, webhook notifications, automatic daily cleanup, **cloud storage integration (Cloudflare R2)**, and **backup management with MongoDB**.

[GitHub](https://github.com/ariefsn/db-backups) | [Docker Hub](https://hub.docker.com/r/ariefsn/db-backup)

## Features

- **Multi-Database Support**: PostgreSQL, MySQL, MongoDB, Redis.
- **Background Backups**: Non-blocking backup operations.
- **Cloud Storage**: Automatic upload to Cloudflare R2 (S3-compatible).
- **Backup Management**: MongoDB-backed metadata storage with pagination.
- **Webhook Notifications**: Receive JSON payloads with object keys and metadata upon backup completion or failure.
- **REST API**: List and delete backups via API endpoints.
- **Automatic Cleanup**: Daily cron job to remove backups older than 7 days.
- **Swagger Documentation**: Interactive API docs.
- **Docker Ready**: Pre-built image with all necessary database tools.

## Getting Started

### Prerequisites

- Go 1.23+
- Database tools (`pg_dump`, `mysqldump`, `mongodump`, `redis-cli`) if running locally without Docker.
- MongoDB instance (for backup metadata storage)
- Cloudflare R2 bucket (optional, for cloud storage)

### Environment Variables

#### Required for MongoDB
- `MONGO_URI` - MongoDB connection string (e.g., `mongodb://localhost:27017`)
- `MONGO_DATABASE` - MongoDB database name (default: `db-backup`)

#### Optional for R2 Storage
- `R2_ENDPOINT` - Cloudflare R2 endpoint URL
- `R2_ACCESS_KEY_ID` - R2 access key
- `R2_SECRET_ACCESS_KEY` - R2 secret access key
- `R2_BUCKET_NAME` - R2 bucket name
- `R2_REGION` - R2 region (default: `auto`)

### Local Run

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ariefsn/db-backup.git
   cd db-backup
   ```

2. **Set environment variables**:
   ```bash
   export MONGO_URI="mongodb://localhost:27017"
   export MONGO_DATABASE="db-backup"
   export R2_ENDPOINT="https://your-account.r2.cloudflarestorage.com"
   export R2_ACCESS_KEY_ID="your-access-key"
   export R2_SECRET_ACCESS_KEY="your-secret-key"
   export R2_BUCKET_NAME="db-backups"
   ```

3. **Run the server**:
   ```bash
   make run
   # OR
   go run cmd/server/main.go
   ```

4. **Access API**:
   - Swagger: [http://localhost:8080/swagger](http://localhost:8080/swagger)
   - Health: [http://localhost:8080/health](http://localhost:8080/health)

### Docker Run

Pull the image from Docker Hub:

```bash
docker pull ariefsn/db-backup:latest
docker run --rm -p 8080:8080 \
  -e MONGO_URI="mongodb://host.docker.internal:27017" \
  -e R2_ENDPOINT="https://your-account.r2.cloudflarestorage.com" \
  -e R2_ACCESS_KEY_ID="your-key" \
  -e R2_SECRET_ACCESS_KEY="your-secret" \
  -e R2_BUCKET_NAME="db-backups" \
  ariefsn/db-backup:latest
```

### Docker Compose (Full Stack)

To run both the backend and web interface:

```bash
docker compose up -d
```

This will run:
- Backend: http://localhost:8080
- Web Interface: http://localhost:3000

## API Usage

### Trigger Backup

**POST** `/backup`

```json
{
  "type": "postgre",
  "host": "postgres-host",
  "port": "5432",
  "username": "user",
  "password": "password",
  "database": "mydb",
  "webhookUrl": "https://your-webhook.com/callback"
}
```

**Supported Types**: `postgre`, `mysql`, `mongo`, `redis`

**Response**: 202 Accepted

### List Backups

**GET** `/backups?page=1&limit=10&statuses=completed,failed`

Returns paginated list of backup metadata with optional status filtering.

**Query Parameters**:
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `statuses` - Comma-separated status values: `pending`, `generating`, `completed`, `failed`

**Response**:
```json
{
  "backups": [
    {
      "id": "507f1f77bcf86cd799439011",
      "type": "postgre",
      "objectKey": "backups/postgre/20231225-120000_mydb.sql",
      "filePath": "/backups/mydb_20231225_120000.sql",
      "fileSize": 1024000,
      "timestamp": "2023-12-25T12:00:00Z",
      "status": "completed",
      "host": "postgres-host",
      "database": "mydb"
    }
  ],
  "total": 42,
  "page": 1,
  "limit": 10
}
```

**Status Filtering Examples**:
```bash
# Get only completed backups
GET /backups?statuses=completed

# Get failed and pending backups
GET /backups?statuses=failed,pending

# Get backups currently being generated
GET /backups?statuses=generating
```

### Get Single Backup

**GET** `/backups/{id}`

Retrieve a single backup by ID.

**Response**:
```json
{
  "id": "507f1f77bcf86cd799439011",
  "type": "postgre",
  "objectKey": "backups/postgre/20231225-120000_mydb.sql",
  "filePath": "/backups/mydb_20231225_120000.sql",
  "fileSize": 1024000,
  "timestamp": "2023-12-25T12:00:00Z",
  "status": "completed",
  "host": "postgres-host",
  "database": "mydb"
}
```

### Delete Backup

**DELETE** `/backups/{id}`

Deletes backup from both MongoDB and R2 storage.

**Response**: 200 OK

### Webhook Payload

When a backup completes, the webhook receives:

```json
{
  "success": true,
  "filePath": "/backups/mydb_20231225_120000.sql",
  "objectKey": "backups/postgre/20231225-120000_mydb.sql",
  "timestamp": "2023-12-25T12:00:00Z",
  "metadata": {
    "database_type": "postgre",
    "host": "postgres-host",
    "database": "mydb",
    "file_size": "1024000",
    "storage": "r2"
  }
}
```

## Development

### Makefile Commands

- `make build`: Build the binary.
- `make run`: Run the application.
- `make docker-build`: Build local Docker image.
- `make docker-push`: Build and push multi-platform Docker image.
- `make swagger`: Regenerate Swagger documentation.

## License

MIT
