package scheduler

import (
	"context"
	"db-backup/internal/database"
	"db-backup/internal/model"
	"db-backup/internal/worker"
	"log"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	cron gocron.Scheduler
	repo *database.Repository
	jobs map[string]uuid.UUID // Map database ID to job ID
}

func NewScheduler(repo *database.Repository) (*Scheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		cron: s,
		repo: repo,
		jobs: make(map[string]uuid.UUID),
	}, nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
	s.loadJobs()
}

func (s *Scheduler) Stop() {
	if err := s.cron.Shutdown(); err != nil {
		log.Printf("Error shutting down scheduler: %v", err)
	}
}

func (s *Scheduler) loadJobs() {
	ctx := context.Background()
	// Fetch all databases (using a large limit for now)
	dbs, _, err := s.repo.ListDatabases(ctx, 1, 1000)
	if err != nil {
		log.Printf("Failed to load databases for scheduler: %v", err)
		return
	}

	for _, db := range dbs {
		if db.IsActive && db.CronExpression != "" {
			if err := s.AddJob(&db); err != nil {
				log.Printf("Failed to schedule job for database %s: %v", db.Name, err)
			}
		}
	}
}

func (s *Scheduler) AddJob(db *model.Database) error {
	// Remove existing job if any
	s.RemoveJob(db.ID.Hex())

	if !db.IsActive || db.CronExpression == "" {
		return nil
	}

	job, err := s.cron.NewJob(
		gocron.CronJob(db.CronExpression, false),
		gocron.NewTask(
			func(db model.Database) {
				log.Printf("Running scheduled backup for %s", db.Name)
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
				worker.ProcessBackup(req)
			},
			*db,
		),
	)

	if err != nil {
		return err
	}

	s.jobs[db.ID.Hex()] = job.ID()
	log.Printf("Scheduled backup for %s with schedule %s", db.Name, db.CronExpression)
	return nil
}

func (s *Scheduler) RemoveJob(dbID string) {
	if jobID, exists := s.jobs[dbID]; exists {
		if err := s.cron.RemoveJob(jobID); err != nil {
			log.Printf("Failed to remove job for database %s: %v", dbID, err)
		}
		delete(s.jobs, dbID)
		log.Printf("Removed scheduled backup for database %s", dbID)
	}
}
