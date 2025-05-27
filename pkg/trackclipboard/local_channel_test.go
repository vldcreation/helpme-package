package trackclipboard

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestEnsureLogFile(t *testing.T) {
	// Helper to create a temporary home directory for tilde expansion tests
	createTempHomeDir := func(t *testing.T) string {
		t.Helper()
		tempHome, err := os.MkdirTemp("", "testhome")
		if err != nil {
			t.Fatalf("Failed to create temp home dir: %v", err)
		}
		// Set HOME env var for os.UserHomeDir() to pick up, then restore
		origHome := os.Getenv("HOME")
		os.Setenv("HOME", tempHome)
		t.Cleanup(func() {
			os.Setenv("HOME", origHome)
			os.RemoveAll(tempHome)
		})
		return tempHome
	}

	// Test cases
	tests := []struct {
		name        string
		path        string
		filename    string
		setup       func(t *testing.T, testPath string) // Optional setup for specific conditions
		wantErr     bool
		wantErrMsg  string // Expected part of the error message
		checkFile   bool   // Whether to check for file existence and permissions
		checkParent bool   // Whether to check for parent directory existence
	}{
		{
			name:      "valid path and name",
			path:      filepath.Join(os.TempDir(), "trackclipboard_test_valid"),
			filename:  "test.log",
			wantErr:   false,
			checkFile: true,
			checkParent: true,
		},
		{
			name:     "tilde expansion",
			path:     "~/trackclipboard_test_tilde",
			filename: "tilde.log",
			setup: func(t *testing.T, testPath string) {
				// createTempHomeDir will be called by the test case itself
			},
			wantErr:   false,
			checkFile: true,
			checkParent: true, // The parent is the temp home dir
		},
		{
			name:        "unwritable path - directory creation fails",
			path:        filepath.Join("/root", "nonexistent_unwritable_dir_test"), // Assuming /root is not writable by test user
			filename:    "unwritable.log",
			wantErr:     true,
			wantErrMsg:  "error creating directory",
			checkFile:   false,
			checkParent: false,
		},
		{
			name:     "unwritable path - file creation fails",
			path:     os.TempDir(), // Writable directory
			filename: "readonly_dir/unwritable_file.log", // Non-existent sub-directory
			setup: func(t *testing.T, testPath string) {
				// Create a read-only directory for the file
				readOnlyDirPath := filepath.Join(testPath, "readonly_dir")
				if err := os.Mkdir(readOnlyDirPath, 0555); err != nil { // Read-only permissions
					t.Fatalf("Failed to create read-only dir: %v", err)
				}
				t.Cleanup(func() {
					// Best effort to remove, may need to change perms first if test failed mid-way
					os.Chmod(readOnlyDirPath, 0755)
					os.RemoveAll(readOnlyDirPath)
				})
			},
			wantErr:    true,
			wantErrMsg: "error opening file",
			checkFile:  false, // File should not be created
			checkParent: true, // Parent (os.TempDir()) exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testPath string
			if strings.HasPrefix(tt.path, "~") {
				tempHome := createTempHomeDir(t)
				// tt.path will be expanded by ensureLogFile using the tempHome
				// For cleanup, we need the actual path it would try to create inside tempHome
				testPath = filepath.Join(tempHome, tt.path[2:])
			} else {
				testPath = tt.path
			}
			fullFilePath := filepath.Join(testPath, tt.filename)

			// General cleanup for files/dirs created by ensureLogFile
			// This runs after specific test cleanups (like for read-only dir)
			t.Cleanup(func() {
				os.Remove(fullFilePath) // Remove the file
				// Attempt to remove the directory if it's not the system temp or a special test dir
				// This is a bit tricky because of tilde expansion and shared temp dirs.
				// We'll remove if it's a subdirectory of our specific test paths.
				if strings.Contains(testPath, "trackclipboard_test_") {
					os.RemoveAll(testPath)
				}
			})
			
			if tt.setup != nil {
				// For "unwritable path - file creation fails", testPath is os.TempDir()
				// For "tilde expansion", testPath is the expanded path inside tempHome
				setupPath := tt.path
				if strings.HasPrefix(tt.path,"~") {
					// The setup for tilde expansion expects the original tilde path
					// but createTempHomeDir handles the actual home dir creation.
					// For the specific "unwritable path - file creation fails", tt.path is os.TempDir()
					// and setup function will create "readonly_dir" inside it.
				} else if tt.name == "unwritable path - file creation fails" {
					setupPath = os.TempDir() // The setup creates a subdir in TempDir
				} else {
					setupPath = testPath
				}
                tt.setup(t, setupPath)
			}


			f, err := ensureLogFile(tt.path, tt.filename)

			if (err != nil) != tt.wantErr {
				t.Errorf("ensureLogFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.wantErrMsg != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErrMsg)) {
					t.Errorf("ensureLogFile() error = %v, wantErrMsg %q", err, tt.wantErrMsg)
				}
				return // Don't proceed with file checks if error was expected
			}

			// If no error was expected, file should be non-nil
			if f == nil {
				t.Fatalf("ensureLogFile() returned nil file, expected non-nil")
			}
			defer f.Close()

			if tt.checkFile {
				stat, err := os.Stat(fullFilePath)
				if err != nil {
					t.Fatalf("os.Stat(%q) failed: %v", fullFilePath, err)
				}
				if stat.IsDir() {
					t.Errorf("Expected file at %q, got directory", fullFilePath)
				}
				// Check permissions (more complex, focus on existence for now)
				// Example: if stat.Mode().Perm()&0644 != 0644 { t.Errorf(...) }
			}

			if tt.checkParent {
				parentDir := filepath.Dir(fullFilePath)
				stat, err := os.Stat(parentDir)
				if err != nil {
					t.Fatalf("os.Stat for parent dir %q failed: %v", parentDir, err)
				}
				if !stat.IsDir() {
					t.Errorf("Expected parent %q to be a directory", parentDir)
				}
			}
		})
	}
}

