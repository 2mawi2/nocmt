package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestConfigPatternMatching(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		patterns []string
		want     bool
	}{
		{
			name:     "match TODO pattern",
			comment:  "// TODO: implement this",
			patterns: []string{"TODO"},
			want:     true,
		},
		{
			name:     "no match without patterns",
			comment:  "// Regular comment",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "match with regex",
			comment:  "// TESTPROJECT-12345: Fix issue",
			patterns: []string{"TESTPROJECT-\\d+"},
			want:     true,
		},
		{
			name:     "no match with non-matching regex",
			comment:  "// KBA-12345: Fix issue",
			patterns: []string{"TESTPROJECT-\\d+"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			err := cfg.SetCLIPatterns(tt.patterns)
			if err != nil {
				t.Fatalf("SetCLIPatterns() error = %v", err)
			}

			if got := cfg.ShouldIgnoreComment(tt.comment); got != tt.want {
				t.Errorf("ShouldIgnoreComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigLoadSave(t *testing.T) {
	tempDir := t.TempDir()
	homeDir := filepath.Join(tempDir, "home")
	projectDir := filepath.Join(tempDir, "project")

	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create home directory: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	globalConfig := filepath.Join(homeDir, ".nocmt")
	if err := os.MkdirAll(globalConfig, 0755); err != nil {
		t.Fatalf("Failed to create global config directory: %v", err)
	}
	globalConfigFile := filepath.Join(globalConfig, "config.json")
	globalConfigContent := `{
		"ignorePatterns": ["TODO", "FIXME"]
	}`
	if err := os.WriteFile(globalConfigFile, []byte(globalConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write global config file: %v", err)
	}

	localConfigFile := filepath.Join(projectDir, ".nocmt.json")
	localConfigContent := `{
		"ignorePatterns": ["TICKET", "JIRA"]
	}`
	if err := os.WriteFile(localConfigFile, []byte(localConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write local config file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	if err := os.Chdir(projectDir); err != nil {
		t.Fatalf("Failed to change to project directory: %v", err)
	}

	oldHome := os.Getenv("HOME")
	defer func() {
		if err := os.Setenv("HOME", oldHome); err != nil {
			t.Logf("Failed to restore HOME environment variable: %v", err)
		}
	}()

	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("Failed to set HOME environment variable: %v", err)
	}

	cfg := New()
	if err := cfg.LoadConfigurations(); err != nil {
		t.Fatalf("LoadConfigurations() error = %v", err)
	}

	testCases := []struct {
		comment string
		want    bool
	}{
		{"// TODO: something", true},
		{"// FIXME: something", true},
		{"// TICKET: 1234", true},
		{"// JIRA: ABC-123", true},
		{"// Regular comment", false},
	}

	for _, tc := range testCases {
		if got := cfg.ShouldIgnoreComment(tc.comment); got != tc.want {
			t.Errorf("ShouldIgnoreComment(%q) = %v, want %v", tc.comment, got, tc.want)
		}
	}

	err = cfg.SetCLIPatterns([]string{"WHY"})
	if err != nil {
		t.Fatalf("SetCLIPatterns() error = %v", err)
	}

	if !cfg.ShouldIgnoreComment("// WHY: Explanation") {
		t.Errorf("SetCLIPatterns() failed to add pattern")
	}
}

func TestShouldIgnoreFile(t *testing.T) {
	tests := []struct {
		name         string
		filename     string
		patterns     []string
		shouldIgnore bool
	}{
		{
			name:         "no patterns",
			filename:     "example.go",
			patterns:     []string{},
			shouldIgnore: false,
		},
		{
			name:         "exact match",
			filename:     "example.go",
			patterns:     []string{"example.go"},
			shouldIgnore: true,
		},
		{
			name:         "extension match",
			filename:     "example.go",
			patterns:     []string{"\\.go$"},
			shouldIgnore: true,
		},
		{
			name:         "directory match",
			filename:     "src/models/example.go",
			patterns:     []string{"models/"},
			shouldIgnore: true,
		},
		{
			name:         "pattern doesn't match",
			filename:     "example.go",
			patterns:     []string{"\\.js$", "\\.py$"},
			shouldIgnore: false,
		},
		{
			name:         "complex pattern match",
			filename:     "test/fixtures/example.test.go",
			patterns:     []string{"test/fixtures/.*\\.test\\."},
			shouldIgnore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := New()
			err := cfg.SetCLIFilePatterns(tt.patterns)
			if err != nil {
				t.Fatalf("Failed to set patterns: %v", err)
			}

			if got := cfg.ShouldIgnoreFile(tt.filename); got != tt.shouldIgnore {
				t.Errorf("ShouldIgnoreFile() = %v, want %v", got, tt.shouldIgnore)
			}
		})
	}
}

func TestFileIgnorePatternConfig(t *testing.T) {
	cfg := New()

	pattern1 := "\\.log$"
	pattern2 := "node_modules/"

	cfg.Local.FileIgnorePatterns = []string{pattern1}
	cfg.Local.FileIgnorePatterns = append(cfg.Local.FileIgnorePatterns, pattern2)

	err := cfg.compilePatterns()
	if err != nil {
		t.Fatalf("Failed to compile patterns: %v", err)
	}

	if !cfg.ShouldIgnoreFile("test.log") {
		t.Errorf("Should ignore .log files")
	}

	if !cfg.ShouldIgnoreFile("src/node_modules/package.json") {
		t.Errorf("Should ignore files in node_modules directory")
	}

	if cfg.ShouldIgnoreFile("test.go") {
		t.Errorf("Should not ignore .go files")
	}

	err = cfg.SetCLIFilePatterns([]string{"temp/"})
	if err != nil {
		t.Fatalf("Failed to set CLI patterns: %v", err)
	}

	if !cfg.ShouldIgnoreFile("temp/cache.dat") {
		t.Errorf("Should ignore files in temp directory from CLI patterns")
	}
}

func TestDirectoryIgnorePatterns(t *testing.T) {
	tempDir := t.TempDir()

	homeDir := filepath.Join(tempDir, "home")
	projectDir := filepath.Join(tempDir, "project")

	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("Failed to create home directory: %v", err)
	}
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	testCases := []struct {
		name          string
		globalPattern string
		localPattern  string
		cliPattern    string
		filePath      string
		shouldIgnore  bool
	}{
		{
			name:         "local config directory pattern",
			localPattern: "^processor/testdata/",
			filePath:     "processor/testdata/bash/original.sh",
			shouldIgnore: true,
		},
		{
			name:         "local config nested directory pattern",
			localPattern: "^processor/testdata/",
			filePath:     "processor/testdata/nested/dir/file.go",
			shouldIgnore: true,
		},
		{
			name:         "local config directory not matching",
			localPattern: "^processor/testdata/",
			filePath:     "processor/src/file.go",
			shouldIgnore: false,
		},
		{
			name:          "global config directory pattern",
			globalPattern: "^vendor/",
			filePath:      "vendor/lib/file.js",
			shouldIgnore:  true,
		},
		{
			name:          "cli pattern overrides config",
			globalPattern: "^processor/",
			localPattern:  "",
			cliPattern:    "^processor/scripts/.*\\.sh$",
			filePath:      "processor/scripts/test.sh",
			shouldIgnore:  true,
		},
		{
			name:          "cli pattern shouldn't ignore exempted file",
			globalPattern: "",
			localPattern:  "",
			cliPattern:    "^processor/scripts/",
			filePath:      "processor/testdata/bash/original.sh",
			shouldIgnore:  false,
		},
		{
			name:         "multiple local patterns",
			localPattern: "^processor/testdata/,^vendor/",
			filePath:     "vendor/lib/file.js",
			shouldIgnore: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			homeDir := filepath.Join(tempDir, "home")
			projectDir := filepath.Join(tempDir, "project")

			if err := os.MkdirAll(homeDir, 0755); err != nil {
				t.Fatalf("Failed to create home directory: %v", err)
			}
			if err := os.MkdirAll(projectDir, 0755); err != nil {
				t.Fatalf("Failed to create project directory: %v", err)
			}

			if tc.globalPattern != "" {
				globalConfigDir := filepath.Join(homeDir, ".nocmt")
				if err := os.MkdirAll(globalConfigDir, 0755); err != nil {
					t.Fatalf("Failed to create global config directory: %v", err)
				}

				globalConfigFile := filepath.Join(globalConfigDir, "config.json")
				var patterns []string
				if tc.globalPattern != "" {
					patterns = splitPatterns(tc.globalPattern)
				}

				globalConfig := CommentConfig{
					FileIgnorePatterns: patterns,
				}

				data, err := json.MarshalIndent(globalConfig, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal global config: %v", err)
				}

				if err := os.WriteFile(globalConfigFile, data, 0644); err != nil {
					t.Fatalf("Failed to write global config file: %v", err)
				}
			}

			if tc.localPattern != "" {
				localConfigFile := filepath.Join(projectDir, ".nocmt.json")
				var patterns []string
				if tc.localPattern != "" {
					patterns = splitPatterns(tc.localPattern)
				}

				localConfig := CommentConfig{
					FileIgnorePatterns: patterns,
				}

				data, err := json.MarshalIndent(localConfig, "", "  ")
				if err != nil {
					t.Fatalf("Failed to marshal local config: %v", err)
				}

				if err := os.WriteFile(localConfigFile, data, 0644); err != nil {
					t.Fatalf("Failed to write local config file: %v", err)
				}
			}

			oldHome := os.Getenv("HOME")
			oldWd, err := os.Getwd()

			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}

			defer func() {
				if err := os.Chdir(oldWd); err != nil {
					t.Logf("Failed to restore working directory: %v", err)
				}
				if err := os.Setenv("HOME", oldHome); err != nil {
					t.Logf("Failed to restore HOME environment: %v", err)
				}
			}()

			if err := os.Setenv("HOME", homeDir); err != nil {
				t.Fatalf("Failed to set HOME environment: %v", err)
			}
			if err := os.Chdir(projectDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			cfg := New()
			if err := cfg.LoadConfigurations(); err != nil {
				t.Fatalf("Failed to load configurations: %v", err)
			}

			if tc.cliPattern != "" {
				patterns := splitPatterns(tc.cliPattern)
				if err := cfg.SetCLIFilePatterns(patterns); err != nil {
					t.Fatalf("Failed to set CLI patterns: %v", err)
				}
			}

			result := cfg.ShouldIgnoreFile(tc.filePath)

			if result != tc.shouldIgnore {
				t.Errorf("ShouldIgnoreFile(%q) = %v, want %v", tc.filePath, result, tc.shouldIgnore)
			}
		})
	}
}

func splitPatterns(patternsStr string) []string {
	if patternsStr == "" {
		return nil
	}
	var patterns []string
	for _, p := range strings.Split(patternsStr, ",") {
		patterns = append(patterns, strings.TrimSpace(p))
	}
	return patterns
}

func TestShouldIgnoreFileForPreCommit(t *testing.T) {
	tmpDir := t.TempDir()

	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	localConfig := CommentConfig{
		FileIgnorePatterns: []string{"^processor/testdata/"},
	}

	configFile := filepath.Join(repoDir, ".nocmt.json")
	data, err := json.MarshalIndent(localConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	testdataDir := filepath.Join(repoDir, "processor", "testdata", "bash")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatalf("Failed to create test directories: %v", err)
	}

	testFile := filepath.Join(testdataDir, "test.sh")
	if err := os.WriteFile(testFile, []byte("# Test file"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	srcDir := filepath.Join(repoDir, "processor", "src")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create src directory: %v", err)
	}

	srcFile := filepath.Join(srcDir, "main.go")
	if err := os.WriteFile(srcFile, []byte("// Main file"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo dir: %v", err)
	}

	cfg := New()
	if err := cfg.LoadConfigurations(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	tests := []struct {
		name         string
		file         string
		shouldIgnore bool
	}{
		{
			name:         "test file in ignored directory",
			file:         "processor/testdata/bash/test.sh",
			shouldIgnore: true,
		},
		{
			name:         "regular source file",
			file:         "processor/src/main.go",
			shouldIgnore: false,
		},
		{
			name:         "absolute path to ignored file",
			file:         filepath.Join(repoDir, "processor/testdata/bash/test.sh"),
			shouldIgnore: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.ShouldIgnoreFile(tt.file); got != tt.shouldIgnore {
				t.Errorf("ShouldIgnoreFile(%q) = %v, want %v", tt.file, got, tt.shouldIgnore)
			}
		})
	}
}

func TestDirectoryPatternIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	testDirs := []string{
		filepath.Join(repoDir, "processor", "testdata", "bash"),
		filepath.Join(repoDir, "processor", "testdata", "go"),
		filepath.Join(repoDir, "processor", "src"),
		filepath.Join(repoDir, "config", "testdata"),
		filepath.Join(repoDir, "walker", "testdata"),
	}

	for _, dir := range testDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	testFiles := map[string]string{
		"processor/testdata/bash/test.sh": "# Bash test file",
		"processor/testdata/go/test.go":   "// Go test file",
		"processor/src/main.go":           "// Main source file",
		"config/testdata/config.json":     "// Config test file",
		"walker/testdata/mock.go":         "// Walker test file",
	}

	for relPath, content := range testFiles {
		fullPath := filepath.Join(repoDir, relPath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	testCases := []struct {
		name         string
		config       CommentConfig
		relativePath string
		tests        []struct {
			file         string
			shouldIgnore bool
		}
	}{
		{
			name: "exact pattern match",
			config: CommentConfig{
				FileIgnorePatterns: []string{"^processor/testdata/bash/"},
			},
			relativePath: "",
			tests: []struct {
				file         string
				shouldIgnore bool
			}{
				{"processor/testdata/bash/test.sh", true},
				{"processor/testdata/go/test.go", false},
				{"processor/src/main.go", false},
			},
		},
		{
			name: "wildcard pattern match",
			config: CommentConfig{
				FileIgnorePatterns: []string{"^processor/testdata/.*"},
			},
			relativePath: "",
			tests: []struct {
				file         string
				shouldIgnore bool
			}{
				{"processor/testdata/bash/test.sh", true},
				{"processor/testdata/go/test.go", true},
				{"processor/src/main.go", false},
			},
		},
		{
			name: "all testdata directories",
			config: CommentConfig{
				FileIgnorePatterns: []string{"testdata"},
			},
			relativePath: "",
			tests: []struct {
				file         string
				shouldIgnore bool
			}{
				{"processor/testdata/bash/test.sh", true},
				{"config/testdata/config.json", true},
				{"walker/testdata/mock.go", true},
				{"processor/src/main.go", false},
			},
		},
		{
			name: "absolute paths",
			config: CommentConfig{
				FileIgnorePatterns: []string{"testdata"},
			},
			relativePath: "",
			tests: []struct {
				file         string
				shouldIgnore bool
			}{
				{filepath.Join(repoDir, "processor/testdata/bash/test.sh"), true},
				{filepath.Join(repoDir, "processor/src/main.go"), false},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			oldWd, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(oldWd); err != nil {
					t.Logf("Failed to restore working directory: %v", err)
				}
			}()

			workDir := repoDir
			if tc.relativePath != "" {
				workDir = filepath.Join(repoDir, tc.relativePath)
			}
			if err := os.Chdir(workDir); err != nil {
				t.Fatalf("Failed to change to working directory: %v", err)
			}

			configFile := filepath.Join(workDir, ".nocmt.json")
			data, err := json.MarshalIndent(tc.config, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal config: %v", err)
			}
			if err := os.WriteFile(configFile, data, 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			cfg := New()
			if err := cfg.LoadConfigurations(); err != nil {
				t.Fatalf("Failed to load configurations: %v", err)
			}

			for _, test := range tc.tests {
				got := cfg.ShouldIgnoreFile(test.file)
				if got != test.shouldIgnore {
					t.Errorf("ShouldIgnoreFile(%q) = %v, want %v", test.file, got, test.shouldIgnore)
				}
			}
		})
	}
}

func TestPreCommitHookIgnorePatterns(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping test")
	}

	tmpDir := t.TempDir()

	repoDir := filepath.Join(tmpDir, "repo")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	cmds := []struct {
		name string
		args []string
	}{
		{"git init", []string{"init"}},
		{"git config user.email", []string{"config", "user.email", "test@example.com"}},
		{"git config user.name", []string{"config", "user.name", "Test User"}},
	}

	for _, cmd := range cmds {
		c := exec.Command("git", cmd.args...)
		c.Dir = repoDir
		if err := c.Run(); err != nil {
			t.Fatalf("Failed to run %s: %v", cmd.name, err)
		}
	}

	configContent := CommentConfig{
		FileIgnorePatterns: []string{
			"^processor/testdata/.*",
			"^config/testdata/.*",
		},
	}

	configFile := filepath.Join(repoDir, ".nocmt.json")
	data, err := json.MarshalIndent(configContent, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	testDirs := []string{
		filepath.Join(repoDir, "processor", "testdata", "bash"),
		filepath.Join(repoDir, "config", "testdata"),
		filepath.Join(repoDir, "src"),
	}

	for _, dir := range testDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	testFiles := map[string]string{
		"processor/testdata/bash/test.sh":  "#!/bin/bash\n# This is a test comment\necho 'Hello'",
		"config/testdata/config.test.json": "// This is a config test file\n{\n  \"test\": true\n}",
		"src/main.go":                      "package main\n\n// This is a comment\nfunc main() {\n  // Another comment\n  println(\"Hello\")\n}",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(repoDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	stageCmd := exec.Command("git", "add", ".")
	stageCmd.Dir = repoDir
	if err := stageCmd.Run(); err != nil {
		t.Fatalf("Failed to stage files: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Failed to restore working directory: %v", err)
		}
	}()

	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	cfg := New()
	if err := cfg.LoadConfigurations(); err != nil {
		t.Fatalf("Failed to load configurations: %v", err)
	}

	stagedCmd := exec.Command("git", "diff", "--name-only", "--cached")
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		t.Fatalf("Failed to get staged files: %v", err)
	}

	stagedFiles := strings.Split(strings.TrimSpace(string(stagedOutput)), "\n")

	var filesToProcess []string
	var ignoredFiles []string
	for _, file := range stagedFiles {
		if cfg.ShouldIgnoreFile(file) {
			ignoredFiles = append(ignoredFiles, file)
		} else {
			filesToProcess = append(filesToProcess, file)
		}
	}

	expectedIgnored := []string{
		"processor/testdata/bash/test.sh",
		"config/testdata/config.test.json",
	}
	expectedToProcess := []string{
		"src/main.go",
		".nocmt.json",
	}

	sort.Strings(ignoredFiles)
	sort.Strings(filesToProcess)
	sort.Strings(expectedIgnored)
	sort.Strings(expectedToProcess)

	if !reflect.DeepEqual(ignoredFiles, expectedIgnored) {
		t.Errorf("Ignored files mismatch.\nGot: %v\nWant: %v", ignoredFiles, expectedIgnored)
	}

	if !reflect.DeepEqual(filesToProcess, expectedToProcess) {
		t.Errorf("Files to process mismatch.\nGot: %v\nWant: %v", filesToProcess, expectedToProcess)
	}
}
