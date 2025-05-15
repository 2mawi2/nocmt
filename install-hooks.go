package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

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

	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	hook := "#!/bin/sh\n\n"
	hook += "# Run nocmt on staged files\n"
	hook += execPath + " --staged\n\n"
	hook += "exit 0\n"

	err = os.WriteFile(hookPath, []byte(hook), 0755)
	if err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	if verbose {
		fmt.Printf("Pre-commit hook installed at %s\n", hookPath)
	}

	return nil
}