func TestLocalChannel_SendProcessClose(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "localchannel_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := &FileConfig{
		Path: tempDir,
		Name: "test_output.log",
	}
	fullLogPath := filepath.Join(cfg.Path, cfg.Name)

	lc := NewLocalChannel(cfg).(*LocalChannel) // Cast to access internal fields for test validation if needed

	// Send some messages
	messages := []string{"hello world", "test message 1", "another line"}
	for _, msg := range messages {
		if err := lc.Send(context.Background(), msg); err != nil {
			t.Errorf("Send() error = %v, want nil", err)
		}
	}

	// Give processMessages some time to write
	// A more robust way would be to check file content changes or use sync mechanisms
	// but for this test, a short sleep is often sufficient.
	// However, Close() should ensure messages are flushed (or written as it processes them).
	
	// Close the channel - this should signal processMessages to finish and close the file.
	if err := lc.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	// After Close(), Send should ideally fail or block indefinitely if context isn't cancellable
	// Test sending after close
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err = lc.Send(ctx, "message after close")
	if err == nil || !errors.Is(err, context.DeadlineExceeded) { // msgChan is closed, send will panic or block. Here it will block until ctx timeout
		// Actually, since msgChan is closed by Close(), sending on it will panic.
		// Let's verify this behavior.
		// We need to run this in a goroutine to catch the panic.
		var wg sync.WaitGroup
		wg.Add(1)
		var sendPanic interface{}
		go func() {
			defer func() {
				sendPanic = recover()
				wg.Done()
			}()
			_ = lc.Send(context.Background(), "message after close panic test") // Use fresh context
		}()
		wg.Wait()
		if sendPanic == nil {
			t.Errorf("Expected Send() to panic on closed channel, but it did not")
		}
	}


	// Verify file content
	file, err := os.Open(fullLogPath)
	if err != nil {
		t.Fatalf("Failed to open log file for verification: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("Error reading log file: %v", err)
	}

	if len(lines) != len(messages) {
		t.Errorf("Expected %d lines in log, got %d. Lines: %v", len(messages), len(lines), lines)
	}

	for i, expectedMsg := range messages {
		if i < len(lines) && lines[i] != expectedMsg {
			t.Errorf("Log line %d: expected %q, got %q", i, expectedMsg, lines[i])
		}
	}
}


