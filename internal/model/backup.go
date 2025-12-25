package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BackupType string

const (
	Postgres BackupType = "postgre"
	MySQL    BackupType = "mysql"
	Mongo    BackupType = "mongo"
	Redis    BackupType = "redis"
)

type BackupRequest struct {
	Type          BackupType `json:"type" example:"postgre"`
	Host          string     `json:"host" example:"localhost"`
	Port          string     `json:"port" example:"5432"`
	Username      string     `json:"username" example:"user"`
	Password      string     `json:"password" example:"pass"`
	WebhookURL    string     `json:"webhookUrl" example:"http://example.com/webhook"`
	Database      string     `json:"database" example:"mydb"`
	ConnectionURI string     `json:"connectionUri" example:"mongodb://user:pass@host:port/db"`
	AuthSource    string     `json:"authSource" example:"admin"`
}

type BackupResult struct {
	Success   bool              `json:"success"`
	Error     string            `json:"error,omitempty"`
	FilePath  string            `json:"filePath"`
	ObjectKey string            `json:"objectKey,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp string            `json:"timestamp"`
}

type BackupResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	ID      string `json:"id,omitempty"`
}

// BackupStatus represents the current status of a backup
type BackupStatus string

const (
	StatusPending    BackupStatus = "pending"
	StatusGenerating BackupStatus = "generating"
	StatusCompleted  BackupStatus = "completed"
	StatusFailed     BackupStatus = "failed"
)

// BackupMetadata represents backup information stored in MongoDB
type BackupMetadata struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      string             `bson:"type" json:"type"`
	ObjectKey string             `bson:"objectKey" json:"objectKey"`
	FilePath  string             `bson:"filePath" json:"filePath"`
	FileSize  int64              `bson:"fileSize" json:"fileSize"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Status    BackupStatus       `bson:"status" json:"status"` // pending, generating, completed, failed
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`
	Host      string             `bson:"host" json:"host"`
	Database  string             `bson:"database" json:"database"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt" swaggertype:"string"`
}

// BackupListResponse represents paginated backup list
type BackupListResponse struct {
	Backups []BackupMetadata `json:"backups"`
	Total   int64            `json:"total"`
	Page    int              `json:"page"`
	Limit   int              `json:"limit"`
}
