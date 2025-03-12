package runtest

import (
	"testing"
)

func TestGetPackageName(t *testing.T) {
	tests := []struct {
		name     string
		fpath    string
		expected string
	}{
		{
			name:     "simple package path",
			fpath:    "/path/to/mypackage/file.go",
			expected: "mypackage",
		},
		{
			name:     "nested package path",
			fpath:    "/path/to/parent/child/file.go",
			expected: "child",
		},
		{
			name:     "relative path",
			fpath:    "./mypackage/file.go",
			expected: "mypackage",
		},
		{
			name:     "file in root",
			fpath:    "/rootfile.go",
			expected: "",
		},
		{
			name:     "empty path",
			fpath:    "",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getPackageName(tt.fpath)
			if got != tt.expected {
				t.Errorf("getPackageName() = %s, want %s", got, tt.expected)
			}
		})
	}
}
