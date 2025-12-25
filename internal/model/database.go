package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database represents a saved database connection configuration
type Database struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id" swaggertype:"string"`
	Name           string             `bson:"name" json:"name" example:"My Production DB"`
	Type           BackupType         `bson:"type" json:"type" example:"postgre"`
	Host           string             `bson:"host" json:"host" example:"localhost"`
	Port           string             `bson:"port" json:"port" example:"5432"`
	Username       string             `bson:"username" json:"username" example:"user"`
	Password       string             `bson:"password" json:"password" example:"pass"`
	Database       string             `bson:"database" json:"database" example:"mydb"`
	ConnectionURI  string             `bson:"connectionUri" json:"connectionUri" example:"mongodb://user:pass@host:port/db"`
	CronExpression string             `bson:"cronExpression" json:"cronExpression" example:"0 0 * * *"`
	IsActive       bool               `bson:"isActive" json:"isActive" example:"true"`
	WebhookURL     string             `bson:"webhookUrl" json:"webhookUrl" example:"http://example.com/webhook"`
	CreatedAt      primitive.DateTime `bson:"createdAt" json:"createdAt" swaggertype:"string"`
	UpdatedAt      primitive.DateTime `bson:"updatedAt" json:"updatedAt" swaggertype:"string"`
}

// CreateDatabaseRequest represents the request body for creating a database
type CreateDatabaseRequest struct {
	Name           string     `json:"name" validate:"required" example:"My Production DB"`
	Type           BackupType `json:"type" validate:"required" example:"postgre"`
	Host           string     `json:"host" example:"localhost"`
	Port           string     `json:"port" example:"5432"`
	Username       string     `json:"username" example:"user"`
	Password       string     `json:"password" example:"pass"`
	Database       string     `json:"database" example:"mydb"`
	ConnectionURI  string     `json:"connectionUri" example:"mongodb://user:pass@host:port/db"`
	CronExpression string     `json:"cronExpression" example:"0 0 * * *"`
	IsActive       bool       `json:"isActive" example:"true"`
	WebhookURL     string     `json:"webhookUrl" example:"http://example.com/webhook"`
}

// UpdateDatabaseRequest represents the request body for updating a database
type UpdateDatabaseRequest struct {
	Name           string     `json:"name" example:"My Production DB"`
	Type           BackupType `json:"type" example:"postgre"`
	Host           string     `json:"host" example:"localhost"`
	Port           string     `json:"port" example:"5432"`
	Username       string     `json:"username" example:"user"`
	Password       string     `json:"password" example:"pass"`
	Database       string     `json:"database" example:"mydb"`
	ConnectionURI  string     `json:"connectionUri" example:"mongodb://user:pass@host:port/db"`
	CronExpression string     `json:"cronExpression" example:"0 0 * * *"`
	IsActive       bool       `json:"isActive" example:"true"`
	WebhookURL     string     `json:"webhookUrl" example:"http://example.com/webhook"`
}

// DatabaseListResponse represents a paginated list of databases
type DatabaseListResponse struct {
	Databases []Database `json:"databases"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
}
