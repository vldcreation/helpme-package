package pull

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Pull(host string, username string, repo string, branch string) error {
	lookUp(&host, &branch)

	// Construct the GitHub repository URL
	repoURL := fmt.Sprintf("%s/%s/%s", host, username, repo)

	// Set GOPRIVATE to allow private repository access
	if err := os.Setenv("GOPRIVATE", host); err != nil {
		return fmt.Errorf("failed to set GOPRIVATE: %w", err)
	}

	getCmd := exec.Command("go", "get", fmt.Sprintf("%s@%s", repoURL, branch))

	// Set environment variables for go get
	getCmd.Env = append(os.Environ(), "GO111MODULE=on")

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	netrcPath := filepath.Join(home, ".netrc")
	if _, err := os.Stat(netrcPath); err != nil {
		return fmt.Errorf("ensure .netrc exists at %s for private repo access: %w", netrcPath, err)
	}

	output, err := getCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull repository: %s: %w", string(output), err)
	}

	return nil
}

func lookUp(host *string, branch *string) {
	if *host == "" {
		*host = GITHUB_URL
	}

	if *branch == "" {
		*branch = "main"
	}
}
