package worker

import (
	"context"
	"db-backup/internal/model"
	"log"
	"os"
	"strconv"
)

const defaultMaxBackupsPerDatabase = 10

// MaxBackupsPerDatabase returns how many completed backups to keep per
// database. A value of 0 (or a negative/invalid one) disables retention.
func MaxBackupsPerDatabase() int {
	raw := os.Getenv("MAX_BACKUPS_PER_DATABASE")
	if raw == "" {
		return defaultMaxBackupsPerDatabase
	}

	limit, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("Invalid MAX_BACKUPS_PER_DATABASE %q, falling back to %d", raw, defaultMaxBackupsPerDatabase)
		return defaultMaxBackupsPerDatabase
	}

	if limit < 0 {
		return 0
	}

	return limit
}

// RemoveLocalFile deletes a local backup file. It reports whether the file is
// gone afterwards, so an already-missing file counts as success.
func RemoveLocalFile(path string) bool {
	if path == "" {
		return true
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return true
		}
		log.Printf("Failed to delete local backup file %s: %v", path, err)
		return false
	}

	log.Printf("Deleted local backup file: %s", path)
	return true
}

// pruneOldBackups keeps only the newest N completed backups for the database
// this request belongs to, deleting the R2 object, the local file and the
// metadata record of everything older.
func pruneOldBackups(ctx context.Context, req model.BackupRequest) {
	limit := MaxBackupsPerDatabase()
	if limit <= 0 || backupRepo == nil {
		return
	}

	backups, err := backupRepo.ListBackupsForRetention(ctx, req.DatabaseID, string(req.Type), req.Host, req.Database)
	if err != nil {
		log.Printf("Retention: failed to list backups: %v", err)
		return
	}

	if len(backups) <= limit {
		return
	}

	for _, backup := range backups[limit:] {
		id := backup.ID.Hex()

		if backup.ObjectKey != "" && storageClient != nil {
			if err := storageClient.Delete(ctx, backup.ObjectKey); err != nil {
				log.Printf("Retention: failed to delete %s from R2: %v", backup.ObjectKey, err)
				// Continue anyway so the record does not linger forever.
			}
		}

		RemoveLocalFile(backup.FilePath)

		if err := backupRepo.DeleteBackup(ctx, id); err != nil {
			log.Printf("Retention: failed to delete backup record %s: %v", id, err)
			continue
		}

		log.Printf("Retention: removed old backup %s (%s)", id, backup.ObjectKey)
	}
}
