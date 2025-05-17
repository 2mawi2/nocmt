package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSSProcessor(t *testing.T) {
	t.Run("FileBasedTestCase", func(t *testing.T) {
		processor := NewCSSProcessor(true)
		RunFileBasedTestCaseVeryLenient(t, processor, "../testdata/css/original.css", "../testdata/css/expected.css")
	})

	t.Run("EmptyInput", func(t *testing.T) {
		processor := NewCSSProcessor(false)
		result, err := processor.StripComments("")
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("OnlyComments", func(t *testing.T) {
		processor := NewCSSProcessor(false)
		input := `/* Comment 1 */
/* Comment 2 */
/* Comment 3 */`
		result, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("UnterminatedComment", func(t *testing.T) {
		processor := NewCSSProcessor(false)
		input := `body {
  color: red;
  /* This comment is not closed
}`
		_, err := processor.StripComments(input)
		assert.Error(t, err)
	})

	t.Run("PreserveDirectives", func(t *testing.T) {
		processor := NewCSSProcessor(true)
		input := `/* Comment */`

		expected := ``

		result, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestCSSProcessorGetLanguageName(t *testing.T) {
	processor := NewCSSProcessor(false)
	assert.Equal(t, "css", processor.GetLanguageName())
}

func TestCSSProcessorPreserveDirectives(t *testing.T) {
	processorWithDirectives := NewCSSProcessor(true)
	assert.True(t, processorWithDirectives.PreserveDirectives())

	processorWithoutDirectives := NewCSSProcessor(false)
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}
