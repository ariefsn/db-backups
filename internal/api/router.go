package api

import (
	"net/http"

	_ "db-backup/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Backup endpoints
	// Backup endpoints
	r.Post("/backup", HandleBackup)
	r.Get("/backups/stats", HandleGetBackupStats)
	r.Get("/backups", HandleListBackups)
	r.Get("/backups/{id}", HandleGetBackup)
	r.Get("/backups/{id}/download", HandleDownloadBackup)
	r.Delete("/backups/{id}", HandleDeleteBackup)

	// Database endpoints
	r.Get("/databases", HandleListDatabases)
	r.Post("/databases", HandleCreateDatabase)
	r.Get("/databases/{id}", HandleGetDatabase)
	r.Put("/databases/{id}", HandleUpdateDatabase)
	r.Delete("/databases/{id}", HandleDeleteDatabase)
	r.Post("/databases/{id}/backup", HandleTriggerBackup)

	return r
}
