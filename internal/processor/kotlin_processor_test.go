package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKotlinProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewKotlinProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../../testdata/kotlin/original.kt", "../../testdata/kotlin/expected.kt")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewKotlinProcessor(false)
		input := `// Regular comment
@Entity // An annotation
class MyClass { /* Block comment */ }`
		expected := `@Entity
class MyClass { /* Block comment */ }`
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestKotlinProcessorGetLanguageName(t *testing.T) {
	processor := NewKotlinProcessor(false)
	assert.Equal(t, "kotlin", processor.GetLanguageName())
}

func TestKotlinProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewKotlinProcessor(true)
	processorWithoutDirectives := NewKotlinProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsKotlinDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"EntityAnnotation", "@Entity", true},
		{"SuppressAnnotation", "@Suppress(\"UNCHECKED_CAST\")", true},
		{"JvmStaticAnnotation", "@JvmStatic", true},
		{"JvmFieldAnnotation", "@JvmField", true},
		{"JvmNameAnnotation", "@JvmName(\"customName\")", true},
		{"JvmOverloadsAnnotation", "@JvmOverloads", true},
		{"DeprecatedAnnotation", "@Deprecated(\"Use new method\")", true},
		{"TargetAnnotation", "@Target(AnnotationTarget.CLASS)", true},
		{"RetentionAnnotation", "@Retention(AnnotationRetention.RUNTIME)", true},
		{"ComponentAnnotation", "@Component", true},
		{"ServiceAnnotation", "@Service", true},
		{"RepositoryAnnotation", "@Repository", true},
		{"ControllerAnnotation", "@Controller", true},
		{"ShebangLine", "#!/usr/bin/env kotlin", true},
		{"SpacedAnnotation", "  @Entity  ", true},
		{"LineWithAnnotation", "class MyClass @Entity", true},
		{"RegularLineComment", "// This is a comment", false},
		{"DocumentationComment", "/** This is a doc comment */", false},
		{"BlockComment", "/* This is a block comment */", false},
		{"CodeLine", "val x = 5", false},
		{"EmptyLine", "", false},
		{"CommentWithAtSymbol", "// @Entity in comment", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isKotlinDirective(tt.line))
		})
	}
}

func TestKotlinProcessor_SingleLineComments(t *testing.T) {
	processor := NewKotlinProcessor(false)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "LineCommentOnly",
			input: `// This comment should be removed
val x = 5`,
			expected: `val x = 5`,
		},
		{
			name:     "InlineComment",
			input:    `val x = 5 // This comment should be removed`,
			expected: `val x = 5`,
		},
		{
			name: "MultipleLineComments",
			input: `// First comment
val x = 5 // Inline comment
// Another comment
val y = 10`,
			expected: `val x = 5
val y = 10`,
		},
		{
			name: "PreserveBlockComments",
			input: `/* This block comment should be preserved */
val x = 5 // This line comment should be removed`,
			expected: `/* This block comment should be preserved */
val x = 5`,
		},
		{
			name: "PreserveDocComments",
			input: `/**
 * This documentation comment should be preserved
 */
// This line comment should be removed
fun myFunction() {}`,
			expected: `/**
 * This documentation comment should be preserved
 */
fun myFunction() {}`,
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
