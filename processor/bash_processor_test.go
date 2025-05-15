package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "strip line comments",
			input: `#!/bin/bash

# This is a line comment
echo "Hello"  # End of line comment
`,
			expected: `#!/bin/bash
echo "Hello"
`,
		},
		{
			name: "comments inside string literals",
			input: `#!/bin/bash

echo "This is not a # comment"
echo 'This is not a # comment either'
echo "Hello"  # This is a real comment
`,
			expected: `#!/bin/bash
echo "This is not a # comment"
echo 'This is not a # comment either'
echo "Hello"
`,
		},
		{
			name: "empty comment lines",
			input: `#!/bin/bash

#
# 
#    
echo "Hello"
`,
			expected: `#!/bin/bash
echo "Hello"
`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `#!/bin/bash

# First comment
# Second comment
# Third comment

echo "Hello"

# Comment group 1
# Comment group 2
# Comment group 3
echo "World"
`,
			expected: `#!/bin/bash
echo "Hello"

echo "World"
`,
		},
		{
			name: "comments with special characters",
			input: `#!/bin/bash

# Comment with UTF-8 characters: 你好, 世界! üñîçøðé
echo "Hello"
`,
			expected: `#!/bin/bash
echo "Hello"
`,
		},
		{
			name: "comments within complex structures",
			input: `#!/bin/bash

if true; then # conditional comment
    for i in {1..10}; do # loop comment
        case $i in # switch comment
            1) # case comment
                echo $i
                ;;
            *) # default comment
                break
                ;;
        esac
    done
else # else comment
    exit 1
fi
`,
			expected: `#!/bin/bash
if true; then 
    for i in {1..10}; do 
        case $i in 
            1) 
                echo $i
                ;;
            *) 
                break
                ;;
        esac
    done
else 
    exit 1
fi
`,
		},
		{
			name: "comments with code-like syntax",
			input: `#!/bin/bash

# function fakeFuncInComment() {
#     echo "fake"
# }

function main() {
    # local x=10
    echo "Hello"
    # for i in {1..10}; do
    #   doSomething
    # done
}
`,
			expected: `#!/bin/bash
function main() {
    echo "Hello"
}
`,
		},
		{
			name: "here documents",
			input: `#!/bin/bash

cat << EOF
This is a here document
# This is not a comment
EOF

# This is a real comment
echo "Done"
`,
			expected: `#!/bin/bash
cat << EOF
This is a here document
# This is not a comment
EOF

echo "Done"
`,
		},
		{
			name: "comments in nested blocks",
			input: `#!/bin/bash

function outer() {
    # Outer comment
    function inner() {
        # Inner comment
        echo "Inner"
    }
    echo "Outer"
}
`,
			expected: `#!/bin/bash
function outer() {
    function inner() {
        echo "Inner"
    }
    echo "Outer"
}
`,
		},
		{
			name: "variable assignments with comments",
			input: `#!/bin/bash

NAME="John" # User name
AGE=30 # User age
echo "$NAME is $AGE years old"
`,
			expected: `#!/bin/bash
NAME="John" 
AGE=30 
echo "$NAME is $AGE years old"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &BashProcessor{preserveDirectives: false}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBashStripCommentsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name: "only comments",
			input: `# Comment 1
# Comment 2
# Comment 3`,
			expected: "",
		},
		{
			name: "shebang handling",
			input: `#!/bin/bash
# Setup script
echo "Setup complete"`,
			expected: `#!/bin/bash
echo "Setup complete"`,
		},
		{
			name: "escaped hash in strings",
			input: `#!/bin/bash
echo "Escaped \# symbol is not a comment"
echo "Regular # inside string also not a comment"
# This is a comment
`,
			expected: `#!/bin/bash
echo "Escaped \# symbol is not a comment"
echo "Regular # inside string also not a comment"
`,
		},
		{
			name: "mixed single and double quotes",
			input: `#!/bin/bash
echo 'Hash # in single quotes'
echo "Hash # in double quotes"
VAR='value with # symbol' # This is a comment
echo "$VAR"
`,
			expected: `#!/bin/bash
echo 'Hash # in single quotes'
echo "Hash # in double quotes"
VAR='value with # symbol' 
echo "$VAR"
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &BashProcessor{preserveDirectives: false}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBashDirectives(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expected           string
		preserveDirectives bool
	}{
		{
			name: "preserve shellcheck directives",
			input: `#!/bin/bash

# Regular comment
# shellcheck disable=SC2034
VAR="unused variable"
echo "Hello"
`,
			expected: `#!/bin/bash
# shellcheck disable=SC2034
VAR="unused variable"
echo "Hello"
`,
			preserveDirectives: true,
		},
		{
			name: "don't preserve directives when not enabled",
			input: `#!/bin/bash

# Regular comment
# shellcheck disable=SC2034
VAR="unused variable"
`,
			expected: `#!/bin/bash
VAR="unused variable"
`,
			preserveDirectives: false,
		},
		{
			name: "preserve various directives",
			input: `#!/bin/bash

# shellcheck disable=SC2034,SC2154
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Testing directives"
`,
			expected: `#!/bin/bash
# shellcheck disable=SC2034,SC2154
# shellcheck source=./lib.sh
# shellcheck shell=bash
echo "Testing directives"
`,
			preserveDirectives: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &BashProcessor{preserveDirectives: tt.preserveDirectives}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBashProcessorGetLanguageName(t *testing.T) {
	processor := NewBashProcessor(false)
	assert.Equal(t, "bash", processor.GetLanguageName())
}

func TestBashProcessorPreserveDirectives(t *testing.T) {
	processor := NewBashProcessor(true)
	assert.True(t, processor.PreserveDirectives())

	processor = NewBashProcessor(false)
	assert.False(t, processor.PreserveDirectives())
}

func TestBashDirectiveDetection(t *testing.T) {
	processor := &BashProcessor{}

	directives := []string{
		"# shellcheck disable=SC2034",
		"# shellcheck source=./lib.sh",
		"# shellcheck shell=bash",
	}

	for _, directive := range directives {
		assert.True(t, processor.isBashDirective(directive), "Should detect: %s", directive)
	}

	nonDirectives := []string{
		"# This is a regular comment",
		"echo 'Not a comment'",
		"#shellcheck", 
		"# not a shellcheck directive",
	}

	for _, nonDirective := range nonDirectives {
		assert.False(t, processor.isBashDirective(nonDirective), "Should not detect: %s", nonDirective)
	}
}