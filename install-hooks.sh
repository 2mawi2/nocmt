#!/bin/sh
set -e

HOOKS_DIR=".git/hooks"
SCRIPT_DIR=$(dirname "$0")

if [ ! -d ".git" ]; then
    echo "Error: Not a git repository"
    echo "Please run this script from the root of your git repository"
    exit 1
fi

mkdir -p "$HOOKS_DIR"

echo "Installing pre-commit hook..."
cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "Git hooks installed successfully!"
echo ""
echo "The pre-commit hook will automatically run nocmt with --staged flag"
echo "to strip comments from modified lines in your commits."
echo ""
echo "To skip the hook, use: git commit --no-verify"