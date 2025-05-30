package interop

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

type JavascriptRunner struct{ I Interop }

func (runner *JavascriptRunner) Run() (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(runner.I.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	cmd := exec.Command("node", append([]string{absPath}, runner.I.Args...)...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run script: %v\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
