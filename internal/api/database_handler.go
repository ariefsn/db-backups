package api

import (
	"context"
	"db-backup/internal/model"
	"db-backup/internal/worker"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// HandleListDatabases godoc
// @Summary List all databases
// @Description List all saved database configurations
// @Tags database
// @Produce json
// @Success 200 {object} model.DatabaseListResponse
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Router /databases [get]
func HandleListDatabases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse pagination params
	page := 1
	limit := 10
	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	databases, total, err := backupRepo.ListDatabases(ctx, page, limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to list databases",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.DatabaseListResponse{
		Databases: databases,
		Total:     total,
		Page:      page,
		Limit:     limit,
	})
}

// HandleCreateDatabase godoc
// @Summary Create a new database
// @Description Save a new database configuration
// @Tags database
// @Accept json
// @Produce json
// @Param request body model.CreateDatabaseRequest true "Database Configuration"
// @Success 201 {object} model.Database
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /databases [post]
func HandleCreateDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req model.CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	db := &model.Database{
		Name:           req.Name,
		Type:           req.Type,
		Host:           req.Host,
		Port:           req.Port,
		Username:       req.Username,
		Password:       req.Password,
		Database:       req.Database,
		ConnectionURI:  req.ConnectionURI,
		CronExpression: req.CronExpression,
		IsActive:       req.IsActive,
		WebhookURL:     req.WebhookURL,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := backupRepo.SaveDatabase(ctx, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to save database",
			Error:   err.Error(),
		})
		return
	}

	// Schedule job if active and has cron
	if jobScheduler != nil && db.IsActive && db.CronExpression != "" {
		jobScheduler.AddJob(db)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(db)
}

// HandleGetDatabase godoc
// @Summary Get a database
// @Description Retrieve a database configuration by ID
// @Tags database
// @Produce json
// @Param id path string true "Database ID"
// @Success 200 {object} model.Database
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Database not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /databases/{id} [get]
func HandleGetDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	db, err := backupRepo.GetDatabase(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database not found",
			Error:   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db)
}

// HandleUpdateDatabase godoc
// @Summary Update a database
// @Description Update an existing database configuration
// @Tags database
// @Accept json
// @Produce json
// @Param id path string true "Database ID"
// @Param request body model.UpdateDatabaseRequest true "Database Configuration"
// @Success 200 {object} model.Database
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Database not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /databases/{id} [put]
func HandleUpdateDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database ID is required",
		})
		return
	}

	var req model.UpdateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Check if exists
	db, err := backupRepo.GetDatabase(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database not found",
			Error:   err.Error(),
		})
		return
	}

	// Update fields
	db.Name = req.Name
	db.Type = req.Type
	db.Host = req.Host
	db.Port = req.Port
	db.Username = req.Username
	db.Password = req.Password
	db.Database = req.Database
	db.ConnectionURI = req.ConnectionURI
	db.CronExpression = req.CronExpression
	db.IsActive = req.IsActive
	db.WebhookURL = req.WebhookURL

	if err := backupRepo.UpdateDatabase(ctx, db); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to update database",
			Error:   err.Error(),
		})
		return
	}

	// Update scheduler
	if jobScheduler != nil {
		if db.IsActive && db.CronExpression != "" {
			jobScheduler.AddJob(db)
		} else {
			jobScheduler.RemoveJob(db.ID.Hex())
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db)
}

// HandleDeleteDatabase godoc
// @Summary Delete a database
// @Description Delete a database configuration
// @Tags database
// @Produce json
// @Param id path string true "Database ID"
// @Success 200 {object} model.BackupResponse "Database deleted successfully"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Database not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /databases/{id} [delete]
func HandleDeleteDatabase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := backupRepo.DeleteDatabase(ctx, id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Failed to delete database",
			Error:   err.Error(),
		})
		return
	}

	// Remove from scheduler
	if jobScheduler != nil {
		jobScheduler.RemoveJob(id)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(model.BackupResponse{
		Success: true,
		Message: "Database deleted successfully",
	})
}

// HandleTriggerBackup godoc
// @Summary Trigger a backup for a database
// @Description Manually trigger a backup for a saved database configuration
// @Tags database
// @Produce json
// @Param id path string true "Database ID"
// @Success 202 {object} model.BackupResponse "Backup job submitted successfully"
// @Failure 400 {object} model.BackupResponse "error: Bad request"
// @Failure 404 {object} model.BackupResponse "error: Database not found"
// @Failure 500 {object} model.BackupResponse "error: Internal server error"
// @Router /databases/{id}/backup [post]
func HandleTriggerBackup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Get database
	db, err := backupRepo.GetDatabase(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(model.BackupResponse{
			Success: false,
			Message: "Database not found",
			Error:   err.Error(),
		})
		return
	}

	// Create backup request
	req := model.BackupRequest{
		Type:          db.Type,
		Host:          db.Host,
		Port:          db.Port,
		Username:      db.Username,
		Password:      db.Password,
		Database:      db.Database,
		ConnectionURI: db.ConnectionURI,
		WebhookURL:    db.WebhookURL,
	}

	// Trigger backup
	backupID := worker.ProcessBackup(req)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(model.BackupResponse{
		Success: true,
		Message: "Backup job submitted successfully",
		ID:      backupID,
	})
}
