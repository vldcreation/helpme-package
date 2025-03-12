package runtest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func getFilename(fpath string) string {
	return filepath.Base(fpath)
}

func getFilenameWithoutExtension(fpath string) string {
	return strings.TrimSuffix(filepath.Base(fpath), filepath.Ext(fpath))
}

func getPackageName2(fpath string) string {
	dir := filepath.Dir(fpath)
	base := filepath.Base(dir)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// getPackageName extracts the package name from the file path
func getPackageName(data []byte, newPackageName string) (modified string, err error) { // Define a regex to match the package declaration
	re := regexp.MustCompile(`(?m)^package\s+\w+`)
	newPackageDecl := fmt.Sprintf("package %s", newPackageName)

	// Replace the package declaration with the new one
	modifiedData := re.ReplaceAll(data, []byte(newPackageDecl))
	return string(modifiedData), nil
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
