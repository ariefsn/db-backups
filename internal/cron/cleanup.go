package cron

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robfig/cron/v3"
)

func StartCleanupCron() {
	c := cron.New()

	// Run daily at midnight
	_, err := c.AddFunc("@daily", func() {
		log.Println("Running daily cleanup...")
		if err := CleanupOldBackups(); err != nil {
			log.Printf("Cleanup failed: %v", err)
		} else {
			log.Println("Cleanup completed successfully.")
		}
	})

	if err != nil {
		log.Printf("Failed to add cron job: %v", err)
		return
	}

	c.Start()
	log.Println("Cleanup cron started.")
}

func CleanupOldBackups() error {
	backupDir := "backups"
	retentionDays := 7

	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			age := time.Since(info.ModTime())
			if age.Hours() > float64(retentionDays*24) {
				log.Printf("Deleting old backup: %s (Age: %.2f hours)", path, age.Hours())
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
