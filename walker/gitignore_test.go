package walker

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitIgnoreCheckerBasic(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "gitignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	gitignoreContent := []byte("*.log\n/node_modules/\n!important.log\n")
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), gitignoreContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	filesToCreate := []string{
		"file.txt",
		"file.log",
		"important.log",
		"node_modules/module.js",
		"subfolder/file.log",
	}

	for _, file := range filesToCreate {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	checker, err := NewGitIgnoreChecker(tempDir)
	if err != nil {
		t.Fatalf("Failed to create GitIgnoreChecker: %v", err)
	}

	testCases := []struct {
		path     string
		expected bool
	}{
		{"file.txt", false},
		{"file.log", true},
		{"important.log", false},
		{"node_modules/module.js", true},
		{"subfolder/file.log", true},
	}

	for _, tc := range testCases {
		fullPath := filepath.Join(tempDir, tc.path)
		isIgnored := checker.IsIgnored(fullPath)
		if isIgnored != tc.expected {
			t.Errorf("IsIgnored(%s) = %v, expected %v", tc.path, isIgnored, tc.expected)
		}
	}
}

func TestHierarchicalGitIgnoreChecker(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "nested-gitignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	rootGitignore := []byte("*.log\n!important.log\n")
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), rootGitignore, 0644)
	if err != nil {
		t.Fatalf("Failed to create root .gitignore file: %v", err)
	}

	subDir := filepath.Join(tempDir, "subfolder")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subfolder: %v", err)
	}

	subGitignore := []byte("!*.log\n*.txt\n")
	err = os.WriteFile(filepath.Join(subDir, ".gitignore"), subGitignore, 0644)
	if err != nil {
		t.Fatalf("Failed to create subfolder .gitignore file: %v", err)
	}

	filesToCreate := []string{
		"file.txt",
		"file.log",
		"important.log",
		"subfolder/file.log",
		"subfolder/file.txt",
	}

	for _, file := range filesToCreate {
		filePath := filepath.Join(tempDir, file)
		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	t.Logf("Root .gitignore: %s", string(rootGitignore))
	t.Logf("Subfolder .gitignore: %s", string(subGitignore))

	checker, err := NewHierarchicalGitIgnoreChecker(tempDir)
	if err != nil {
		t.Fatalf("Failed to create HierarchicalGitIgnoreChecker: %v", err)
	}

	t.Logf("Found %d gitignore files:", len(checker.gitignoreFiles))
	for dir, ignorer := range checker.gitignoreFiles {
		t.Logf("  - %s", dir)

		for _, file := range filesToCreate {
			relPath, err := filepath.Rel(tempDir, filepath.Join(tempDir, file))
			if err != nil {
				continue
			}

			var dirRelPath string
			if dir == "" {
				dirRelPath = relPath
			} else {
				if strings.HasPrefix(relPath, dir+"/") {
					dirRelPath = relPath[len(dir)+1:]
				} else if relPath == dir {
					dirRelPath = "."
				} else {
					continue
				}
			}

			dirRelPath = filepath.ToSlash(dirRelPath)
			matches, pattern := ignorer.MatchesPathHow(dirRelPath)
			if matches && pattern != nil {
				t.Logf("    - %s matches pattern in %s: \"%s\" (negate: %v)",
					file, dir, pattern.Line, pattern.Negate)
			} else {
				t.Logf("    - %s does NOT match any pattern in %s (relPath: %s)",
					file, dir, dirRelPath)
			}
		}
	}

	t.Logf("Patterns extracted from gitignore files:")
	for dir, patterns := range checker.patterns {
		t.Logf("  Directory %q:", dir)
		for _, pat := range patterns {
			t.Logf("    - Pattern: %q, Negated: %v, Wildcard: %v", pat.pattern, pat.isNegated, pat.isWildcard)
		}
	}

	problemFile := "subfolder/file.log"
	problemFilePath := filepath.Join(tempDir, problemFile)

	t.Logf("Manual testing for %s:", problemFile)
	rootIgnorer := checker.gitignoreFiles[""]
	if rootIgnorer != nil {
		relToRoot, _ := filepath.Rel(tempDir, problemFilePath)
		relToRoot = filepath.ToSlash(relToRoot)
		rootMatches, rootPattern := rootIgnorer.MatchesPathHow(relToRoot)
		if rootMatches && rootPattern != nil {
			t.Logf("  - Matches root pattern: \"%s\" (negate: %v)",
				rootPattern.Line, rootPattern.Negate)
		} else {
			t.Logf("  - Does NOT match any root pattern (relPath: %s)", relToRoot)
		}
	}

	subFolderIgnorer := checker.gitignoreFiles["subfolder"]
	if subFolderIgnorer != nil {
		relToSubfolder := "file.log"
		subFolderMatches, subFolderPattern := subFolderIgnorer.MatchesPathHow(relToSubfolder)
		if subFolderMatches && subFolderPattern != nil {
			t.Logf("  - Matches subfolder pattern: \"%s\" (negate: %v)",
				subFolderPattern.Line, subFolderPattern.Negate)
		} else {
			t.Logf("  - Does NOT match any subfolder pattern (relPath: %s)", relToSubfolder)
		}
	}

	relProblemPath, _ := filepath.Rel(tempDir, problemFilePath)
	relProblemPath = filepath.ToSlash(relProblemPath)
	dirPath := filepath.Dir(relProblemPath)
	t.Logf("  - Problem file: %s, Dir path: %s", relProblemPath, dirPath)
	if patterns, ok := checker.patterns[dirPath]; ok {
		t.Logf("  - Found patterns for directory %s:", dirPath)
		for _, pat := range patterns {
			t.Logf("    - Pattern: %q, Negated: %v", pat.pattern, pat.isNegated)
			if pat.isNegated && pat.pattern == "*.log" {
				t.Logf("    - Should match special case and NOT be ignored")
			}
		}
	} else {
		t.Logf("  - No patterns found for directory %s", dirPath)
	}

	testCases := []struct {
		path     string
		expected bool
	}{
		{"file.txt", false},
		{"file.log", true},
		{"important.log", false},
		{"subfolder/file.log", false},
		{"subfolder/file.txt", true},
	}

	for _, tc := range testCases {
		fullPath := filepath.Join(tempDir, tc.path)
		isIgnored := checker.IsIgnored(fullPath)

		t.Logf("Path: %s, Expected ignored: %v, Actually ignored: %v", tc.path, tc.expected, isIgnored)

		if isIgnored != tc.expected {
			t.Errorf("IsIgnored(%s) = %v, expected %v", tc.path, isIgnored, tc.expected)
		}
	}
}

