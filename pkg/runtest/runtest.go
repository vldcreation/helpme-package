package runtest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/vldcreation/helpme-package/pkg/consts"
)

func RunTest(fpath string, funcName string, mustReturnOutput bool, inputPath string, sampleOutputPath string) (string, string, error) {
	var output strings.Builder
	if inputPath == "" || !strings.HasSuffix(inputPath, ".in") {
		buildFailureMessage(&output, "input file invalid")
		return output.String(), "", nil
	}

	if sampleOutputPath == "" || !strings.HasSuffix(sampleOutputPath, ".out") {
		buildFailureMessage(&output, "sample output file invalid")
		return output.String(), "", nil
	}

	if funcName == "" {
		funcName = getFilenameWithoutExtension(fpath)
	}

	input, err := os.ReadFile(inputPath)
	if err != nil {
		return "", "", err
	}

	sampleOutput, err := os.ReadFile(sampleOutputPath)
	if err != nil {
		return "", "", err
	}

	res, err := runTestWithOutput(fpath, funcName, input)
	if err != nil {
		return "", "", err
	}
	if mustReturnOutput {
		output.WriteString("Output:")
		output.WriteString("\n")
		output.WriteString(consts.BLUE)
		output.WriteString(res)
		output.WriteString(consts.RESET)
	}

	if res != string(sampleOutput) {
		buildDiffMessage(&output, string(sampleOutput), res)
		buildFailureMessage(&output, "test failed")
		return output.String(), res, nil
	}

	buildSuccessMessage(&output, "test passed")
	return output.String(), res, nil
}

func runTestWithOutput(fpath string, funcName string, input []byte) (string, error) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "runtest")
	if err != nil {
		return "", fmt.Errorf("error creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary main.go file
	fl, err := os.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	flString, err := getPackageName(fl, "main")
	if err != nil {
		return "", err
	}

	tmpMainPath := filepath.Join(tmpDir, "main.go")
	mainContent := fmt.Sprintf(`
%s
func main() {
	%s()
}
`, flString, funcName)

	if err := os.WriteFile(tmpMainPath, []byte(mainContent), 0644); err != nil {
		return "", fmt.Errorf("error writing temp main file: %v", err)
	}

	// Run the temporary main package
	cmd := exec.Command("go", "run", tmpMainPath)

	// Set up input pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("error creating stdin pipe: %v", err)
	}

	// Set up output capture
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting command: %v", err)
	}

	// Write input to stdin
	if _, err := stdin.Write(input); err != nil {
		return "", fmt.Errorf("error writing to stdin: %v", err)
	}
	stdin.Close()

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("error running command: %v\nstderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

func buildDiffMessage(s *strings.Builder, expected string, got string) {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(got, "\n")

	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	for i := 0; i < maxLines; i++ {
		expectedLine := ""
		actualLine := ""

		if i < len(expectedLines) {
			expectedLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actualLine = actualLines[i]
		}

		if expectedLine == actualLine {
			s.WriteString(expectedLine + "\n")
		} else {
			s.WriteString(consts.RED + actualLine + consts.RESET + "\n")
		}
	}
}

func buildSuccessMessage(s *strings.Builder, msg string) {
	s.WriteString(consts.GREEN)
	s.WriteString(msg)
	s.WriteString(consts.RESET)
}

func buildFailureMessage(s *strings.Builder, msg string) {
	s.WriteString(consts.RED)
	s.WriteString(msg)
	s.WriteString(consts.RESET)
}
