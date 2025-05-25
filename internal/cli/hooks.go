package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const preCommitHookContent = `#!/bin/sh
#
# pre-commit hook that runs nocmt on staged files
#

# Exit on error
set -e

# Check if nocmt is in PATH, otherwise use the local binary
NOCMT_CMD="./nocmt"
if [ ! -x "$NOCMT_CMD" ]; then
    NOCMT_CMD="nocmt"
    if ! command -v $NOCMT_CMD >/dev/null 2>&1; then
        echo "Error: nocmt not found in PATH or current directory"
        echo "Please build the nocmt binary first or add it to your PATH"
        exit 1
    fi
fi

echo "Running nocmt to remove comments from staged files..."

# Use the --staged flag to process all staged files at once
$NOCMT_CMD --staged --verbose

# Exit with success status
exit 0
`

func InstallPreCommitHook(verbose bool) error {
	if !IsGitRepo() {
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

	err = os.WriteFile(hookPath, []byte(preCommitHookContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write hook file: %w", err)
	}

	if verbose {
		fmt.Printf("Pre-commit hook installed at %s\n", hookPath)
	}

	return nil
}
