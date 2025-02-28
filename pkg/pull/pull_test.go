package pull

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPull(t *testing.T) {
	// Create temporary .netrc file for testing
	tmpHome, err := os.MkdirTemp("", "test-home")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpHome)

	// Create .netrc file
	netrcPath := filepath.Join(tmpHome, ".netrc")
	if err := os.WriteFile(netrcPath, []byte("machine github.com\nlogin test\npassword test"), 0600); err != nil {
		t.Fatal(err)
	}

	// Override home directory for testing
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpHome)
	defer os.Setenv("HOME", originalHome)

	tests := []struct {
		name     string
		host     string
		username string
		repo     string
		branch   string
		wantErr  bool
	}{
		{
			name:     "Valid public repository with default host and branch",
			host:     "",
			username: "golang",
			repo:     "go",
			branch:   "",
			wantErr:  false,
		},
		{
			name:     "Valid repository with custom host and branch",
			host:     "github.com",
			username: "golang",
			repo:     "go",
			branch:   "master",
			wantErr:  false,
		},
		{
			name:     "Invalid repository",
			host:     "github.com",
			username: "invalid-user-12345",
			repo:     "invalid-repo-12345",
			branch:   "main",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Pull(tt.host, tt.username, tt.repo, tt.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pull() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLookUp(t *testing.T) {
	tests := []struct {
		name        string
		inputHost   string
		inputBranch string
		wantHost    string
		wantBranch  string
	}{
		{
			name:        "Empty host and branch",
			inputHost:   "",
			inputBranch: "",
			wantHost:    GITHUB_URL,
			wantBranch:  "main",
		},
		{
			name:        "Custom host and empty branch",
			inputHost:   "gitlab.com",
			inputBranch: "",
			wantHost:    "gitlab.com",
			wantBranch:  "main",
		},
		{
			name:        "Empty host and custom branch",
			inputHost:   "",
			inputBranch: "develop",
			wantHost:    GITHUB_URL,
			wantBranch:  "develop",
		},
		{
			name:        "Custom host and branch",
			inputHost:   "gitlab.com",
			inputBranch: "develop",
			wantHost:    "gitlab.com",
			wantBranch:  "develop",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host := tt.inputHost
			branch := tt.inputBranch
			lookUp(&host, &branch)

			if host != tt.wantHost {
				t.Errorf("lookUp() host = %v, want %v", host, tt.wantHost)
			}
			if branch != tt.wantBranch {
				t.Errorf("lookUp() branch = %v, want %v", branch, tt.wantBranch)
			}
		})
	}
}