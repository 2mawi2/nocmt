# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

nocmt is a language-aware comment removal tool designed to clean up AI-generated comments from code. It supports selective comment removal from git-staged changes and processes multiple programming languages using Tree-sitter parsers.

## Development Commands

### Build and Run
- `just build` or `go build -o nocmt` - Build the binary
- `just run [args]` or `go run main.go [args]` - Run with arguments
- `just clean` - Remove build artifacts

### Testing
- `just test` or `go test ./...` - Run all tests
- `go test ./processor` - Run processor-specific tests
- Individual test files follow pattern `*_test.go` using testify/assert

### Code Quality
- `just lint` - Run golangci-lint (required for CI)
- `gofmt` formatting is mandatory

### Benchmarking
- `just bench [run|quick|compare|history]` - Run benchmarks via benchmark.sh

## Architecture

### Core Components

**Processors (`processor/`)**
- Each language implements `LanguageProcessor` interface
- Uses Tree-sitter grammars for parsing
- Base functionality in `common.go` with `BaseProcessor`
- Factory pattern in `processor.go` manages processor creation and registration

**Walker (`walker/`)**
- File system traversal and git integration
- Gitignore support and repository validation
- Selective processing of staged changes only

**Config (`config/`)**
- Global (`~/.nocmt/config.json`) and local (`.nocmt.json`) configuration
- Regex patterns for preserving specific comments
- File ignore patterns for skipping files

### Key Processing Modes
1. **Staged files** (default) - Process only git-staged changes on modified lines
2. **Single file** - Process specific file completely  
3. **Directory** - Process all supported files recursively

### Adding New Language Support
1. Create processor in `processor/` implementing `LanguageProcessor`
2. Add Tree-sitter grammar dependency to `go.mod`
3. Register in `ProcessorFactory` with file extension mapping
4. Add comprehensive tests in `*_processor_test.go`

## Development Patterns

- Go 1.22+ required (see go.mod)
- Error wrapping with `fmt.Errorf("...: %w", err)`
- File-based tests using testdata/ directory with original/expected pairs
- Processor tests use `RunFileBasedTestCase()` helper
- Preserve compiler directives by default (e.g., `//go:generate`)
- Use `--dry-run` flag for testing changes safely
- Always run tests after making changes
- Verify builds complete successfully
- Follow existing code patterns and conventions⏎    

## Git Integration

Tool designed for git workflows:
- Processes staged files by default via `git diff --cached`
- Re-stages processed files automatically
- Supports pre-commit hook installation via `nocmt install`
- Selective comment removal only on modified lines (not entire file)

## GitHub CLI Access
- Use `gh` command to interact with GitHub (issues, PRs, comments)
- GitHub credentials are available from environment variables
- Can comment on issues, create PRs, and manage repository interactions
