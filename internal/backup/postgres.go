package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os"
	"os/exec"
)

type PostgresBackup struct{}

func (b *PostgresBackup) Backup(ctx context.Context, req model.BackupRequest) (string, error) {
	filename := generateFilename(req, "sql")

	binPath := resolveExecutable("pg_dump")
	// PGPASSWORD environment variable is the safest way to pass password to pg_dump
	cmd := exec.CommandContext(ctx, binPath,
		"-h", req.Host,
		"-p", req.Port,
		"-U", req.Username,
		"-f", filename,
		req.Database,
	)

	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", req.Password))

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("pg_dump failed: %s, output: %s", err, string(output))
	}

	return filename, nil
}
