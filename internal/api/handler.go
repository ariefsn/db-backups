package api

import (
	"context"
	"db-backup/internal/database"
	"db-backup/internal/model"
	"db-backup/internal/scheduler"
	"db-backup/internal/storage"
	"db-backup/internal/worker"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	backupRepo    *database.Repository
	storageClient *storage.Client
	jobScheduler  *scheduler.Scheduler
)

// SetScheduler sets the scheduler instance for the API package
func SetScheduler(s *scheduler.Scheduler) {
	jobScheduler = s
}

// InitializeHandlers initializes the API handler dependencies
func InitializeHandlers() {
	backupRepo = database.NewRepository()

	// Initialize storage client (optional)
	var err error
	storageClient, err = storage.NewClient()
	if err != nil {
		log.Printf("Warning: Storage client not initialized: %v", err)
	}
}

// HandleBackup godoc
// @Summary Trigger a database backup
// @Description Queue a backup job for the specified database
// @Tags backup
// @Accept json
// @Produce json
// @Param request body model.BackupRequest true "Backup Request"
// @Success 202 {object} model.BackupResponse "Backup started"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Router /backup [post]
func HandleBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.BackupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if req.Type == "" || (req.Host == "" && req.ConnectionURI == "") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Missing required fields",
			Error:   "type and (host or connectionUri) are required",
		})
		return
	}

	backupID := worker.ProcessBackup(req)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(model.BackupResponse{
		Success: true,
		Message: "Backup job submitted successfully",
		ID:      backupID,
	})
}

// HandleListBackups godoc
// @Summary List all backups
// @Description Get a paginated list of all backups with optional filtering
// @Tags backup
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param statuses query string false "Comma-separated status values (pending,generating,completed,failed)"
// @Param search query string false "Search keyword (searches in database, host, type)"
// @Param orderBy query string false "Field to order by" default(createdAt)
// @Param orderDir query string false "Order direction (asc/desc)" default(desc)
// @Param startDate query string false "Start date for filtering (RFC3339 format)"
// @Param endDate query string false "End date for filtering (RFC3339 format)"
// @Success 200 {object} model.BackupListResponse "List of backups"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /backups [get]
func HandleListBackups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	page := 1
	limit := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Parse status filter
	var statuses []model.BackupStatus
	if statusesStr := r.URL.Query().Get("statuses"); statusesStr != "" {
		statusStrings := strings.Split(statusesStr, ",")
		for _, s := range statusStrings {
			s = strings.TrimSpace(s)
			if s != "" {
				statuses = append(statuses, model.BackupStatus(s))
			}
		}
	}

	// Parse types filter
	var types []string
	if typesStr := r.URL.Query().Get("types"); typesStr != "" {
		typeStrings := strings.Split(typesStr, ",")
		for _, t := range typeStrings {
			t = strings.TrimSpace(t)
			if t != "" {
				types = append(types, t)
			}
		}
	}

	// Parse search, orderBy, orderDir
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	orderDir := r.URL.Query().Get("orderDir")

	// Parse date range
	var startDate, endDate *primitive.DateTime
	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			dt := primitive.NewDateTimeFromTime(t)
			startDate = &dt
		}
	}
	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			dt := primitive.NewDateTimeFromTime(t)
			endDate = &dt
		}
	}

	// Get backups from repository
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	backups, total, err := backupRepo.ListBackups(ctx, page, limit, statuses, types, search, orderBy, orderDir, startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to list backups",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.BackupListResponse{
		Backups: backups,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}

// HandleGetBackup godoc
// @Summary Get a single backup
// @Description Retrieve a single backup by ID
// @Tags backup
// @Produce json
// @Param id path string true "Backup ID"
// @Success 200 {object} model.BackupMetadata "Backup details"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Backup not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /backups/{id} [get]
func HandleGetBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	backupID := chi.URLParam(r, "id")
	if backupID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	backup, err := backupRepo.GetBackup(ctx, backupID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup not found",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(backup)
}

// HandleDeleteBackup godoc
// @Summary Delete a backup
// @Description Delete a backup from both MongoDB and R2 storage
// @Tags backup
// @Produce json
// @Param id path string true "Backup ID"
// @Success 200 {object} model.BackupResponse "Backup deleted successfully"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Backup not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /backups/{id} [delete]
func HandleDeleteBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	backupID := chi.URLParam(r, "id")
	if backupID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Get backup metadata first
	backup, err := backupRepo.GetBackup(ctx, backupID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup not found",
			Error:   err.Error(),
		})
		return
	}

	// Delete from R2 if object key exists
	if backup.ObjectKey != "" && storageClient != nil {
		if err := storageClient.Delete(ctx, backup.ObjectKey); err != nil {
			log.Printf("Failed to delete from R2: %v", err)
			// Continue with MongoDB deletion even if R2 deletion fails
		}
	}

	// Delete from MongoDB
	if err := backupRepo.DeleteBackup(ctx, backupID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to delete backup",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.BackupResponse{
		Success: true,
		Message: "Backup deleted successfully",
	})
}

// HandleDownloadBackup godoc
// @Summary Download a backup file
// @Description Generate a presigned URL to download a backup file from R2 storage
// @Tags backup
// @Produce json
// @Param id path string true "Backup ID"
// @Success 200 {object} map[string]string "Download URL"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Backup not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /backups/{id}/download [get]
func HandleDownloadBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	backupID := chi.URLParam(r, "id")
	if backupID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Get backup metadata
	backup, err := backupRepo.GetBackup(ctx, backupID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup not found",
			Error:   err.Error(),
		})
		return
	}

	// Check if storage client is available
	if storageClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Storage client not available",
		})
		return
	}

	// Check if backup has an object key
	if backup.ObjectKey == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Backup file not found in storage",
		})
		return
	}

	// Generate presigned URL (valid for 1 hour)
	url, err := storageClient.GetPresignedURL(ctx, backup.ObjectKey, 1*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to generate download URL",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"url": url,
	})
}

// HandleGetBackupStats godoc
// @Summary Get backup statistics
// @Description Retrieve aggregated backup statistics by type and status
// @Tags backup
// @Produce json
// @Param startDate query string false "Start date for filtering (RFC3339 format)"
// @Param endDate query string false "End date for filtering (RFC3339 format)"
// @Success 200 {object} database.BackupStats "Backup statistics"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /backups/stats [get]
func HandleGetBackupStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse date range
	var startDate, endDate *primitive.DateTime
	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			dt := primitive.NewDateTimeFromTime(t)
			startDate = &dt
		}
	}
	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			dt := primitive.NewDateTimeFromTime(t)
			endDate = &dt
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	stats, err := backupRepo.GetBackupStats(ctx, startDate, endDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to get backup statistics",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
