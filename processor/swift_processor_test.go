package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwiftProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewSwiftProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/swift/original.swift", "../testdata/swift/expected.swift")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewSwiftProcessor(false)
		input := `// Regular comment
import Foundation // Import comment
class MyClass { /* Block comment */ }`
		expected := `import Foundation
class MyClass { /* Block comment */ }`
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestSwiftProcessorGetLanguageName(t *testing.T) {
	processor := NewSwiftProcessor(false)
	assert.Equal(t, "swift", processor.GetLanguageName())
}

func TestSwiftProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewSwiftProcessor(true)
	processorWithoutDirectives := NewSwiftProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsSwiftDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"SwiftLintDirective", "// swiftlint:disable line_length", true},
		{"SourceryDirective", "// sourcery: AutoMockable", true},
		{"TODODirective", "// TODO: Implement this", true},
		{"FIXMEDirective", "// FIXME: Fix this bug", true},
		{"MARKDirective", "// MARK: - Section", true},
		{"WARNINGDirective", "// WARNING: Deprecated", true},
		{"NOTEDirective", "// NOTE: Important note", true},
		{"AttributeDirective", "// @available(iOS 13.0, *)", true},
		{"SpacedDirective", "  // TODO: Fix this  ", true},
		{"InlineDirective", "func test() { // swiftlint:disable:next force_cast", true},
		{"RegularLineComment", "// This is a comment", false},
		{"DocumentationComment", "/// This is a doc comment", false},
		{"BlockComment", "/* This is a block comment */", false},
		{"CodeLine", "let x = 5", false},
		{"EmptyLine", "", false},
		{"CommentWithAtSymbol", "// @override in comment", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isSwiftDirective(tt.line))
		})
	}
}

func TestSwiftProcessor_SingleLineComments(t *testing.T) {
	processor := NewSwiftProcessor(false)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "LineCommentOnly",
			input: `// This comment should be removed
let x = 5`,
			expected: `let x = 5`,
		},
		{
			name:     "InlineComment",
			input:    `let x = 5 // This comment should be removed`,
			expected: `let x = 5`,
		},
		{
			name: "MultipleLineComments",
			input: `// First comment
let x = 5 // Inline comment
// Another comment
let y = 10`,
			expected: `let x = 5
let y = 10`,
		},
		{
			name: "PreserveBlockComments",
			input: `/* This block comment should be preserved */
let x = 5 // This line comment should be removed`,
			expected: `/* This block comment should be preserved */
let x = 5`,
		},
		{
			name: "PreserveDocComments",
			input: `/// This documentation comment should be preserved
// This line comment should be removed
func myMethod() {}`,
			expected: `/// This documentation comment should be preserved
func myMethod() {}`,
		},
		{
			name: "PreserveJavaDocStyle",
			input: `/**
 * This documentation comment should be preserved
 */
// This line comment should be removed
func myMethod() {}`,
			expected: `/**
 * This documentation comment should be preserved
 */
func myMethod() {}`,
		},
		{
			name: "MixedComments",
			input: `/// Documentation comment
// Line comment to remove
/* Block comment to preserve */
func test() { // Another line comment to remove
    let x = 5
}`,
			expected: `/// Documentation comment
/* Block comment to preserve */
func test() {
    let x = 5
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
