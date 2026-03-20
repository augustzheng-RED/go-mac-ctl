package observe

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"go-mac-ctl/internal/executor"
)

const (
	screenshotDir = "tmp/screenshots"
	maxScreenshots = 20
)

func Screenshot() (string, error) {
	if err := os.MkdirAll(screenshotDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	filename := filepath.Join(screenshotDir, fmt.Sprintf("screenshot_%d.png", time.Now().UnixNano()))
	if err := executor.Screenshot(filename); err != nil {
		return "", err
	}

	if err := trimOldScreenshots(); err != nil {
		return "", err
	}

	return filename, nil
}

func trimOldScreenshots() error {
	entries, err := os.ReadDir(screenshotDir)
	if err != nil {
		return fmt.Errorf("failed to read screenshot directory: %w", err)
	}

	type screenshotFile struct {
		path    string
		modTime time.Time
	}

	files := make([]screenshotFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("failed to stat screenshot file: %w", err)
		}

		files = append(files, screenshotFile{
			path:    filepath.Join(screenshotDir, entry.Name()),
			modTime: info.ModTime(),
		})
	}

	if len(files) <= maxScreenshots {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	for _, file := range files[:len(files)-maxScreenshots] {
		if err := os.Remove(file.path); err != nil {
			return fmt.Errorf("failed to remove old screenshot %q: %w", file.path, err)
		}
	}

	return nil
}
