package runtest

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/vldcreation/helpme-package/pkg/consts"
)

func TestRunTest(t *testing.T) {
	testDataDir := "testdata"
	tests := []struct {
		name             string
		fpath            string
		funcName         string
		mustReturnOutput bool
		inputPath        string
		sampleOutputPath string
		expectedContains []string
		expectError      bool
	}{
		{
			name:             "Success test with helloworld",
			fpath:            filepath.Join(testDataDir, "helloworld.go"),
			funcName:         "HelloWorld",
			mustReturnOutput: true,
			inputPath:        filepath.Join(testDataDir, "helloworld.in"),
			sampleOutputPath: filepath.Join(testDataDir, "helloworld.out"),
			expectedContains: []string{
				consts.GREEN + "test passed" + consts.RESET,
			},
			expectError: false,
		},
		{
			name:             "Test with invalid input file",
			fpath:            filepath.Join(testDataDir, "helloworld.go"),
			funcName:         "HelloWorld",
			mustReturnOutput: false,
			inputPath:        "",
			sampleOutputPath: filepath.Join(testDataDir, "helloworld.out"),
			expectedContains: []string{
				consts.RED + "input file invalid" + consts.RESET,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, out, err := RunTest(tt.fpath, tt.funcName, tt.mustReturnOutput, tt.inputPath, tt.sampleOutputPath)

			_ = out
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			for _, expected := range tt.expectedContains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected result to contain %q, got %q", expected, result)
				}
			}
		})
	}
}
