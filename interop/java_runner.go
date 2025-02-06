package interop

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type JavaRunner struct{ I Interop }

func (runner *JavaRunner) Run() (string, error) {
	// Get absolute path
	absPath, err := filepath.Abs(runner.I.FilePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	javaProgram, ok := strings.CutSuffix(absPath, ".java")

	if !ok {
		return "", fmt.Errorf("file is not a java program")
	}

	cmd := exec.Command("javac", absPath, "&&", "java", javaProgram)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run script: %v\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