func TestDefaultIgnorePatternsBehavior(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "default-ignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	filesToIgnore := []string{
		"node_modules/module.js",
		".vscode/settings.json",
		".git/config",
		"build/output.js",
		"file.tmp",
		"file.swp",
		"Thumbs.db",
		".DS_Store",
		"package-lock.json",
	}

	filesToProcess := []string{
		"src/main.go",
		"README.md",
		"index.html",
		"assets/image.png",
	}

	for _, file := range append(filesToIgnore, filesToProcess...) {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		err = os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	checker, err := NewHierarchicalGitIgnoreChecker(tempDir)
	if err != nil {
		t.Fatalf("Failed to create HierarchicalGitIgnoreChecker: %v", err)
	}

	for _, file := range filesToIgnore {
		fullPath := filepath.Join(tempDir, file)
		isIgnored := checker.IsIgnored(fullPath)
		if !isIgnored {
			t.Errorf("Expected %s to be ignored by default, but it wasn't", file)
		}
	}

	for _, file := range filesToProcess {
		fullPath := filepath.Join(tempDir, file)
		isIgnored := checker.IsIgnored(fullPath)
		if isIgnored {
			t.Errorf("Expected %s not to be ignored by default, but it was", file)
		}
	}

	customGitignore := []byte("*.md\n")
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), customGitignore, 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	checker, err = NewHierarchicalGitIgnoreChecker(tempDir)
	if err != nil {
		t.Fatalf("Failed to create HierarchicalGitIgnoreChecker: %v", err)
	}

	for _, file := range filesToIgnore {
		fullPath := filepath.Join(tempDir, file)
		isIgnored := checker.IsIgnored(fullPath)
		if !isIgnored {
			t.Errorf("With custom .gitignore, expected %s to be ignored by default, but it wasn't", file)
		}
	}

	readmePath := filepath.Join(tempDir, "README.md")
	isIgnored := checker.IsIgnored(readmePath)
	if !isIgnored {
		t.Errorf("Expected README.md to be ignored by custom .gitignore, but it wasn't")
	}

	for _, file := range filesToProcess {
		if file == "README.md" {
			continue
		}
		fullPath := filepath.Join(tempDir, file)
		isIgnored := checker.IsIgnored(fullPath)
		if isIgnored {
			t.Errorf("With custom .gitignore, expected %s not to be ignored, but it was", file)
		}
	}
}
