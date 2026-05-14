// Package tests provides end-to-end integration tests for the name-counter CLI.
// These tests compile the actual binary and execute it to verify the entire pipeline.
package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binPath string

// TestMain compiles the application binary before running the tests.
func TestMain(m *testing.M) {
	absPath, err := filepath.Abs("name-counter-e2e")
	if err != nil {
		os.Stderr.WriteString("Failed to resolve absolute path for test binary\n")
		os.Exit(1)
	}
	binPath = absPath

	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/name-counter")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Stderr.WriteString("Failed to build CLI binary for testing\n")
		os.Exit(1)
	}

	exitCode := m.Run()

	_ = os.Remove(binPath)
	os.Exit(exitCode)
}

// cliTestCase defines the structure for a single CLI test scenario.
type cliTestCase struct {
	name           string
	args           []string
	expectedStdout string
	expectedStderr string
	expectError    bool
}

func TestCLI_Integration(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "names.txt")
	testData := "Alice\nBob\nAlice\nCharlie\nBob\nBob\n"

	if err := os.WriteFile(inputFile, []byte(testData), 0o600); err != nil {
		t.Fatalf("Failed to write temporary test file: %v", err)
	}

	tests := []cliTestCase{
		{
			name:           "Default execution mode (descending)",
			args:           []string{inputFile},
			expectedStdout: "Bob:3\nAlice:2\nCharlie:1\n",
			expectError:    false,
		},
		{
			name:           "Ascending execution mode",
			args:           []string{inputFile, "--mode", "asc"},
			expectedStdout: "Charlie:1\nAlice:2\nBob:3\n",
			expectError:    false,
		},
		{
			name:        "None execution mode",
			args:        []string{inputFile, "--mode", "none"},
			expectError: false,
		},
		{
			name:           "Invalid mode flag",
			args:           []string{inputFile, "--mode", "unknown"},
			expectedStderr: "invalid mode: unknown",
			expectError:    true,
		},
		{
			name:           "Invalid batch size flag",
			args:           []string{inputFile, "--batch-size", "0"},
			expectedStderr: "invalid batch-size",
			expectError:    true,
		},
		{
			name:           "Missing file argument",
			args:           []string{},
			expectedStderr: "accepts 1 arg(s), received 0",
			expectError:    true,
		},
		{
			name:           "Non-existent file execution",
			args:           []string{"non_existent_file.txt"},
			expectedStderr: "failed to open file",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runAndAssertCLICommand(t, &tt)
		})
	}
}

// runAndAssertCLICommand executes the CLI binary and delegates output validation.
func runAndAssertCLICommand(t *testing.T, tc *cliTestCase) {
	t.Helper()

	cmd := exec.Command(binPath, tc.args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	assertExecutionStatus(t, tc, err, stderr.String())
	assertStdout(t, tc, stdout.String())
	assertStderr(t, tc, stderr.String())
}

// assertExecutionStatus verifies the exit code matches expectations.
func assertExecutionStatus(t *testing.T, tc *cliTestCase, err error, gotStderr string) {
	t.Helper()
	if tc.expectError && err == nil {
		t.Fatalf("Expected execution to fail, but it succeeded")
	}
	if !tc.expectError && err != nil {
		t.Fatalf("Expected successful execution, got error: %v\nStderr: %s", err, gotStderr)
	}
}

// assertStdout verifies the standard output matches expectations.
func assertStdout(t *testing.T, tc *cliTestCase, gotStdout string) {
	t.Helper()
	if tc.expectedStdout != "" && gotStdout != tc.expectedStdout {
		t.Errorf("Unexpected stdout.\nGot:\n%s\nExpected:\n%s", gotStdout, tc.expectedStdout)
	}
}

// assertStderr verifies the standard error contains the expected substring.
func assertStderr(t *testing.T, tc *cliTestCase, gotStderr string) {
	t.Helper()
	if tc.expectedStderr != "" && !strings.Contains(gotStderr, tc.expectedStderr) {
		t.Errorf("Unexpected stderr.\nGot:\n%s\nExpected to contain:\n%s", gotStderr, tc.expectedStderr)
	}
}
