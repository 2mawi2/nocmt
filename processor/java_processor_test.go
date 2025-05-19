package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJavaProcessor_FileBased(t *testing.T) {
	processor := NewJavaProcessor(false)
	RunFileBasedTestCaseNormalized(t, processor, "../testdata/java/original.java", "../testdata/java/expected.java")
}

func TestJavaProcessorGetLanguageName(t *testing.T) {
	processor := NewJavaProcessor(false)
	assert.Equal(t, "java", processor.GetLanguageName())
}
