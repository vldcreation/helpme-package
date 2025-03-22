package trackclipboard

import (
	"context"
	"fmt"
	"os"
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

	if t.Cfg.Channel == "" {
		t.Cfg.Channel = "local"
	}

	if t.Cfg.Idle <= 0 {
		t.Cfg.Idle = 10 * time.Second
	}

	if t.Cfg.Channel == "local" {
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
				Path: homeDir + "/Downloads/",
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

	if t.Cfg.Channel == "telegram" {
		if t.Cfg.Telegram == nil {
			panic("telegram config is nil")
		}
	}

	switch t.Cfg.Channel {
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
	defer cancel()

	// Create timer for idle timeout
	timer := time.NewTimer(t.Cfg.Idle)
	defer timer.Stop()

	// Watch clipboard changes
	content := clipboard.Watch(ctx, clipboard.FmtText)

	for {
		select {
		case data := <-content:
			// Reset timer when clipboard content changes
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(t.Cfg.Idle)
			t.Channel.Send(ctx, string(data))
		case <-timer.C:
			fmt.Println("idle timeout")
			return
		}
	}
}

func writeToFile(fpath string, message string) {
	// create temporary file
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(file *os.File, cause error) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", cause)
		}
	}(file, err)

	_, err = file.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
