# nocmt - Remove Comments Without Breaking Code

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/2mawi2/nocmt)](https://goreportcard.com/report/github.com/2mawi2/nocmt)

A fast, language-aware tool to remove comments from source code while preserving structure and special directives.

## Why nocmt?

- **Code Cleanup**: Remove explanatory comments before committing/merging
- **Keep Only What Matters**: Preserve important directives like `//go:generate`
- **Selective Processing**: Remove comments only from changed lines in git-staged files
- **Customizable**: Powerful regex patterns to keep specific comments

## Installation

### Homebrew

You can install nocmt using Homebrew:

```bash
brew tap 2mawi2/tap
brew install nocmt
```

```bash
# Using go install
go install github.com/2mawi2/nocmt@latest

# From source
git clone https://github.com/2mawi2/nocmt.git
cd nocmt
go build
```

## Quick Start

```bash
# Remove comments from a file
nocmt -path main.go

# Process directory, keeping compiler directives
nocmt -path ./src -preserve-directives

# Only clean comments in staged git changes
nocmt --staged

# Add a pattern to keep TODO comments
nocmt -add-ignore "TODO"
```

## Supported Languages

- Go
- JavaScript/TypeScript
- Java
- Python
- Rust
- Bash
- CSS

## Usage

```
nocmt -path <filepath|directory> [-preserve-directives] [-dry-run] [-verbose] [-force] [-ignore pattern1,pattern2] [-add-ignore pattern] [-add-ignore-global pattern] [-staged]
```

### Options

- `-path`: Path to the source file or directory to process (required unless using `--staged`)
- `-preserve-directives`: Preserve compiler directives in the output
- `-dry-run`: Don't write changes, just show what would be done
- `-verbose`: Show verbose output
- `-force`: Force processing even if not a git repository
- `-ignore`: Comma-separated list of regex patterns to ignore comments
- `-add-ignore`: Add a regex pattern to the local ignore list
- `-add-ignore-global`: Add a regex pattern to the global ignore list
- `-staged`: Process only staged files and remove comments only from changed lines

## Configuration

nocmt supports both global (`~/.nocmt/config.json`) and local (`.nocmt.json`) configuration:

```json
{
  "ignorePatterns": [
    "TODO",
    "^\\s*//\\s*WHY",
    "#\\d+",
    "TESTPROJECT-\\d+"
  ]
}
```

## Git Integration

### Pre-commit Hook

```bash
# Install git pre-commit hook
nocmt --install-hooks
```

### Using with pre-commit framework

```yaml
repos:
- repo: https://github.com/2mawi2/nocmt
  rev: v1.0.0  # Use the latest version
  hooks:
  - id: nocmt
    name: nocmt
    description: Remove comments from source code
    entry: nocmt --staged --preserve-directives
    language: golang
    files: \.(go|js|ts|java|py|cs|rs|kt|swift|sh|css)$
```

## License

[MIT License](LICENSE) 