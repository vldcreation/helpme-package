package trackclipboard

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ensureLogFile handles the creation and opening of the log file.
// It expands tilde in the path, creates the directory if necessary,
// and opens the file in append mode.
func ensureLogFile(path, name string) (*os.File, error) {
	// Expand tilde to home directory if present
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Ensure directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("error creating directory: %w", err)
	}

	// Construct full file path and open file
	fullPath := filepath.Join(path, name)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	return f, nil
}

type LocalChannel struct {
	Path     string
	Name     string
	file     *os.File
	msgChan  chan string
	doneChan chan struct{}
}

func NewLocalChannel(cfg *FileConfig) TrackChannel {
	channel := &LocalChannel{
		Path:     cfg.Path,
		Name:     cfg.Name,
		msgChan:  make(chan string, 100),
		doneChan: make(chan struct{}),
	}

	// Start the message processing goroutine
	go channel.processMessages()

	return channel
}

func (l *LocalChannel) Send(ctx context.Context, msg string) error {
	select {
	case l.msgChan <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *LocalChannel) processMessages() {
	f, err := ensureLogFile(l.Path, l.Name)
	if err != nil {
		fmt.Printf("Error ensuring log file: %v\n", err)
		return
	}
	l.file = f // Store the file handle in the struct
	defer l.file.Close()

	for {
		select {
		case msg := <-l.msgChan:
			if _, err := l.file.Write([]byte(msg + "\n")); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
		case <-l.doneChan:
			return
		}
	}
}

func (l *LocalChannel) Close() error {
	// Signal the processMessages goroutine to stop
	if l.doneChan != nil {
		close(l.doneChan)
	}

	// Close message channel after ensuring processMessages has stopped
	if l.msgChan != nil {
		close(l.msgChan)
	}
	
	// The file will be closed by the defer statement in processMessages.
	// If processMessages didn't run (e.g. due to error in ensureLogFile),
	// l.file might be nil.
	if l.file != nil {
		// We might want to ensure a final flush or any other cleanup specific to the file here,
		// but the primary closing is handled by defer in processMessages.
		// For now, we'll rely on the defer in processMessages.
		// If ensureLogFile fails, l.file is nil and processMessages returns early.
		// If ensureLogFile succeeds, defer will close it when processMessages returns.
	}

	return nil
}
