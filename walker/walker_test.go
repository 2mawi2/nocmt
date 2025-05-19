package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkerBasic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "walker-basic-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	testFiles := map[string]string{
		"file1.go":           "package main\nfunc main() {}\n",
		"file2.go":           "package util\nfunc Util() {}\n",
		"subfolder/file3.go": "package sub\nfunc Sub() {}\n",
	}

	for file, content := range testFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	processedFiles := make(map[string]bool)
	mockProcessor := func(path string) error {
		processedFiles[path] = true
		return nil
	}

	walker := &Walker{}
	err = walker.Walk(tempDir, mockProcessor)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	for file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if !processedFiles[fullPath] {
			t.Errorf("Expected file to be processed but wasn't: %s", fullPath)
		}
	}
}

func TestWalkerWithGitignore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "walker-gitignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	testFiles := map[string]string{
		"file.go":                 "// This is a comment\npackage main\n",
		"file.log":                "log content",
		"important.log":           "important log",
		"node_modules/module.js":  "// JS comment\nfunction test() {}\n",
		"subfolder/file.go":       "// Another comment\npackage sub\n",
		"subfolder/file.log":      "sub log content",
		"subfolder/.gitignore":    "!*.log\n*.txt\n",
		"subfolder/file.txt":      "text content",
		"subfolder/important.txt": "important text",
		".vscode/settings.json":   "{ \"setting\": true }",
		"file.tmp":                "temporary content",
	}

	for file, content := range testFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	processedFiles := make(map[string]bool)
	mockProcessor := func(path string) error {
		processedFiles[path] = true
		return nil
	}

	walker := &Walker{}
	err = walker.Walk(tempDir, mockProcessor)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	expectedProcessed := []string{
		filepath.Join(tempDir, "file.go"),
		filepath.Join(tempDir, "subfolder/file.go"),
		filepath.Join(tempDir, "subfolder/file.log"),
	}

	expectedNotProcessed := []string{
		filepath.Join(tempDir, "file.log"),
		filepath.Join(tempDir, "important.log"),
		filepath.Join(tempDir, "node_modules/module.js"),
		filepath.Join(tempDir, "subfolder/file.txt"),
		filepath.Join(tempDir, "subfolder/important.txt"),
		filepath.Join(tempDir, ".vscode/settings.json"),
		filepath.Join(tempDir, "file.tmp"),
	}

	for _, path := range expectedProcessed {
		if !processedFiles[path] {
			t.Errorf("Expected file to be processed but wasn't: %s", path)
		}
	}

	for _, path := range expectedNotProcessed {
		if processedFiles[path] {
			t.Errorf("Expected file NOT to be processed but was: %s", path)
		}
	}
}

func TestWalkerWithProcessorError(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "walker-error-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	testFiles := []string{
		"file1.go",
		"file2.go",
		"file3.go",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	errorFile := filepath.Join(tempDir, "file2.go")
	processedFiles := make(map[string]bool)
	errorProcessor := func(path string) error {
		if path == errorFile {
			return os.ErrNotExist
		}
		processedFiles[path] = true
		return nil
	}

	walker := &Walker{}
	err = walker.Walk(tempDir, errorProcessor)

	if err == nil {
		t.Errorf("Expected Walk to fail with an error, but it succeeded")
	}

	if processedFiles[errorFile] {
		t.Errorf("Error file was processed despite returning an error")
	}

	if !processedFiles[filepath.Join(tempDir, "file1.go")] {
		t.Errorf("File before error file wasn't processed")
	}
}