func TestNewLocalChannel(t *testing.T) {
	cfg := &FileConfig{
		Path: "/tmp/testpath",
		Name: "testname.log",
	}
	lc := NewLocalChannel(cfg)

	if lc == nil {
		t.Fatal("NewLocalChannel returned nil")
	}

	localCh, ok := lc.(*LocalChannel)
	if !ok {
		t.Fatal("NewLocalChannel did not return a *LocalChannel")
	}

	if localCh.Path != cfg.Path {
		t.Errorf("Expected Path %q, got %q", cfg.Path, localCh.Path)
	}
	if localCh.Name != cfg.Name {
		t.Errorf("Expected Name %q, got %q", cfg.Name, localCh.Name)
	}
	if localCh.msgChan == nil {
		t.Error("Expected msgChan to be initialized, got nil")
	}
	if localCh.doneChan == nil {
		t.Error("Expected doneChan to be initialized, got nil")
	}

	// Important: Close the channel to stop the goroutine
	// This also requires a temporary directory for the log file to be created.
	tempDir, err := os.MkdirTemp("", "newlocalchannel_test_log")
	if err != nil {
		t.Fatalf("Failed to create temp dir for NewLocalChannel test: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	cfgActual := &FileConfig{Path: tempDir, Name: "actual.log"}
	lcToClose := NewLocalChannel(cfgActual)

	if err := lcToClose.Close(); err != nil {
		t.Errorf("Failed to close LocalChannel in test: %v", err)
	}
}

// Minimal TrackChannel interface definition for Config struct
type TrackChannel interface {
	Send(ctx context.Context, msg string) error
	Close() error
}

// Minimal Config struct definitions for FileConfig to be used in tests
type FileConfig struct {
	Path string
	Name string
}

// APPConfig and TelegramConfig are not needed for local_channel_test.go
// type APPConfig struct { ... }
// type TelegramConfig struct { ... }
// type Config struct { ... }
// Stubs for other types if they were in this file (they are in types.go or track.go)
// For example, if Config, APPConfig, TelegramConfig were defined here, they'd need stubs or actual definitions.
// Since they are likely in types.go or similar, this test file will compile if those types are accessible.
// The provided local_channel.go only uses FileConfig for NewLocalChannel.
// We need TrackChannel interface for NewLocalChannel's return type.
// We need FileConfig for NewLocalChannel's argument.
// The rest of the structs (APPConfig, TelegramConfig, Config) are not directly used by local_channel.go's functions.
// So, adding minimal stubs for FileConfig and TrackChannel interface here for self-containment of the test logic.
// Note: In a real Go module, these types would be imported from their respective packages/files.
// This test file assumes it's part of the 'trackclipboard' package and can access its types.
// If types.go exists and defines these, these stubs are not strictly necessary here but don't harm.
// The error message "TrackChannel not defined" implies that types.go is not being considered part of this test run's package scope
// or that it's missing. Let's assume types.go exists in the same package.
// To resolve potential "type not defined" issues when running tests in isolation,
// it's common to have a types.go or model.go file in the package.
// If this test file is run as `go test ./pkg/trackclipboard`, it should find other .go files in that package.

// For the sake of ensuring the test file is self-contained for analysis,
// if these types are not in the provided context, stubs are useful.
// However, the original file shows `package trackclipboard`, so it implies these types
// should be available from other files in the same package.
// The `LocalChannel` struct itself implements `TrackChannel`.
// Let's remove the stubs and assume the package structure is correct.
