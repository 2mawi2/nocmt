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
		input := `// @formatter:off
package com.example; // Regular comment
public class Main{}`
		expected := `package com.example;
public class Main{}
`
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
		{"FormatterOff", "// @formatter:off", true},
		{"FormatterOn", "// @formatter:on", true},
		{"SuppressWarningsLine", "// @SuppressWarnings(\"all\")", true},
		{"SuppressWarningsBlock", "@SuppressWarnings(\"unused\") // in a comment", true},
		{"CheckstyleSimple", "//CHECKSTYLE:OFF", true},
		{"CheckstyleOn", "// CHECKSTYLE.ON", true},
		{"CheckstyleOffSpecific", "// CHECKSTYLE.OFF: LineLengthCheck", true},
		{"NoCheckstyle", "// NOCHECKSTYLE Javadoc", true},
		{"NoSonar", "// NOSONAR", true},
		{"NoFolint", "// NOFOLINT", true},
		{"RegularComment", "// This is a normal comment", false},
		{"CodeLine", "public void test() { // Some code", false},
		{"EmptyComment", "//", false},
		{"SpacedDirective", "  // @formatter:off  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isJavaDirective(tt.line))
		})
	}
}
