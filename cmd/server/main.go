package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "db-backup/docs" // Import generated docs
	"db-backup/internal/api"
	"db-backup/internal/database"
	"db-backup/internal/scheduler"
	"db-backup/internal/worker"

	"github.com/joho/godotenv"
)

// @title Database Backup API
// @version 1.0
// @description API for triggering database backups and managing them.
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	// Initialize MongoDB connection
	ctx := context.Background()
	if err := database.Connect(ctx); err != nil {
		log.Printf("Warning: MongoDB connection failed: %v", err)
		log.Println("Application will continue but backup metadata will not be saved")
	} else {
		log.Println("Connected to MongoDB successfully")
		defer func() {
			if err := database.Disconnect(ctx); err != nil {
				log.Printf("Error disconnecting from MongoDB: %v", err)
			}
		}()
	}

	// Initialize worker dependencies
	if err := worker.InitializeWorker(); err != nil {
		log.Printf("Warning: Worker initialization failed: %v", err)
	}

	// Initialize API handlers
	api.InitializeHandlers()

	// Initialize Scheduler
	schedulerRepo := database.NewRepository()
	jobScheduler, err := scheduler.NewScheduler(schedulerRepo)
	if err != nil {
		log.Printf("Warning: Failed to create scheduler: %v", err)
	} else {
		jobScheduler.Start()
		defer jobScheduler.Stop()
		api.SetScheduler(jobScheduler)
	}

	// Start cleanup cron
	// Disabled, since we have capabilities to delete backups
	// cron.StartCleanupCron()

	router := api.NewRouter()

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
