package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJavaProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewJavaProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/java/original.java", "../testdata/java/expected.java")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewJavaProcessor(false)
		input := `// Regular comment
@Override // An annotation
public class MyClass { /* Block comment */ }`
		expected := `@Override
public class MyClass { /* Block comment */ }`
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestJavaProcessor_Comprehensive(t *testing.T) {
	processor := NewJavaProcessor(true)
	RunFileBasedTestCaseNormalized(t, processor, "../testdata/java/comprehensive_original.java", "../testdata/java/comprehensive_expected.java")
}

func TestJavaProcessor_EdgeCases(t *testing.T) {
	processor := NewJavaProcessor(true)
	RunFileBasedTestCaseNormalized(t, processor, "../testdata/java/edge_cases_original.java", "../testdata/java/edge_cases_expected.java")
}

func TestJavaProcessorGetLanguageName(t *testing.T) {
	processor := NewJavaProcessor(false)
	assert.Equal(t, "java", processor.GetLanguageName())
}

func TestJavaProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewJavaProcessor(true)
	processorWithoutDirectives := NewJavaProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsJavaDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"FormatterDirective", "// @formatter:off", true},
		{"SuppressWarningsDirective", "// @SuppressWarnings", true},
		{"CheckstyleDirective", "//CHECKSTYLE:OFF", true},
		{"SuppressWarningsAnnotation", "@SuppressWarnings(\"unchecked\")", true},
		{"CheckstyleOff", "CHECKSTYLE.OFF", true},
		{"CheckstyleOn", "CHECKSTYLE.ON", true},
		{"NoCheckstyle", "NOCHECKSTYLE", true},
		{"NoSonar", "NOSONAR", true},
		{"NoFollowLint", "NOFOLINT", true},
		{"SpacedDirective", "  // @formatter:on  ", true},
		{"LineWithDirective", "System.out.println(); // NOSONAR", true},
		{"RegularLineComment", "// This is a comment", false},
		{"DocumentationComment", "/** This is a doc comment */", false},
		{"BlockComment", "/* This is a block comment */", false},
		{"CodeLine", "int x = 5;", false},
		{"EmptyLine", "", false},
		{"CommentWithAtSymbol", "// @Override in comment", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isJavaDirective(tt.line))
		})
	}
}

func TestJavaProcessor_SingleLineComments(t *testing.T) {
	processor := NewJavaProcessor(false)

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
 */
// This line comment should be removed
public void myMethod() {}`,
			expected: `/**
 * This documentation comment should be preserved
 */
public void myMethod() {}`,
		},
		{
			name: "PreserveAnnotations",
			input: `@Override // This line comment should be removed
@SuppressWarnings("unchecked") // And this one too
public void myMethod() {}`,
			expected: `@Override
@SuppressWarnings("unchecked")
public void myMethod() {}`,
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
