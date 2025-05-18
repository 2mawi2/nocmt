package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTSXFileBased(t *testing.T) {
	processor := NewTSXProcessor(true)
	RunFileBasedTestCaseNormalized(t, processor, "../testdata/tsx/original.tsx", "../testdata/tsx/expected.tsx")
}

func TestTSXProcessorGetLanguageName(t *testing.T) {
	processor := NewTSXProcessor(false)
	assert.Equal(t, "tsx", processor.GetLanguageName())
}

func TestTSXProcessorPreserveDirectives(t *testing.T) {
	processorWithDirectives := NewTSXProcessor(true)
	processorWithoutDirectives := NewTSXProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}
