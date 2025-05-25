package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCppProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewCppProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../../testdata/cpp/original.cpp", "../../testdata/cpp/expected.cpp")
	})
	t.Run("WithoutDirectives", func(t *testing.T) {
		processor := NewCppProcessor(false)
		input := `// Regular comment
#include <iostream>
int main() { /* Block comment */ }`
		expected := `#include <iostream>
int main() { /* Block comment */ }`
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestCppProcessor_GetLanguageName(t *testing.T) {
	processor := NewCppProcessor(false)
	assert.Equal(t, "cpp", processor.GetLanguageName())
}

func TestCppProcessor_PreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewCppProcessor(true)
	processorWithoutDirectives := NewCppProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsCppDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"TODO", "// TODO: Fix this issue", true},
		{"FIXME", "// FIXME: Memory leak", true},
		{"NOTE", "// NOTE: This is important", true},
		{"HACK", "// HACK: Temporary solution", true},
		{"XXX", "// XXX: Review this code", true},
		{"BUG", "// BUG: Known issue", true},
		{"WARNING", "// WARNING: Deprecated", true},
		{"Pragma", "// pragma once", true},
		{"PragmaHash", "// #pragma pack(1)", true},
		{"SpacedTODO", "  // TODO: With spaces", true},
		{"RegularComment", "// This is a regular comment", false},
		{"BlockComment", "/* This is a block comment */", false},
		{"CodeLine", "int x = 5;", false},
		{"EmptyLine", "", false},
		{"IncludeDirective", "#include <iostream>", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isCppDirective(tt.line))
		})
	}
}

func TestCppProcessor_SingleLineComments(t *testing.T) {
	processor := NewCppProcessor(false)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "LineCommentOnly",
			input: `// This comment should be removed
int x = 5;`,
			expected: `int x = 5;`,
		},
		{
			name:     "InlineComment",
			input:    `int x = 5; // This comment should be removed`,
			expected: `int x = 5;`,
		},
		{
			name: "MultipleLineComments",
			input: `// First comment
int x = 5; // Inline comment
// Another comment
int y = 10;`,
			expected: `int x = 5;
int y = 10;`,
		},
		{
			name: "PreserveBlockComments",
			input: `/* This block comment should be preserved */
int x = 5; // This line comment should be removed`,
			expected: `/* This block comment should be preserved */
int x = 5;`,
		},
		{
			name: "PreserveDocComments",
			input: `/**
 * This documentation comment should be preserved
 * @param value The input value
 */
// This line comment should be removed
int getValue(int value) {}`,
			expected: `/**
 * This documentation comment should be preserved
 * @param value The input value
 */
int getValue(int value) {}`,
		},
		{
			name: "PreserveIncludeStatements",
			input: `#include <iostream> // This line comment should be removed
#include <vector>
// Regular comment to remove
int main() {}`,
			expected: `#include <iostream>
#include <vector>
int main() {}`,
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

func TestCppProcessor_DirectivePreservation(t *testing.T) {
	processor := NewCppProcessor(true)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "PreserveTODO",
			input: `// TODO: Implement this feature
// Regular comment to remove
int main() {}`,
			expected: `// TODO: Implement this feature
int main() {}`,
		},
		{
			name: "PreserveFIXME",
			input: `// FIXME: Memory leak here
// Another regular comment
void function() {}`,
			expected: `// FIXME: Memory leak here
void function() {}`,
		},
		{
			name: "PreserveMultipleDirectives",
			input: `// TODO: Feature A
// FIXME: Bug B
// Regular comment
// NOTE: Important info
int x = 5;`,
			expected: `// TODO: Feature A
// FIXME: Bug B
// NOTE: Important info
int x = 5;`,
		},
		{
			name: "PreservePragma",
			input: `// pragma once
// #pragma pack(1)
// Regular comment
class MyClass {};`,
			expected: `// pragma once
// #pragma pack(1)
class MyClass {};`,
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

func TestCppProcessor_ComplexScenarios(t *testing.T) {
	processor := NewCppProcessor(true)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "MixedCommentTypes",
			input: `#include <iostream>
// TODO: Optimize this
/* Block comment
   multiline */
class Example {
    int value;

    /**
     * Documentation comment
     */
    void method() {
        // FIXME: Handle edge case
        std::cout << "Hello";
    }
};`,
			expected: `#include <iostream>
// TODO: Optimize this
/* Block comment
   multiline */
class Example {
    int value;

    /**
     * Documentation comment
     */
    void method() {
        // FIXME: Handle edge case
        std::cout << "Hello";
    }
};`,
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
