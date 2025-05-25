package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllFlagFunctionality(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-all-flag-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	initGitRepo(t, tempDir)

	testFiles := map[string]string{
		"test.go": `package test

// This is a comment
func TestFunc() {
    // Another comment
    println("Hello")  // End of line comment
}`,
		"subdir/test.js": `// JavaScript comment
function test() {
    // Function comment
    console.log("test");
}`,
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", filepath.Dir(fullPath), err)
		}

		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	cmd := exec.Command(binaryPath, "-all")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt with -all flag: %v\nOutput: %s", err, output)
	}

	for path := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("Failed to read processed file %s: %v", path, err)
		}

		if strings.Contains(string(content), "//") {
			t.Errorf("Comments were not removed from file %s", path)
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Processing complete") {
		t.Errorf("Expected processing complete message, got: %s", outputStr)
	}
}

func TestDryRunFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-dry-run-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	initGitRepo(t, tempDir)

	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package test

// This is a comment
func TestFunc() {
    // Another comment
    println("Hello")  // End of line comment
}`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	originalContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read original file: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	cmd := exec.Command(binaryPath, testFile, "-dry-run")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt with -dry-run flag: %v\nOutput: %s", err, output)
	}

	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file after dry run: %v", err)
	}

	if string(originalContent) != string(afterContent) {
		t.Errorf("Dry run modified the file when it shouldn't have")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "package test") {
		t.Errorf("Dry run should output the processed content")
	}
}

func TestIgnorePatternFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-ignore-pattern-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	initGitRepo(t, tempDir)

	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package test

// TODO: This comment should be preserved
func TestFunc() {
    // Regular comment that should be removed
    // TODO: Another comment to preserve
    println("Hello")  // End of line comment
}`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	configFile := filepath.Join(tempDir, ".nocmt.json")
	configContent := `{
  "ignorePatterns": [
    "TODO"
  ]
}`

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	cmd := exec.Command(binaryPath, testFile)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Command output: %s", string(output))
		t.Fatalf("Failed to run nocmt: %v", err)
	}

	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read processed file: %v", err)
	}

	processedContent := string(afterContent)
	if !strings.Contains(processedContent, "TODO: This comment should be preserved") {
		t.Errorf("Ignore pattern didn't preserve the TODO comment")
	}
	if !strings.Contains(processedContent, "TODO: Another comment to preserve") {
		t.Errorf("Ignore pattern didn't preserve the second TODO comment")
	}
	if strings.Contains(processedContent, "Regular comment that should be removed") {
		t.Errorf("Regular comment wasn't removed")
	}
	if strings.Contains(processedContent, "End of line comment") {
		t.Errorf("End of line comment wasn't removed")
	}
}

func TestVerboseFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-verbose-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	initGitRepo(t, tempDir)

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	testFile := filepath.Join(subDir, "test.go")
	testContent := `package test

// This is a comment
func TestFunc() {
    println("Hello")
}`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	cmd := exec.Command(binaryPath, subDir, "-verbose")
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt with -verbose flag: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Processing") {
		t.Errorf("Verbose output should include processing details")
	}
}

func TestRemoveDirectivesFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-remove-directives-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	initGitRepo(t, tempDir)

	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package test

//go:generate mockgen -source=test.go
// Regular comment
func TestFunc() {
    //go:noinline
    println("Hello")
}`

	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	testFile2 := filepath.Join(tempDir, "test2.go")
	err = os.WriteFile(testFile2, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create second test file: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	cmd := exec.Command(binaryPath, testFile)
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt without -remove-directives: %v", err)
	}

	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read processed file: %v", err)
	}

	processedContent := string(afterContent)
	if !strings.Contains(processedContent, "//go:generate") {
		t.Errorf("go:generate directive should be preserved by default")
	}
	if !strings.Contains(processedContent, "//go:noinline") {
		t.Errorf("go:noinline directive should be preserved by default")
	}
	if strings.Contains(processedContent, "// Regular comment") {
		t.Errorf("Regular comment should be removed")
	}

	cmd = exec.Command(binaryPath, testFile2, "-remove-directives")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt with -remove-directives: %v", err)
	}

	afterContent, err = os.ReadFile(testFile2)
	if err != nil {
		t.Fatalf("Failed to read processed file after remove-directives: %v", err)
	}

	processedContent = string(afterContent)
	if strings.Contains(processedContent, "//go:generate") {
		t.Errorf("go:generate directive should be removed with -remove-directives")
	}
	if strings.Contains(processedContent, "//go:noinline") {
		t.Errorf("go:noinline directive should be removed with -remove-directives")
	}
}

func TestStagedFilesFlag(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nocmt-staged-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Printf("Failed to remove temp directory: %v\n", err)
		}
	}()

	initGitRepo(t, tempDir)

	stagedFile := filepath.Join(tempDir, "staged.go")
	stagedContent := `package test

// This staged comment should be removed
func StagedFunc() {
    // Another comment
    println("Hello")
}`

	unstagedFile := filepath.Join(tempDir, "unstaged.go")
	unstagedContent := `package test

// This unstaged comment should remain
func UnstagedFunc() {
    // Another comment
    println("Hello")
}`

	if err := os.WriteFile(stagedFile, []byte(stagedContent), 0644); err != nil {
		t.Fatalf("Failed to create staged file: %v", err)
	}

	if err := os.WriteFile(unstagedFile, []byte(unstagedContent), 0644); err != nil {
		t.Fatalf("Failed to create unstaged file: %v", err)
	}

	stageCmd := exec.Command("git", "add", "staged.go")
	stageCmd.Dir = tempDir
	if err := stageCmd.Run(); err != nil {
		t.Fatalf("Failed to stage file: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "nocmt-test")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/nocmt")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build nocmt binary: %v", err)
	}

	runCmd := exec.Command(binaryPath, "-staged", "-verbose")
	runCmd.Dir = tempDir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run nocmt with -staged flag: %v\nOutput: %s", err, output)
	}

	stagedAfter, err := os.ReadFile(stagedFile)
	if err != nil {
		t.Fatalf("Failed to read staged file after processing: %v", err)
	}

	unstagedAfter, err := os.ReadFile(unstagedFile)
	if err != nil {
		t.Fatalf("Failed to read unstaged file after processing: %v", err)
	}

	if strings.Contains(string(stagedAfter), "This staged comment should be removed") {
		t.Errorf("Comments were not removed from staged file")
	}

	if !strings.Contains(string(unstagedAfter), "This unstaged comment should remain") {
		t.Errorf("Unstaged file was incorrectly modified")
	}
}

func initGitRepo(t *testing.T, dir string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	configCmd := exec.Command("git", "config", "user.email", "test@example.com")
	configCmd.Dir = dir
	if err := configCmd.Run(); err != nil {
		t.Fatalf("Failed to configure git email: %v", err)
	}

	configNameCmd := exec.Command("git", "config", "user.name", "Test User")
	configNameCmd.Dir = dir
	if err := configNameCmd.Run(); err != nil {
		t.Fatalf("Failed to configure git name: %v", err)
	}
}
