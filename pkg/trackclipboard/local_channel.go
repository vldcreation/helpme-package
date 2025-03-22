package trackclipboard

import (
	"context"
	"os"
	"path/filepath"
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
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(l.Path), 0755); err != nil {
		return err
	}

	// Construct full file path
	fullPath := filepath.Join(l.Path, l.Name)

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
