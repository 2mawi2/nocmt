
I'll examine the existing processors to understand the patterns and create a concise PRD for a PHP processor implementation.

Read file: internal/processor/processor.go

Read file: internal/processor/common.go

Read file: internal/processor/javascript_processor.go

Read file: internal/processor/python_single_processor.go

Ran tool
Now I'll look at existing test files to understand testing patterns:

Read file: internal/processor/javascript_processor_test.go
Now I have enough information to create a comprehensive PRD for a PHP processor. Let me generate this based on the patterns I've observed.

# Product Requirements Document: PHP Processor Implementation

## Overview

Implement a PHP language processor for the `nocmt` tool that removes comments from PHP source code while preserving important directives and maintaining proper code structure.

## Core Requirements

### 1. **Interface Compliance**
- Implement the `LanguageProcessor` interface:
  - `StripComments(source string) (string, error)`
  - `GetLanguageName() string` - return "php"
  - `PreserveDirectives() bool`
  - `SetCommentConfig(cfg *config.Config)`

### 2. **Comment Types Support**
PHP supports three comment syntaxes that must be handled:
- **Single-line comments**: `// comment` and `# comment`
- **Multi-line comments**: `/* multi-line comment */`
- **DocBlock comments**: `/** documentation */`

### 3. **Directive Preservation**
When `preserveDirectives` is enabled, preserve:
- **Shebang lines**: `#!/usr/bin/env php`, `#!/usr/bin/php`
- **PHP directives**: `<?php`, `?>`
- **Magic comments**: 
  - `# @codingStandardsIgnoreStart`
  - `# @codingStandardsIgnoreEnd`
  - `# @phan-*` (Phan static analyzer)
  - `# @phpstan-*` (PHPStan static analyzer)
  - `# @psalm-*` (Psalm static analyzer)

### 4. **File Extension Support**
Register support for: `.php`, `.phtml`, `.php3`, `.php4`, `.php5`, `.phps`

### 5. **Tree-sitter Integration**
- Use `github.com/smacker/go-tree-sitter/php` grammar
- Embed `SingleLineCoreProcessor` for consistent behavior
- Handle comment node types: `comment`, `shell_comment_line`

## Implementation Structure

### Files to Create
1. `internal/processor/php_processor.go` - Main processor
2. `internal/processor/php_processor_test.go` - Unit tests  
3. `testdata/php/original.php` - Test input file
4. `testdata/php/expected.php` - Expected output file

### Code Structure Pattern
```go
type PHPProcessor struct {
    *SingleLineCoreProcessor
}

func NewPHPProcessor(preserveDirectivesFlag bool) *PHPProcessor {
    // Use SingleLineCoreProcessor pattern like other processors
}

func isPHPCommentNode(node *sitter.Node, sourceText string) bool {
    // Detect PHP comment nodes
}

func isPHPDirective(line string) bool {
    // Check for PHP-specific directives to preserve
}
```

## Technical Specifications

### Dependencies
- Add `github.com/smacker/go-tree-sitter/php` to `go.mod`
- Follow existing processor patterns in codebase

### Registration Requirements
- Add to `ProcessorFactory` in `processor.go`
- Map file extensions in `GetProcessorByExtension()`
- Include constructor in `NewProcessorFactory()`

### Testing Requirements
- File-based tests using testdata
- Unit tests for directive detection
- Language name and preserve directives flag tests
- Edge cases: nested comments, mixed comment styles

## Success Criteria

### Must Have
- ✅ Removes all three PHP comment types correctly
- ✅ Preserves original trailing newlines
- ✅ Handles directive preservation when enabled
- ✅ Passes all unit tests
- ✅ Works with existing CLI options
- ✅ Maintains consistent performance

## Edge Cases to Handle
1. Comments within PHP strings (should not be removed)
2. Heredoc/Nowdoc syntax with comment-like content
3. Mixed comment styles in same file
4. Comments at EOF without trailing newlines
5. Unicode content in comments

## Future Considerations
- Support for additional PHP-specific static analysis tools
- Framework-specific comment directives (Laravel, Symfony)
- Integration with PHP-CS-Fixer style comments