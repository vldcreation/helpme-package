package runtest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetPackageName(t *testing.T) {
	testDataDir := "testdata"
	tests := []struct {
		name             string
		fpath            string
		expectedContains string
	}{
		{
			name:             "simple package",
			fpath:            "helloworld.go",
			expectedContains: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(testDataDir, tt.fpath))
			if err != nil {
				t.Errorf("error reading file: %v", err)
			}
			got, err := getPackageName(data, "package main")
			if err != nil {
				t.Errorf("error getting package name: %v", err)
			}

			if !strings.Contains(got, tt.expectedContains) {
				t.Errorf("expected %s to contain %s", got, tt.expectedContains)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	tests := []struct {
		name     string
		fpath    string
		expected string
	}{
		{
			name:     "simple file",
			fpath:    "mydir/mysubdir/mysubsubdir/helloworld.go",
			expected: "helloworld.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFilename(tt.fpath)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestGetPackageName2(t *testing.T) {
	tests := []struct {
		name     string
		fpath    string
		expected string
	}{
		{
			name:     "simple file",
			fpath:    "mydir/mysubdir/mysubsubdir/helloworld.go",
			expected: "mysubsubdir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPackageName2(tt.fpath)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
