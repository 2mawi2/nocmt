package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPythonProcessor_FileBased(t *testing.T) {
	t.Run("Default_PreserveDirectives", func(t *testing.T) {
		pyProc := NewPythonSingleProcessor(true)
		RunFileBasedTestCaseNormalized(t, pyProc, "../../testdata/python/original.py", "../../testdata/python/expected.py")
	})
	t.Run("RemoveAll_NoDirectives", func(t *testing.T) {
		pyProc := NewPythonSingleProcessor(false)
		input := `#!/usr/bin/env python3
# noqa: E123
"""Module docstring."""
print("Hello") # A comment`
		expected := `"""Module docstring."""
print("Hello")
`
		actual, err := pyProc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestPythonProcessorGetLanguageName(t *testing.T) {
	processor := NewPythonSingleProcessor(false)
	assert.Equal(t, "python", processor.GetLanguageName())
}

func TestPythonProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewPythonSingleProcessor(true)
	processorWithoutDirectives := NewPythonSingleProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsPythonDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Shebang", "#!/usr/bin/env python3", true},
		{"ShebangWithSpace", "  #!/usr/bin/env python", true},
		{"Noqa", "# noqa: F401", true},
		{"NoqaInline", "x = 1  # noqa: E731", true},
		{"TypeComment", "# type: ignore", true},
		{"TypeCommentVar", "my_var = None  # type: Optional[str]", true},
		{"PylintDisable", "# pylint: disable=import-error", true},
		{"Flake8Noqa", "# flake8: noqa", true},
		{"MypyIgnore", "# mypy: ignore-errors", true},
		{"YapfOff", "# yapf: disable", true},
		{"IsortSkip", "# isort: skip_file", true},
		{"RuffNoqa", "# ruff: noqa: E501", true},
		{"FmtOff", "# fmt: off", true},
		{"FmtOn", "# fmt: on", true},
		{"NormalComment", "# This is just a normal comment", false},
		{"CommentWithHash", "# This has a # in it but not a directive", false},
		{"EmptyComment", "#", false},
		{"CodeLine", "print(\"hello\") # not a directive start", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, checkPythonSingleLineDirective(tt.line))
		})
	}
}
