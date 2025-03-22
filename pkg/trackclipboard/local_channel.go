package trackclipboard

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
	// Expand tilde to home directory if present
	path := l.Path
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Error getting home directory: %v\n", err)
			return
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Ensure directory exists
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Construct full file path and open file
	fullPath := filepath.Join(path, l.Name)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer f.Close()

	for {
		select {
		case msg := <-l.msgChan:
			if _, err := f.Write([]byte(msg + "\n")); err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
			}
		case <-l.doneChan:
			return
		}
	}
}

func (l *LocalChannel) Close() error {
	// Signal the processMessages goroutine to stop
	close(l.doneChan)

	// Close message channel after ensuring processMessages has stopped
	close(l.msgChan)

	return nil
}
