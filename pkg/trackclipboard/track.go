package trackclipboard

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

type TrackClipboard struct {
	Cfg     *Config
	Channel TrackChannel
}

func NewTrackClipboard(cfg *Config) *TrackClipboard {
	t := &TrackClipboard{
		Cfg: cfg,
	}

	if t.Cfg.App == nil {
		t.Cfg.App = &APPConfig{}
	}

	if t.Cfg.App.Channel == "" {
		t.Cfg.App.Channel = "local"
	}

	if t.Cfg.App.Idle <= 0 {
		t.Cfg.App.Idle = 10 * time.Second
	}

	if t.Cfg.App.Channel == "local" {
		if t.Cfg.File == nil {
			t.Cfg.File = &FileConfig{}
		}

		// if user did not provide a fpath and fname, use a default one
		if t.Cfg.File.Path == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}
			t.Cfg.File = &FileConfig{
				Path: filepath.Join(homeDir, "Downloads"),
				Name: "",
			}
		}
		if t.Cfg.File.Name == "" {
			var s strings.Builder
			s.WriteString("ressource-")
			s.WriteString(time.Now().Format("2006-01-02-15-04-05"))
			s.WriteString(".txt")
			t.Cfg.File.Name = s.String()
		}
	}

	if t.Cfg.App.Channel == "telegram" {
		if t.Cfg.Telegram == nil {
			panic("telegram config is nil")
		}
	}

	switch t.Cfg.App.Channel {
	case "local":
		t.Channel = NewLocalChannel(t.Cfg.File)
	case "telegram":
		t.Channel = NewTelegramChannel(t.Cfg.Telegram)
	}

	return t
}

func (t *TrackClipboard) Track() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	// Create a parent context for overall control
	parentCtx := context.Background()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(parentCtx)
	defer func() {
		// Cleanup resources
		cancel()
		if err := t.Channel.Close(); err != nil {
			fmt.Printf("Error closing channel: %v\n", err)
		}
	}()

	// Create timer for idle timeout
	timer := time.NewTimer(t.Cfg.App.Idle)
	defer timer.Stop()

	// Watch clipboard changes
	content := clipboard.Watch(ctx, clipboard.FmtText)

	for {
		select {
		case data := <-content:
			// Reset timer when clipboard content changes
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			if err := t.Channel.Send(ctx, string(data)); err != nil {
				fmt.Printf("Error sending message: %v\n", err)
			}
			timer.Reset(t.Cfg.App.Idle)
			if t.Cfg.App.Debug {
				fmt.Printf("%s: Sucess sending clipboard content => %s\n", t.Cfg.App.Channel, string(data))
			}
		case <-ctx.Done():
			fmt.Println("Context canceled, stopping clipboard tracking")
			return
		case <-timer.C:
			fmt.Println("Idle timeout, stopping clipboard tracking")
			return
		}
	}
}
