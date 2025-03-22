package trackclipboard

import (
	"context"
	"os"
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
	f, err := os.OpenFile(l.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(msg); err != nil {
		return err
	}
	return nil
}
