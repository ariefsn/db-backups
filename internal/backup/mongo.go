package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os/exec"
)

type MongoBackup struct{}

func (b *MongoBackup) Backup(ctx context.Context, req model.BackupRequest) (string, error) {
	// mongodump creates a directory usually, or an archive. Archive is better for single file.
	filename := generateFilename(req, "gz")

	var args []string
	if req.ConnectionURI != "" {
		args = append(args, fmt.Sprintf("--uri=%s", req.ConnectionURI))
	} else {
		args = append(args,
			fmt.Sprintf("--host=%s", req.Host),
			fmt.Sprintf("--port=%s", req.Port),
		)
		if req.Username != "" {
			args = append(args, fmt.Sprintf("--username=%s", req.Username))
		}
		if req.Password != "" {
			args = append(args, fmt.Sprintf("--password=%s", req.Password))
		}
		if req.Database != "" {
			args = append(args, fmt.Sprintf("--db=%s", req.Database))
		}
	}

	args = append(args, fmt.Sprintf("--archive=%s", filename), "--gzip")

	binPath := resolveExecutable("mongodump")
	cmd := exec.CommandContext(ctx, binPath, args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mongodump failed: %s, output: %s", err, string(output))
	}

	return filename, nil
}
