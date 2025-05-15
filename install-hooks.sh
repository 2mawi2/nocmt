#!/bin/sh
#
# Install git hooks for nocmt

# Exit on error
set -e

HOOKS_DIR=".git/hooks"
SCRIPT_DIR=$(dirname "$0")

# Check if .git directory exists
if [ ! -d ".git" ]; then
    echo "Error: Not a git repository"
    echo "Please run this script from the root of your git repository"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Copy pre-commit hook
echo "Installing pre-commit hook..."
cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "Git hooks installed successfully!"
echo ""
echo "The pre-commit hook will automatically run nocmt with --staged flag"
echo "to strip comments from modified lines in your commits."
echo ""
echo "To skip the hook, use: git commit --no-verify" 