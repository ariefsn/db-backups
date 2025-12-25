package worker

import (
	"bytes"
	"context"
	"db-backup/internal/backup"
	"db-backup/internal/database"
	"db-backup/internal/model"
	"db-backup/internal/storage"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	storageClient *storage.Client
	backupRepo    *database.Repository
)

// InitializeWorker initializes the worker dependencies
func InitializeWorker() error {
	var err error

	// Initialize storage client
	storageClient, err = storage.NewClient()
	if err != nil {
		log.Printf("Warning: Failed to initialize storage client: %v", err)
		// Don't fail if storage is not configured
	}

	// Initialize backup repository
	backupRepo = database.NewRepository()

	return nil
}

func ProcessBackup(req model.BackupRequest) string {
	timestamp := time.Now()

	// Create a context for initial save
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Save initial pending status and get the ID
	var backupID string
	if backupRepo != nil {
		backupID = saveBackupMetadata(ctx, req, "", "", 0, model.StatusPending, "", timestamp)
	}
	cancel()

	// Run the actual backup in a goroutine
	go func() {
		// Create a context with timeout, e.g. 1 hour max for backup
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
		defer cancel()

		log.Printf("Starting backup for %s (%s)", req.Type, req.Host)

		strategy, err := backup.NewStrategy(req.Type)
		if err != nil {
			notifyWebhook(req.WebhookURL, model.BackupResult{
				Success: false,
				Error:   err.Error(),
			})
			log.Printf("Failed to create strategy: %v", err)

			// Update to failed status
			if backupRepo != nil && backupID != "" {
				updateBackupStatus(ctx, backupID, model.StatusFailed, err.Error())
			}
			return
		}

		// Update to generating status
		if backupRepo != nil && backupID != "" {
			updateBackupStatus(ctx, backupID, model.StatusGenerating, "")
		}

		filePath, err := strategy.Backup(ctx, req)
		result := model.BackupResult{
			Success:   err == nil,
			FilePath:  filePath,
			Timestamp: timestamp.Format(time.RFC3339),
			Metadata:  make(map[string]string),
		}

		if err != nil {
			result.Error = err.Error()
			log.Printf("Backup failed for %s: %v", req.Type, err)

			// Update to failed status
			if backupRepo != nil && backupID != "" {
				updateBackupStatus(ctx, backupID, model.StatusFailed, err.Error())
			}
		} else {
			log.Printf("Backup completed for %s: %s", req.Type, filePath)

			// Get file size
			fileInfo, _ := os.Stat(filePath)
			var fileSize int64
			if fileInfo != nil {
				fileSize = fileInfo.Size()
			}

			// Upload to R2 if configured
			var objectKey string
			if storageClient != nil {
				objectKey, err = storageClient.Upload(ctx, filePath, storage.UploadMetadata{
					DatabaseType: string(req.Type),
					Host:         req.Host,
					Database:     req.Database,
					Timestamp:    timestamp,
					FileSize:     fileSize,
				})

				if err != nil {
					log.Printf("Failed to upload to R2: %v", err)
					result.Metadata["upload_error"] = err.Error()
				} else {
					result.ObjectKey = objectKey
					result.Metadata["storage"] = "r2"
					log.Printf("Uploaded to R2: %s", objectKey)
				}
			}

			// Update backup metadata to completed
			if backupRepo != nil && backupID != "" {
				updateBackupMetadata(ctx, backupID, filePath, objectKey, fileSize, model.StatusCompleted, "")
			}

			// Add metadata
			result.Metadata["database_type"] = string(req.Type)
			result.Metadata["host"] = req.Host
			result.Metadata["database"] = req.Database
			result.Metadata["file_size"] = string(rune(fileSize))
		}

		notifyWebhook(req.WebhookURL, result)
	}()

	return backupID
}

func saveBackupMetadata(ctx context.Context, req model.BackupRequest, filePath, objectKey string, fileSize int64, status model.BackupStatus, errorMsg string, timestamp time.Time) string {
	metadata := &model.BackupMetadata{
		Type:      string(req.Type),
		ObjectKey: objectKey,
		FilePath:  filePath,
		FileSize:  fileSize,
		Timestamp: timestamp,
		Status:    status,
		Error:     errorMsg,
		Host:      req.Host,
		Database:  req.Database,
		CreatedAt: primitive.NewDateTimeFromTime(timestamp),
	}

	if err := backupRepo.SaveBackup(ctx, metadata); err != nil {
		log.Printf("Failed to save backup metadata: %v", err)
		return ""
	} else {
		log.Printf("Saved backup metadata to MongoDB: %s", metadata.ID.Hex())
		return metadata.ID.Hex()
	}
}

func updateBackupStatus(ctx context.Context, id string, status model.BackupStatus, errorMsg string) {
	if err := backupRepo.UpdateBackupStatusByID(ctx, id, status, errorMsg); err != nil {
		log.Printf("Failed to update backup status: %v", err)
	}
}

func updateBackupMetadata(ctx context.Context, id, filePath, objectKey string, fileSize int64, status model.BackupStatus, errorMsg string) {
	if err := backupRepo.UpdateBackupMetadataByID(ctx, id, filePath, objectKey, fileSize, status, errorMsg); err != nil {
		log.Printf("Failed to update backup metadata: %v", err)
	}
}

func notifyWebhook(url string, result model.BackupResult) {
	if url == "" {
		return
	}

	body, _ := json.Marshal(result)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to notify webhook: %v", err)
		return
	}
	defer resp.Body.Close()
}
