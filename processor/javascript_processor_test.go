package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJavaScriptProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewJavaScriptProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/javascript/original.js", "../testdata/javascript/expected.js")
	})
	t.Run("WithoutDirectives", func(t *testing.T) {
		processor := NewJavaScriptProcessor(false)
		input := "// @license MIT\nconsole.log(\"hello\"); // comment\n"
		expected := "console.log(\"hello\");\n"
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestJavaScriptProcessorGetLanguageName(t *testing.T) {
	processor := NewJavaScriptProcessor(false)
	assert.Equal(t, "javascript", processor.GetLanguageName())
}

func TestJavaScriptProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewJavaScriptProcessor(true)
	processorWithoutDirectives := NewJavaScriptProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsJSDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"SourceMappingURL", "//# sourceMappingURL=foo.js.map", true},
		{"License single line comment", "// @license MIT", true},
		{"License block comment", "/* @license GPL */", true},
		{"Preserve single line comment", "// @preserve", true},
		{"Preserve block comment", "/* @preserve */", true},
		{"Generic @ in single line", "// @foo", true},
		{"Generic @ in block comment", "/* @bar */", true},
		{"Simple line comment", "// This is a normal comment", false},
		{"Simple block comment", "/* This is a normal comment */", false},
		{"Code line", "console.log(\"hello // @world\");", false},
		{"Empty line", "", false},
		{"Whitespace line", "   ", false},
		{"SourceMappingURL with space", "  //# sourceMappingURL=foo.js.map  ", true},
		{"Equal sign directive", "// =require foo", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isJSDirective(tt.line))
		})
	}
}
