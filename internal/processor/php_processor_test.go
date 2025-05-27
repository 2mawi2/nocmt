package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPHPProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewPHPProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../../testdata/php/original.php", "../../testdata/php/expected.php")
	})
	t.Run("WithoutDirectives", func(t *testing.T) {
		processor := NewPHPProcessor(false)
		input := "<?php\n// @license MIT\necho \"hello\"; // comment\n"
		expected := "<?php\necho \"hello\";\n"
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestPHPProcessorGetLanguageName(t *testing.T) {
	processor := NewPHPProcessor(false)
	assert.Equal(t, "php", processor.GetLanguageName())
}

func TestPHPProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewPHPProcessor(true)
	processorWithoutDirectives := NewPHPProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsPHPDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"Shebang", "#!/usr/bin/env php", true},
		{"PHP tag", "<?php", true},
		{"PHP close tag", "?>", true},
		{"License single line comment", "// @license MIT", true},
		{"License block comment", "/* @license GPL */", true},
		{"License shell comment", "# @license Apache", true},
		{"Preserve single line comment", "// @preserve", true},
		{"Preserve block comment", "/* @preserve */", true},
		{"CodingStandards ignore", "# @codingStandardsIgnoreStart", true},
		{"Phan directive", "// @phan-ignore-next-line", true},
		{"PHPStan directive", "# @phpstan-var string", true},
		{"Psalm directive", "/* @psalm-suppress PossiblyNullReference */", true},
		{"Generic @ in single line", "// @foo", true},
		{"Generic @ in block comment", "/* @bar */", true},
		{"Generic @ in shell comment", "# @baz", true},
		{"Simple line comment", "// This is a normal comment", false},
		{"Simple block comment", "/* This is a normal comment */", false},
		{"Simple shell comment", "# This is a normal comment", false},
		{"Code line", "echo \"hello // @world\";", false},
		{"Empty line", "", false},
		{"Whitespace line", "   ", false},
		{"Shebang with space", "  #!/usr/bin/php  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isPHPDirective(tt.line))
		})
	}
}
