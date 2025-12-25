package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os"
	"os/exec"
)

type RedisBackup struct{}

func (b *RedisBackup) Backup(ctx context.Context, req model.BackupRequest) (string, error) {
	// Redis backup is tricky remotely without just triggering SAVE and downloading dump.rdb.
	// However, `redis-cli --rdb filename` is a standard way to do remote backup.

	filename := generateFilename(req, "rdb")

	binPath := resolveExecutable("redis-cli")
	cmd := exec.CommandContext(ctx, binPath,
		"-h", req.Host,
		"-p", req.Port,
		"-a", req.Password, // Warning: password on CLI
		"--rdb", filename,
	)

	// redis-cli might output warning about password on CLI, but it's the standard CLI flag.

	if output, err := cmd.CombinedOutput(); err != nil {
		// If redis-cli fails, we check output.
		return "", fmt.Errorf("redis-cli failed: %s, output: %s", err, string(output))
	}

	// Verify file exists and is not empty
	info, err := os.Stat(filename)
	if err != nil {
		return "", fmt.Errorf("failed to verify backup file: %w", err)
	}
	if info.Size() == 0 {
		os.Remove(filename)
		return "", fmt.Errorf("backup file is empty")
	}

	return filename, nil
}
