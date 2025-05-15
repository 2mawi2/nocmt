package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKotlinDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "file annotations",
			input: `// @file:JvmName("TestFile")
// @file:Suppress("unused")
package example

fun main() {
    // This is a regular comment
    println("Hello")
}`,
			expected: `// @file:JvmName("TestFile")
// @file:Suppress("unused")
package example

fun main() {

    println("Hello")
}`,
		},
		{
			name: "function annotations",
			input: `package example

// This is a regular comment
// @OptIn(ExperimentalTime::class)
fun test() {
    // Another regular comment
    println("Test")
}`,
			expected: `package example

// @OptIn(ExperimentalTime::class)
fun test() {

    println("Test")
}`,
		},
		{
			name: "suppress warnings",
			input: `package example

class Example {
    // @Suppress("UNUSED_PARAMETER")
    fun test(unused: String) {
        // Regular comment
    }
}`,
			expected: `package example

class Example {
    // @Suppress("UNUSED_PARAMETER")
    fun test(unused: String) {

    }
}`,
		},
		{
			name: "mixed annotations and comments",
			input: `// Copyright notice
// @file:JvmName("Example")
package example

// Regular comment
// @Suppress("UNUSED_VARIABLE")
fun test() {
    // Comment here
    val x = 1
}`,
			expected: `
// @file:JvmName("Example")
package example

// @Suppress("UNUSED_VARIABLE")
fun test() {

    val x = 1
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewKotlinProcessor(true)
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
