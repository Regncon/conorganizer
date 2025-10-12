package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Regncon/conorganizer/backup-service/models"
)

type Flyctl struct {
	Config models.Config
	Logger *slog.Logger
}

func FlyCommands(cfg models.Config, logger *slog.Logger) (*Flyctl, error) {
	var flyctl = &Flyctl{
		Config: cfg,
		Logger: logger,
	}

	if err := flyctl.verifyFlyctlPath(); err != nil {
		return nil, err
	}

	return flyctl, nil
}

func (flyctl *Flyctl) verifyFlyctlPath() error {
	if _, err := exec.LookPath("flyctl"); err != nil {
		return fmt.Errorf("flyctl not found in PATH: %w", err)
	}
	return nil
}

func (flyctl *Flyctl) exec(ctx context.Context, cmdString string) (string, error) {
	// Check if flyctl is installed and accessable in path
	var cmdArr = strings.Split(cmdString, " ")

	output, err := exec.CommandContext(ctx, cmdArr[0], cmdArr[1:]...).CombinedOutput()
	if err != nil {
		if output != nil {
			return "", fmt.Errorf("failed to execute flyctl commands: %w\n%s", err, string(output))
		}
		return "", fmt.Errorf("flyclt failed, unable to read output: %w\n%s", err, string(output))
	}
	return string(output), nil
}

// DownloadDatabaseFromVolume todo: pass ctx
func (flyctl *Flyctl) DownloadDatabaseFromVolume(ctx context.Context) (string, error) {
	// Construct tmp folder name based on date
	var tmpDirName = time.Now().Format("20060102T150405")
	var tmpDir = filepath.Join(os.TempDir(), tmpDirName)
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		return "", fmt.Errorf("unable to create tmp dir: %w", err)
	}

	// Iterate over these files
	var targetFiles = []string{"events.db", "events.db-wal", "events.db-shm"}
	for _, fileName := range targetFiles {
		var cmdString = createDownloadQuery(fileName, tmpDir)
		fmt.Println(cmdString)
		if _, err := flyctl.exec(ctx, cmdString); err != nil {
			return "", err
		}
	}
	return tmpDir, nil
}

func createDownloadQuery(fileName string, targetDir string) string {
	return fmt.Sprintf(`flyctl -a regncon ssh sftp get /data/regncon/database/%s %s/%s`, fileName, targetDir, fileName)
}
