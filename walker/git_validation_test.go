package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-validation-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	isGit := IsGitRepository(tempDir)
	if isGit {
		t.Errorf("Expected non-git directory to return false, but got true")
	}

	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	isGit = IsGitRepository(tempDir)
	if !isGit {
		t.Errorf("Expected git directory to return true, but got false")
	}

	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	isGit = IsGitRepository(subDir)
	if !isGit {
		t.Errorf("Expected subdirectory of git repository to return true, but got false")
	}
}

func TestConfirmNonGitUsage(t *testing.T) {
	confirmed := confirmNonGitUsage("y\n")
	if !confirmed {
		t.Errorf("Expected 'y' to confirm, but got rejection")
	}

	confirmed = confirmNonGitUsage("yes\n")
	if !confirmed {
		t.Errorf("Expected 'yes' to confirm, but got rejection")
	}

	confirmed = confirmNonGitUsage("n\n")
	if confirmed {
		t.Errorf("Expected 'n' to reject, but got confirmation")
	}

	confirmed = confirmNonGitUsage("no\n")
	if confirmed {
		t.Errorf("Expected 'no' to reject, but got confirmation")
	}

	confirmed = confirmNonGitUsage("invalid\ny\n")
	if !confirmed {
		t.Errorf("Expected 'invalid\\ny' to confirm, but got rejection")
	}
}

func TestValidateGitRepository(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "git-validation-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	result := ValidateGitRepository(tempDir, true)
	if !result {
		t.Errorf("Expected ValidateGitRepository with force=true to return true, but got false")
	}

	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	result = ValidateGitRepository(tempDir, false)
	if !result {
		t.Errorf("Expected ValidateGitRepository for git repository to return true, but got false")
	}

}
