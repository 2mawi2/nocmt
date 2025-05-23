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
		{"OverrideAnnotation", "@Override", true},
		{"SuppressWarningsAnnotation", "@SuppressWarnings(\"unchecked\")", true},
		{"ComponentAnnotation", "@Component", true},
		{"ServiceAnnotation", "@Service", true},
		{"RepositoryAnnotation", "@Repository", true},
		{"ControllerAnnotation", "@Controller", true},
		{"InjectAnnotation", "@Inject", true},
		{"DeprecatedAnnotation", "@Deprecated(\"Use new method\")", true},
		{"PreDestroyAnnotation", "@PreDestroy", true},
		{"PostConstructAnnotation", "@PostConstruct", true},
		{"EntityAnnotation", "@Entity", true},
		{"TableAnnotation", "@Table(name=\"users\")", true},
		{"ShebangLine", "#!/usr/bin/env java", true},
		{"SpacedAnnotation", "  @Override  ", true},
		{"LineWithAnnotation", "public class MyClass @Entity", true},
		{"RegularLineComment", "// This is a comment", false},
		{"DocumentationComment", "/** This is a doc comment */", false},
		{"BlockComment", "/* This is a block comment */", false},
		{"CodeLine", "String x = \"hello\";", false},
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
String x = "hello";`,
			expected: `String x = "hello";`,
		},
		{
			name:     "InlineComment",
			input:    `String x = "hello"; // This comment should be removed`,
			expected: `String x = "hello";`,
		},
		{
			name: "MultipleLineComments",
			input: `// First comment
String x = "hello"; // Inline comment
// Another comment
int y = 10;`,
			expected: `String x = "hello";
int y = 10;`,
		},
		{
			name: "PreserveBlockComments",
			input: `/* This block comment should be preserved */
String x = "hello"; // This line comment should be removed`,
			expected: `/* This block comment should be preserved */
String x = "hello";`,
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
			input: `@Override // This comment should be removed
@SuppressWarnings("unchecked") // Another comment to remove
public String toString() {
    return "test"; // Inline comment to remove
}`,
			expected: `@Override
@SuppressWarnings("unchecked")
public String toString() {
    return "test";
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