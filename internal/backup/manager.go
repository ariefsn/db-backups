package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Strategy interface {
	Backup(ctx context.Context, req model.BackupRequest) (string, error)
}

func NewStrategy(t model.BackupType) (Strategy, error) {
	switch t {
	case model.Postgres:
		return &PostgresBackup{}, nil
	case model.MySQL:
		return &MySQLBackup{}, nil
	case model.Mongo:
		return &MongoBackup{}, nil
	case model.Redis:
		return &RedisBackup{}, nil
	default:
		return nil, fmt.Errorf("unsupported backup type: %s", t)
	}
}

func resolveExecutable(binName string) string {
	path, err := exec.LookPath(binName)
	if err == nil {
		return path
	}

	commonPaths := []string{
		"/opt/homebrew/bin",
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
	}

	for _, dir := range commonPaths {
		fullPath := filepath.Join(dir, binName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return binName
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func generateFilename(req model.BackupRequest, ext string) string {
	timestamp := time.Now().Format("20060102_150405")
	dir := filepath.Join("backups", string(req.Type))
	ensureDir(dir)
	var dbPart string
	if req.Database != "" {
		dbPart = fmt.Sprintf("_%s", req.Database)
	}
	return filepath.Join(dir, fmt.Sprintf("%s_%s%s.%s", req.Host, timestamp, dbPart, ext))
}
