package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBashProcessor(t *testing.T) {
	t.Run("BasicCommentStripping", func(t *testing.T) {
		processor := NewBashProcessor(true)
		RunFileBasedTestCaseVeryLenient(t, processor, "../../testdata/bash/original.sh", "../../testdata/bash/expected.sh")
	})

	t.Run("DirectiveHandling", func(t *testing.T) {
		const bashWithDirectives = `#!/bin/bash
# Regular comment
# shellcheck disable=SC2034
VAR="unused variable"
# shellcheck source=./lib.sh
echo "Hello"
`

		const expectedPreserved = `#!/bin/bash
# shellcheck disable=SC2034
VAR="unused variable"
# shellcheck source=./lib.sh
echo "Hello"
`

		const expectedRemoved = `#!/bin/bash
VAR="unused variable"
echo "Hello"
`

		t.Run("PreserveDirectives", func(t *testing.T) {
			processor := NewBashProcessor(true)
			result, err := processor.StripComments(bashWithDirectives)
			assert.NoError(t, err)
			assert.Equal(t, expectedPreserved, result)
		})

		t.Run("RemoveDirectives", func(t *testing.T) {
			processor := NewBashProcessor(false)
			result, err := processor.StripComments(bashWithDirectives)
			assert.NoError(t, err)
			assert.Equal(t, expectedRemoved, result)
		})
	})

	t.Run("NoComments_Unchanged", func(t *testing.T) {
		const noComments = `#!/bin/bash
		echo "Hello"
		VAR=1
		`
		proc := NewBashProcessor(false)
		result, err := proc.StripComments(noComments)
		assert.NoError(t, err)
		assert.Equal(t, noComments, result)
		procPreserve := NewBashProcessor(true)
		resultPreserve, err := procPreserve.StripComments(noComments)
		assert.NoError(t, err)
		assert.Equal(t, noComments, resultPreserve)
	})

	t.Run("IndentationPreserved", func(t *testing.T) {
		const script = `#!/bin/bash
# a comment
    echo "Hello"
    VAR=1
`
		const expected = `#!/bin/bash
    echo "Hello"
    VAR=1
`
		proc := NewBashProcessor(false)
		result, err := proc.StripComments(script)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
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

func TestBashProcessorComplexCaseStatement(t *testing.T) {
	t.Run("ComplexCasePatternWithQuotedPipes", func(t *testing.T) {
		const bashWithComplexCase = `#!/usr/bin/env sh
# Session management for para

# Load session information from state file
get_session_info() {
  SESSION_ID="$1"
  STATE_FILE="$STATE_DIR/$SESSION_ID.state"
  [ -f "$STATE_FILE" ] || die "session '$SESSION_ID' not found"

  # Read state file with backward compatibility
  STATE_CONTENT=$(cat "$STATE_FILE")
  case "$STATE_CONTENT" in
  *"|"*"|"*"|"*)
    # New format with merge mode
    IFS='|' read -r TEMP_BRANCH WORKTREE_DIR BASE_BRANCH MERGE_MODE <"$STATE_FILE"
    ;;
  *)
    # Old format without merge mode, default to squash
    IFS='|' read -r TEMP_BRANCH WORKTREE_DIR BASE_BRANCH <"$STATE_FILE"
    MERGE_MODE="squash"
    ;;
  esac
}
`

		const expected = `#!/usr/bin/env sh

get_session_info() {
  SESSION_ID="$1"
  STATE_FILE="$STATE_DIR/$SESSION_ID.state"
  [ -f "$STATE_FILE" ] || die "session '$SESSION_ID' not found"

  STATE_CONTENT=$(cat "$STATE_FILE")
  case "$STATE_CONTENT" in
  *"|"*"|"*"|"*)
    IFS='|' read -r TEMP_BRANCH WORKTREE_DIR BASE_BRANCH MERGE_MODE <"$STATE_FILE"
    ;;
  *)
    IFS='|' read -r TEMP_BRANCH WORKTREE_DIR BASE_BRANCH <"$STATE_FILE"
    MERGE_MODE="squash"
    ;;
  esac
}
`

		processor := NewBashProcessor(false)
		result, err := processor.StripComments(bashWithComplexCase)
		assert.NoError(t, err, "Should parse complex case statement without syntax error")
		assert.Equal(t, expected, result)
	})

	t.Run("FallbackHandlesQuotedHashes", func(t *testing.T) {
		const bashWithQuotedHashes = `#!/bin/bash
# This is a comment
echo "This # is not a comment"
echo 'This # is also not a comment'
VAR="value # with hash" # This is a comment
echo "Escaped \" quote # still not a comment" # But this is
`

		const expected = `#!/bin/bash
echo "This # is not a comment"
echo 'This # is also not a comment'
VAR="value # with hash"
echo "Escaped \" quote # still not a comment"
`

		processor := NewBashProcessor(false)
		result, err := processor.StripComments(bashWithQuotedHashes)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
