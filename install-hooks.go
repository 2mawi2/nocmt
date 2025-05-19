package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed pre-commit
var preCommitHookContent []byte

func InstallPreCommitHook(verbose bool) error {
	if !isGitRepo() {
		return fmt.Errorf("not a git repository (or any of the parent directories)")
	}

	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find git root directory: %w", err)
	}
	gitRootDir := strings.TrimSpace(string(output))

	hooksDir := filepath.Join(gitRootDir, ".git", "hooks")
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, "pre-commit")

	if _, err := os.Stat(hookPath); err == nil {
		if verbose {
			fmt.Println("Pre-commit hook already exists, backing it up...")
		}
		backupPath := hookPath + ".backup"
		if err := os.Rename(hookPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup existing hook: %w", err)
		}
		if verbose {
			fmt.Printf("Existing hook backed up to %s\n", backupPath)
		}
	}

	err = os.WriteFile(hookPath, preCommitHookContent, 0755)
	if err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	if verbose {
		fmt.Printf("Pre-commit hook installed at %s\n", hookPath)
	}

	return nil
}
