package backup

import (
	"context"
	"db-backup/internal/model"
	"fmt"
	"os"
	"os/exec"
)

type MySQLBackup struct{}

func (b *MySQLBackup) Backup(ctx context.Context, req model.BackupRequest) (string, error) {
	filename := generateFilename(req, "sql")

	// mysqldump -u [username] -p[password] [database_name] > [filename]
	// Note: Putting password in command line is insecure but common for simple tools.
	// A better way is using a config file, but environment variables for mysqldump
	// vary by version (MYSQL_PWD is deprecated/insecure often warning).
	// We will try using MYSQL_PWD env var to avoid command line arg if possible,
	// or fallback. Here we use the env var approach which is generally supported.

	binPath := resolveExecutable("mysqldump")
	cmd := exec.CommandContext(ctx, binPath,
		"-h", req.Host,
		"-P", req.Port,
		"-u", req.Username,
		req.Database,
	)

	// Redirect stdout to file
	// cmd.Stdout is not set, we will write to file manually or let shell handle it?
	// exec.Command doesn't support > redirection directly.
	// We should write stdout to the file.

	outfile, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Env = append(os.Environ(), fmt.Sprintf("MYSQL_PWD=%s", req.Password))

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("mysqldump failed: %w", err)
	}

	return filename, nil
}
