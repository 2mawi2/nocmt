package config

import (
	"os"
	"path/filepath"
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
