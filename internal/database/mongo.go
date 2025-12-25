package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	database *mongo.Database
)

// Connect initializes MongoDB connection from environment variables
func Connect(ctx context.Context) error {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DATABASE")

	if mongoURI == "" {
		return fmt.Errorf("MONGO_URI environment variable is required")
	}

	if dbName == "" {
		dbName = "db-backups"
	}

	// Set client options
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(10).
		SetMinPoolSize(2).
		SetMaxConnIdleTime(5 * time.Minute)

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database = client.Database(dbName)
	return nil
}

// GetDatabase returns the MongoDB database instance
func GetDatabase() *mongo.Database {
	return database
}

// Disconnect closes the MongoDB connection
func Disconnect(ctx context.Context) error {
	if client != nil {
		return client.Disconnect(ctx)
	}
	return nil
}

// HealthCheck verifies the MongoDB connection is alive
func HealthCheck(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("MongoDB client is not initialized")
	}
	return client.Ping(ctx, nil)
}
