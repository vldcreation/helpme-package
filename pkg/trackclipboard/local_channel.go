package trackclipboard

import (
	"context"
	"os"
	"path/filepath"
	"strings"
)

type LocalChannel struct {
	Path string
	Name string
}

func NewLocalChannel(cfg *FileConfig) TrackChannel {
	return &LocalChannel{
		Path: cfg.Path,
		Name: cfg.Name,
	}
}

func (l *LocalChannel) Send(ctx context.Context, msg string) error {
	// Expand tilde to home directory if present
	path := l.Path
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Ensure directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// Construct full file path
	fullPath := filepath.Join(path, l.Name)

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(msg); err != nil {
		return err
	}
	return nil
}
