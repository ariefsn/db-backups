package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// slugify turns an arbitrary label into a filesystem and object-key safe
// segment: lowercase, non-alphanumerics collapsed into single dashes.
func slugify(s string) string {
	var b strings.Builder
	lastDash := true // leading dashes are trimmed by starting as if one was written

	for _, r := range strings.ToLower(s) {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			lastDash = false
		case !lastDash:
			b.WriteByte('-')
			lastDash = true
		}
	}

	slug := strings.Trim(b.String(), "-")
	if len(slug) > maxSlugLength {
		slug = strings.Trim(slug[:maxSlugLength], "-")
	}

	return slug
}

const maxSlugLength = 40

func generateFilename(req model.BackupRequest, ext string) string {
	timestamp := time.Now().Format("20060102_150405")
	dir := filepath.Join("backups", string(req.Type))
	ensureDir(dir)

	slug := slugify(req.DisplayName())
	if slug == "" {
		slug = string(req.Type)
	}

	// The backup ID suffix keeps two concurrent jobs from writing the same
	// path, which would otherwise overwrite each other in R2 as well.
	var idPart string
	if req.BackupID != "" {
		id := req.BackupID
		if len(id) > 8 {
			id = id[:8]
		}
		idPart = fmt.Sprintf("_%s", id)
	}

	return filepath.Join(dir, fmt.Sprintf("%s_%s%s.%s", slug, timestamp, idPart, ext))
}
