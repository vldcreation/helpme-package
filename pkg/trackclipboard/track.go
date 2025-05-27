package trackclipboard

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"errors"
	"strings"
	"time"

	"golang.design/x/clipboard"
)

// ChannelFactory defines the signature for functions that create TrackChannel instances.
type ChannelFactory func(config *Config) (TrackChannel, error)

// channelFactories is a registry for TrackChannel factory functions.
var channelFactories = make(map[string]ChannelFactory)

// newLocalChannelFactory creates a LocalChannel.
func newLocalChannelFactory(config *Config) (TrackChannel, error) {
	// Assumes config.File has been defaulted by getDefaultedFileConfig
	return NewLocalChannel(config.File), nil
}

// newTelegramChannelFactory creates a TelegramChannel.
func newTelegramChannelFactory(config *Config) (TrackChannel, error) {
	if config.Telegram == nil {
		return nil, errors.New("telegram config is missing for telegram channel")
	}
	return NewTelegramChannel(config.Telegram), nil
}

func init() {
	channelFactories["local"] = newLocalChannelFactory
	channelFactories["telegram"] = newTelegramChannelFactory
}

type TrackClipboard struct {
	Cfg     *Config
	Channel TrackChannel
}

func getDefaultedAPPConfig(cfg *APPConfig) *APPConfig {
	if cfg == nil {
		cfg = &APPConfig{}
	}
	if cfg.Channel == "" {
		cfg.Channel = "local"
	}
	if cfg.Idle <= 0 {
		cfg.Idle = 10 * time.Second
	}
	return cfg
}

func getDefaultedFileConfig(cfg *FileConfig) *FileConfig {
	if cfg == nil {
		cfg = &FileConfig{}
	}

	if cfg.Path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Consider how to handle this error more gracefully if panic is not desired
			panic(fmt.Errorf("error getting user home directory: %w", err))
		}
		cfg.Path = filepath.Join(homeDir, "Downloads")
	} else if strings.HasPrefix(cfg.Path, "~") { // Handle tilde expansion for Path
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Errorf("error getting user home directory for tilde expansion: %w", err))
		}
		cfg.Path = filepath.Join(homeDir, cfg.Path[2:])
	}

	if cfg.Name == "" {
		var s strings.Builder
		s.WriteString("ressource-")
		s.WriteString(time.Now().Format("2006-01-02-15-04-05"))
		s.WriteString(".txt")
		cfg.Name = s.String()
	}
	return cfg
}

func NewTrackClipboard(cfg *Config) *TrackClipboard {
	t := &TrackClipboard{
		Cfg: cfg,
	}

	t.Cfg.App = getDefaultedAPPConfig(t.Cfg.App)

	if t.Cfg.App.Channel == "local" {
		t.Cfg.File = getDefaultedFileConfig(t.Cfg.File)
		// Note: No specific error check needed here for FileConfig itself,
		// as getDefaultedFileConfig panics on critical errors like home dir resolution.
		// If it were to return an error, we'd handle it here.
	}

	// Look up the factory for the configured channel type
	factory, ok := channelFactories[t.Cfg.App.Channel]
	if !ok {
		// Or return an error: return nil, fmt.Errorf("unknown channel type: %s", t.Cfg.App.Channel)
		panic(fmt.Sprintf("unknown channel type: %s", t.Cfg.App.Channel))
	}

	// Create the channel using the factory
	channel, err := factory(t.Cfg)
	if err != nil {
		// Or return an error: return nil, fmt.Errorf("error creating channel %s: %w", t.Cfg.App.Channel, err)
		panic(fmt.Sprintf("error creating channel %s: %v", t.Cfg.App.Channel, err))
	}
	t.Channel = channel

	return t
}

func (t *TrackClipboard) handleClipboardData(ctx context.Context, data []byte, timer *time.Timer) {
	// Reset timer when clipboard content changes
	if !timer.Stop() {
		// If Stop() returns false, the timer has already fired or been stopped.
		// We might need to drain the channel if it fired, though in this select
		// structure, it's less likely to be a problem if we immediately Reset.
		select {
		case <-timer.C:
		default:
		}
	}

	if err := t.Channel.Send(ctx, string(data)); err != nil {
		// Error handling for t.Channel.Send
		// The current error message is reasonably informative.
		fmt.Printf("Error sending message via channel %s: %v\n", t.Cfg.App.Channel, err)
	}
	timer.Reset(t.Cfg.App.Idle)

	if t.Cfg.App.Debug {
		fmt.Printf("%s: Success sending clipboard content => %s\n", t.Cfg.App.Channel, string(data))
	}
}

func (t *TrackClipboard) Track() {
	err := clipboard.Init()
	if err != nil {
		// Panicking here is acceptable for a CLI tool if clipboard is essential.
		// For a library, returning an error might be more appropriate.
		panic(fmt.Errorf("failed to initialize clipboard: %w", err))
	}

	// Create a parent context for overall control.
	// Using context.Background() as the base is standard.
	parentCtx := context.Background()

	// Create a cancellable context for managing the lifecycle of clipboard watching.
	ctx, cancel := context.WithCancel(parentCtx)
	// Defer cancel() to ensure resources are released when Track() exits.
	// Also defer closing the channel.
	defer func() {
		cancel() // Signal all operations using this context to stop.
		if t.Channel != nil {
			if err := t.Channel.Close(); err != nil {
				fmt.Printf("Error closing channel %s: %v\n", t.Cfg.App.Channel, err)
			}
		}
	}()

	// Create timer for idle timeout. This timer determines how long to wait
	// for new clipboard activity before exiting.
	timer := time.NewTimer(t.Cfg.App.Idle)
	// Defer timer.Stop() to release timer resources.
	defer timer.Stop()

	// Watch clipboard changes. clipboard.Watch returns a channel that receives
	// clipboard content when it changes.
	content := clipboard.Watch(ctx, clipboard.FmtText)

	fmt.Printf("Starting clipboard tracking. Idle timeout: %s. Channel: %s\n", t.Cfg.App.Idle, t.Cfg.App.Channel)

	for {
		select {
		case data := <-content:
			t.handleClipboardData(ctx, data, timer)
		case <-ctx.Done():
			// This case handles explicit cancellation of the context,
			// for example, if the application is shutting down.
			fmt.Println("Context canceled, stopping clipboard tracking.")
			return
		case <-timer.C:
			// This case handles the idle timeout. If no clipboard activity
			// occurs for the duration of t.Cfg.App.Idle, the timer fires.
			fmt.Println("Idle timeout reached, stopping clipboard tracking.")
			return
		}
	}
}
