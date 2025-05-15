package walker

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func IsGitRepository(dir string) bool {
	gitDir := filepath.Join(dir, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		return true
	}

	parentDir := filepath.Dir(dir)
	if parentDir == dir {
		return false
	}

	return IsGitRepository(parentDir)
}

func confirmNonGitUsage(input string) bool {
	if input != "" {
		scanner := bufio.NewScanner(strings.NewReader(input))
		for scanner.Scan() {
			response := strings.ToLower(strings.TrimSpace(scanner.Text()))
			switch response {
			case "y", "yes":
				return true
			case "n", "no":
				return false
			default:
				fmt.Println("Please enter 'y' or 'n':")
			}
		}
		return false
	}

	fmt.Println("WARNING: The target directory is not a Git repository or a subdirectory of one.")
	fmt.Println("Running nocmt on non-Git repositories is not recommended, as it may modify files")
	fmt.Println("that you don't want to modify and there's no version control to revert changes.")
	fmt.Println("Do you want to continue anyway? (y/n):")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))
		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Please enter 'y' or 'n':")
		}
	}

	return false
}

func ValidateGitRepository(dir string, force bool) bool {
	if IsGitRepository(dir) {
		return true
	}

	if force {
		fmt.Println("WARNING: Processing a non-Git repository (force mode enabled)")
		return true
	}

	return confirmNonGitUsage("")
}
