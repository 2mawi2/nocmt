package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewGoProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/go/original.go", "../testdata/go/expected.go")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewGoProcessor(false)
		input := `//go:generate echo "hello"
package main // comment
func main(){}`
		expected := `package main
func main(){}
`
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("No_Line_Artifacts_Strict", func(t *testing.T) {
		processor := NewGoProcessor(true)
		RunFileBasedTestCase(t, processor, "../testdata/go/original_no_line_artifacts.go", "../testdata/go/expected_no_line_artifacts.go")
	})
}

func TestGoProcessorGetLanguageName(t *testing.T) {
	processor := NewGoProcessor(false)
	assert.Equal(t, "go", processor.GetLanguageName())
}

func TestGoProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewGoProcessor(true)
	processorWithoutDirectives := NewGoProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsGoDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"GoGenerate", "//go:generate echo hello", true},
		{"GoBuild", "//go:build linux", true},
		{"PlusBuild", "// +build ignore", true},
		{"GoEmbed", "//go:embed file.txt", true},
		{"CgoInclude", "// #include <stdio.h>", true},
		{"CgoLdflags", "// #cgo LDFLAGS: -lm", true},
		{"CgoAnything", "// #cgo CFLAGS: -DPNG_DEBUG=1", true},
		{"StandardComment", "// This is a standard comment", false},
		{"CommentWithGoKeyword", "// The go keyword is here", false},
		{"CommentWithBuildKeyword", "// Build this thing", false},
		{"EmptyComment", "//", false},
		{"SpacedGoDirective", "  //go:generate echo spaced", true},
		{"NoSpaceGoDirective", "//go:generate", true},
		{"MalformedGoDirective", "//go: generate", true},
		{"NotADirective", "func main() { //go:fmt off }", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, checkGoDirective(tt.line))
		})
	}
}
