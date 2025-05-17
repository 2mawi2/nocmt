package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRustProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewRustProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/rust/original.rs", "../testdata/rust/expected.rs")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewRustProcessor(false)
		input := `// Regular comment
#![allow(unused_variables)] // A directive
fn main() { /* Another comment */ }`
		expected := `fn main() { }
` 
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestRustProcessorGetLanguageName(t *testing.T) {
	processor := NewRustProcessor(false) 
	assert.Equal(t, "rust", processor.GetLanguageName())
}

func TestRustProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewRustProcessor(true)
	processorWithoutDirectives := NewRustProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsRustDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"OuterAttribute", "#[derive(Debug)]", true},
		{"InnerAttribute", "#![allow(dead_code)]", true},
		{"SpacedOuterAttribute", "  #[cfg(test)]  ", true},
		{"SpacedInnerAttribute", "  #![feature(custom_derive)]  ", true},
		{"RegularLineComment", "// This is a comment", false},
		{"DocCommentOuter", "/// This is an outer doc comment", false}, 
		{"DocCommentInner", "//! This is an inner doc comment", false}, 
		{"CommentWithHashBracket", "// #[not_an_attribute]", false},
		{"CodeLine", "let x = 5; // #[attribute_in_comment]", false},
		{"EmptyLine", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isRustDirective(tt.line))
		})
	}
}

