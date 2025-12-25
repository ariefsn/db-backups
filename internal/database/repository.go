package database

import (
	"context"
	"db-backup/internal/model"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const backupsCollection = "backups"

type Repository struct {
	db *mongo.Database
}

// NewRepository creates a new backup repository
func NewRepository() *Repository {
	return &Repository{
		db: GetDatabase(),
	}
}

// SaveBackup inserts a new backup record
func (r *Repository) SaveBackup(ctx context.Context, backup *model.BackupMetadata) error {
	collection := r.db.Collection(backupsCollection)

	// Set creation time if not set
	if backup.CreatedAt == 0 {
		backup.CreatedAt = primitive.NewDateTimeFromTime(backup.Timestamp)
	}

	result, err := collection.InsertOne(ctx, backup)
	if err != nil {
		return fmt.Errorf("failed to save backup: %w", err)
	}

	// Set the ID from the insert result
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		backup.ID = oid
	}

	return nil
}

// ListBackups retrieves backups with pagination and optional filtering
func (r *Repository) ListBackups(ctx context.Context, page, limit int, statuses []model.BackupStatus, types []string, search, orderBy, orderDir string, startDate, endDate *primitive.DateTime) ([]model.BackupMetadata, int64, error) {
	collection := r.db.Collection(backupsCollection)

	// Calculate skip
	skip := (page - 1) * limit

	// Build filter
	filter := bson.M{}

	// Status filtering
	if len(statuses) > 0 {
		filter["status"] = bson.M{"$in": statuses}
	}

	// Type filtering
	if len(types) > 0 {
		filter["type"] = bson.M{"$in": types}
	}

	// Search filtering (search in database, host, and type fields)
	if search != "" {
		filter["$or"] = []bson.M{
			{"database": bson.M{"$regex": search, "$options": "i"}},
			{"host": bson.M{"$regex": search, "$options": "i"}},
			{"type": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Date range filtering
	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}
		if startDate != nil {
			dateFilter["$gte"] = *startDate
		}
		if endDate != nil {
			dateFilter["$lte"] = *endDate
		}
		filter["createdAt"] = dateFilter
	}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count backups: %w", err)
	}

	// Determine sort order
	sortField := "createdAt"
	sortOrder := -1 // descending by default

	if orderBy != "" {
		sortField = orderBy
	}
	if orderDir == "asc" {
		sortOrder = 1
	}

	// Find with pagination and sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list backups: %w", err)
	}
	defer cursor.Close(ctx)

	var backups []model.BackupMetadata
	if err := cursor.All(ctx, &backups); err != nil {
		return nil, 0, fmt.Errorf("failed to decode backups: %w", err)
	}

	return backups, total, nil
}

// GetBackup retrieves a single backup by ID
func (r *Repository) GetBackup(ctx context.Context, id string) (*model.BackupMetadata, error) {
	collection := r.db.Collection(backupsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid backup ID: %w", err)
	}

	var backup model.BackupMetadata
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&backup)
	if err != nil {
		return nil, fmt.Errorf("failed to get backup: %w", err)
	}

	return &backup, nil
}

// DeleteBackup removes a backup record
func (r *Repository) DeleteBackup(ctx context.Context, id string) error {
	collection := r.db.Collection(backupsCollection)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid backup ID: %w", err)
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("backup not found")
	}

	return nil
}

// UpdateStatus updates the status and error message of the most recent backup for a given host/database/type
func (r *Repository) UpdateStatus(ctx context.Context, host, database, dbType string, status model.BackupStatus, errorMsg string) error {
	collection := r.db.Collection(backupsCollection)

	// Find the most recent backup for this host/database/type
	filter := bson.M{
		"host":     host,
		"database": database,
		"type":     dbType,
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
			"error":  errorMsg,
		},
	}

	// Sort by timestamp descending to get the most recent
	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var result model.BackupMetadata
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to update backup status: %w", err)
	}

	return nil
}

// UpdateMetadata updates the complete metadata of the most recent backup
func (r *Repository) UpdateMetadata(ctx context.Context, host, database, dbType, filePath, objectKey string, fileSize int64, status model.BackupStatus, errorMsg string) error {
	collection := r.db.Collection(backupsCollection)

	// Find the most recent backup for this host/database/type
	filter := bson.M{
		"host":     host,
		"database": database,
		"type":     dbType,
	}

	update := bson.M{
		"$set": bson.M{
			"filePath":  filePath,
			"objectKey": objectKey,
			"fileSize":  fileSize,
			"status":    status,
			"error":     errorMsg,
		},
	}

	// Sort by timestamp descending to get the most recent
	opts := options.FindOneAndUpdate().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var result model.BackupMetadata
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return fmt.Errorf("failed to update backup metadata: %w", err)
	}

	return nil
}

// BackupStats represents aggregated backup statistics
type BackupStats struct {
	Total    int64            `json:"total"`
	ByType   map[string]int64 `json:"byType"`
	ByStatus map[string]int64 `json:"byStatus"`
}

// GetBackupStats retrieves aggregated backup statistics
func (r *Repository) GetBackupStats(ctx context.Context, startDate, endDate *primitive.DateTime) (*BackupStats, error) {
	collection := r.db.Collection(backupsCollection)

	// Build filter for date range
	filter := bson.M{}
	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}
		if startDate != nil {
			dateFilter["$gte"] = *startDate
		}
		if endDate != nil {
			dateFilter["$lte"] = *endDate
		}
		filter["createdAt"] = dateFilter
	}

	// Get total count
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to count total backups: %w", err)
	}

	// Aggregate by type
	byTypePipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$type"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	byTypeResults, err := collection.Aggregate(ctx, byTypePipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate by type: %w", err)
	}
	defer byTypeResults.Close(ctx)

	byType := make(map[string]int64)
	for byTypeResults.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := byTypeResults.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode type result: %w", err)
		}
		byType[result.ID] = result.Count
	}

	// Aggregate by status
	byStatusPipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$status"},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
	}

	byStatusResults, err := collection.Aggregate(ctx, byStatusPipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to aggregate by status: %w", err)
	}
	defer byStatusResults.Close(ctx)

	byStatus := make(map[string]int64)
	for byStatusResults.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		}
		if err := byStatusResults.Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode status result: %w", err)
		}
		byStatus[result.ID] = result.Count
	}

	return &BackupStats{
		Total:    total,
		ByType:   byType,
		ByStatus: byStatus,
	}, nil
}
