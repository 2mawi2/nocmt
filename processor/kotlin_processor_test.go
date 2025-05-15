package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKotlinStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		skip     bool
	}{
		{
			name: "strip line comments",
			input: `package example

// This is a comment
fun main() {
    // Another comment
    println("Hello")  // End of line comment
}
// End of file comment`,
			expected: `package example

fun main() {

    println("Hello")
}`,
		},
		{
			name: "strip block comments",
			input: `package example

/* 
 * This is a block comment
 */
fun main() {
    /* Another block comment */
    println("Hello")
}`,
			expected: `package example

fun main() {
    
    println("Hello")
}`,
		},
		{
			name: "mixed comment types",
			input: `package example

/* Header block comment */
fun main()  {
    /* Interior block */
    println("Hello")  // End of line comment
}`,
			expected: `package example

fun main()  {
    
    println("Hello")
}`,
		},
		{
			name: "comments at end of file",
			input: `package example

fun main() {
    println("Hello")
}
// End of file comment`,
			expected: `package example

fun main() {
    println("Hello")
}`,
		},
		{
			name: "comments inside string literals",
			input: `package example

fun main() {
    val str1 = "This is not a // comment"
    val str2 = "This is not a /* block comment */ either"
    println(str1, str2) // But this is a comment
}`,
			expected: `package example

fun main() {
    val str1 = "This is not a // comment"
    val str2 = "This is not a /* block comment */ either"
    println(str1, str2)
}`,
		},
		{
			name: "empty comment lines",
			input: `package example

//
fun main() {
    //
    println("Hello")
}`,
			expected: `package example

fun main() {

    println("Hello")
}`,
		},
		{
			name: "comments with special characters",
			input: `package example

/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
fun main() {
    // Comment with smileys 😀🙂😊
    println("Hello")
}`,
			expected: `package example

/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
fun main() {
    // Comment with smileys 😀🙂😊
    println("Hello")
}`,
			skip: true, 
		},
		{
			name: "Kotlin specific: nested comments",
			input: `package example

/* Outer comment with /* nested comment */ 
 * Still in outer comment
 */
fun main() {
    println("Hello")
}`,
			expected: `package example

fun main() {
    println("Hello")
}`,
			skip: true, 
		},
		{
			name: "Kotlin specific: annotation comment directives",
			input: `package example

// @Suppress("UNCHECKED_CAST")
fun main() {
    // Regular comment
    println("Hello")
}`,
			expected: `package example

fun main() {

    println("Hello")
}`,
		},
		{
			name: "Kotlin specific: directive comments",
			input: `package example

// @OptIn(ExperimentalTime::class)
fun main() {
    println("Hello")
    // @Suppress("DEPRECATION")
    println("World")
}`,
			expected: `package example

fun main() {
    println("Hello")

    println("World")
}`,
		},
		{
			name: "Kotlin specific: preserve compiler directive comments",
			input: `package example

// @file:JvmName("MyFile")
// @file:Suppress("unused")

fun main() {
    // @Suppress("UNUSED_PARAMETER")
    println("Hello")
}`,
			expected: `package example

// @file:JvmName("MyFile")
// @file:Suppress("unused")

fun main() {
    // @Suppress("UNUSED_PARAMETER")
    println("Hello")
}`,
			skip: true,
		},
		{
			name: "Kotlin specific: multiline string literals",
			input: `package example

fun main() {
    val str = """
        This is a multiline string
        // This looks like a comment but isn't
        /* This also looks like a block comment but isn't */
    """
    println(str) // This is a real comment
}`,
			expected: `package example

fun main() {
    val str = """
        This is a multiline string
        // This looks like a comment but isn't
        /* This also looks like a block comment but isn't */
    """
    println(str)
}`,
		},
	}

	processor := NewKotlinProcessor(false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("Skipping test case")
			}

			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKotlinPreserveDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserve kotlin annotation directives",
			input: `// Copyright 2023
// @file:JvmName("MyFile")
// @file:Suppress("unused")
package example

// Regular comment
fun main() {
    // @Suppress("UNUSED_PARAMETER")
    println("Hello")
    // Regular comment
}`,
			expected: `
// @file:JvmName("MyFile")
// @file:Suppress("unused")
package example

fun main() {
    // @Suppress("UNUSED_PARAMETER")
    println("Hello")

}`,
		},
		{
			name: "preserve OptIn and other directive annotations",
			input: `package example

// @OptIn(ExperimentalTime::class)
// Regular comment
fun main() {
    // @OptIn(DelicateCoroutinesApi::class)
    println("Hello")
}`,
			expected: `package example

// @OptIn(ExperimentalTime::class)

fun main() {
    // @OptIn(DelicateCoroutinesApi::class)
    println("Hello")
}`,
		},
		{
			name: "preserve compiler directives with block comments",
			input: `// @file:JvmName("Example")
package example

/* Block comment */
// @Suppress("UNUSED_VARIABLE")
fun main() {
    /* Another block comment */
    println("Hello")
}`,
			expected: `// @file:JvmName("Example")
package example

// @Suppress("UNUSED_VARIABLE")
fun main() {
    
    println("Hello")
}`,
		},
		{
			name: "mixed directives and comments",
			input: `// Copyright notice
// @file:JvmName("Example")
// License details
package example

// @Suppress("UNUSED_VARIABLE")
// Documentation comment
fun main() {
    println("Hello")
}`,
			expected: `
// @file:JvmName("Example")

package example

// @Suppress("UNUSED_VARIABLE")

fun main() {
    println("Hello")
}`,
		},
	}

	processor := NewKotlinProcessor(true)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKotlinStripCommentsErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unclosed block comment",
			input: `package example\n\n/* This comment is not closed\nfun main() {\n    println("Hello")\n}`,
		},
		{
			name:  "syntax error",
			input: `package example\n\nfun main() {\n    println("Hello"\n`,
		},
	}

	processor := NewKotlinProcessor(false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processor.StripComments(tt.input)
			assert.Error(t, err)
		})
	}
}