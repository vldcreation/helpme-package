package runtest

import (
	"os"
	"path/filepath"
	"strings"
)

// getPackagePath returns the full import path for a Go file
func getPackagePath(fpath string) string {
	// Convert to absolute path if it's not already
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return ""
	}

	// Find the "go.mod" file by walking up the directory tree
	dir := filepath.Dir(absPath)
	for dir != "/" {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Found go.mod, now construct the package path
			relPath, _ := filepath.Rel(dir, filepath.Dir(absPath))
			return filepath.Join(getModuleName(dir), relPath)
		}
		dir = filepath.Dir(dir)
	}

	return ""
}

// getPackageName extracts the package name from the file path
func getPackageName(fpath string) string {
	return filepath.Base(filepath.Dir(fpath))
}

// getModuleName reads the module name from go.mod file
func getModuleName(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}

	// Find the module line
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}

	return ""
}
