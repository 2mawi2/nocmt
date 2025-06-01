package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellProcessor(t *testing.T) {
	t.Run("BasicCommentStripping", func(t *testing.T) {
		processor := NewShellProcessor(true)
		RunFileBasedTestCaseVeryLenient(t, processor, "../../testdata/shell/original.sh", "../../testdata/shell/expected.sh")
	})

	t.Run("DirectiveHandling", func(t *testing.T) {
		const shellWithDirectives = `#!/bin/bash
# Regular comment  
# shellcheck disable=SC2034
VAR="unused variable"
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Hello"
`

		const expectedPreserved = `#!/bin/bash
# shellcheck disable=SC2034
VAR="unused variable"
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Hello"
`

		const expectedRemoved = `#!/bin/bash
VAR="unused variable"
echo "Hello"
`

		t.Run("PreserveDirectives", func(t *testing.T) {
			processor := NewShellProcessor(true)
			result, err := processor.StripComments(shellWithDirectives)
			assert.NoError(t, err)
			assert.Equal(t, expectedPreserved, result)
		})

		t.Run("RemoveDirectives", func(t *testing.T) {
			processor := NewShellProcessor(false)
			result, err := processor.StripComments(shellWithDirectives)
			assert.NoError(t, err)
			assert.Equal(t, expectedRemoved, result)
		})
	})

	t.Run("MultipleShellTypes", func(t *testing.T) {
		tests := []struct {
			name     string
			script   string
			expected string
		}{
			{
				name: "BashScript",
				script: `#!/bin/bash
# This is bash
echo "Hello from bash"
`,
				expected: `#!/bin/bash
echo "Hello from bash"
`,
			},
			{
				name: "ZshScript",
				script: `#!/bin/zsh
# This is zsh  
echo "Hello from zsh"
`,
				expected: `#!/bin/zsh
echo "Hello from zsh"
`,
			},
			{
				name: "FishScript",
				script: `#!/usr/bin/fish
# This is fish
echo "Hello from fish"
`,
				expected: `#!/usr/bin/fish
echo "Hello from fish"
`,
			},
			{
				name: "DashScript",
				script: `#!/bin/dash
# This is dash
echo "Hello from dash"
`,
				expected: `#!/bin/dash
echo "Hello from dash"
`,
			},
		}

		processor := NewShellProcessor(false)
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				result, err := processor.StripComments(test.script)
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			})
		}
	})

	t.Run("ComplexShellFeatures", func(t *testing.T) {
		const complexScript = `#!/bin/bash
# Function with comments
function test_func() {
    # Local variable comment
    local var="value"
    echo "$var"
}

# Array with comments
ARRAY=("item1" "item2") # Inline comment
echo "${ARRAY[0]}" # Access first element

# Here document with comments
cat << EOF
This is content
# This should remain
EOF

# Process substitution comment
diff <(ls dir1) <(ls dir2) # Compare directories
`

		const expected = `#!/bin/bash
function test_func() {
    local var="value"
    echo "$var"
}

ARRAY=("item1" "item2")
echo "${ARRAY[0]}"

cat << EOF
This is content
# This should remain
EOF

diff <(ls dir1) <(ls dir2)
`

		processor := NewShellProcessor(false)
		result, err := processor.StripComments(complexScript)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("StringHandling", func(t *testing.T) {
		const scriptWithStrings = `#!/bin/bash
echo "This # is not a comment"
echo 'Single quotes # also protect'
VAR='value # with hash' # This is a comment
echo "Escaped \" quote # still not a comment"
`

		const expected = `#!/bin/bash
echo "This # is not a comment"
echo 'Single quotes # also protect'
VAR='value # with hash'
echo "Escaped \" quote # still not a comment"
`

		processor := NewShellProcessor(false)
		result, err := processor.StripComments(scriptWithStrings)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("NoComments_Unchanged", func(t *testing.T) {
		const noComments = `#!/bin/bash
echo "Hello"
VAR=1
for i in {1..3}; do
    echo $i
done
`
		processor := NewShellProcessor(false)
		result, err := processor.StripComments(noComments)
		assert.NoError(t, err)
		assert.Equal(t, noComments, result)

		processorPreserve := NewShellProcessor(true)
		resultPreserve, err := processorPreserve.StripComments(noComments)
		assert.NoError(t, err)
		assert.Equal(t, noComments, resultPreserve)
	})

	t.Run("IndentationPreserved", func(t *testing.T) {
		const script = `#!/bin/bash
# Top level comment
if [ true ]; then
    # Indented comment
    echo "Hello"
    VAR=1
fi
`
		const expected = `#!/bin/bash
if [ true ]; then
    echo "Hello"
    VAR=1
fi
`
		processor := NewShellProcessor(false)
		result, err := processor.StripComments(script)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("EmptyAndWhitespaceComments", func(t *testing.T) {
		const script = `#!/bin/bash
#
#   
#	
# Regular comment
echo "test"
`
		const expected = `#!/bin/bash
echo "test"
`
		processor := NewShellProcessor(false)
		result, err := processor.StripComments(script)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestShellProcessorGetLanguageName(t *testing.T) {
	processor := NewShellProcessor(false)
	assert.Equal(t, "shell", processor.GetLanguageName())
}

func TestShellProcessorPreserveDirectives(t *testing.T) {
	processor := NewShellProcessor(true)
	assert.True(t, processor.PreserveDirectives())

	processor = NewShellProcessor(false)
	assert.False(t, processor.PreserveDirectives())
}

func TestShellDirectiveDetection(t *testing.T) {
	processor := &ShellProcessor{}

	directives := []string{
		"# shellcheck disable=SC2034",
		"# shellcheck source=./lib.sh",
		"# shellcheck shell=bash",
		"# shellcheck disable=SC1091,SC2034",
	}

	for _, directive := range directives {
		assert.True(t, processor.isShellDirective(directive), "Should detect: %s", directive)
	}

	nonDirectives := []string{
		"# This is a regular comment",
		"echo 'Not a comment'",
		"#shellcheck",
		"# not a shellcheck directive",
		"# shell check disable=SC2034", // with space
	}

	for _, nonDirective := range nonDirectives {
		assert.False(t, processor.isShellDirective(nonDirective), "Should not detect: %s", nonDirective)
	}
}
