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

// This is a line comment
fun main() {
    // Another line comment
    println("Hello")  // End of line comment
}`,
			expected: `package example

fun main() {
    println("Hello")  
}`,
		},
		{
			name: "strip block comments",
			input: `package example

/* This is a
   block comment */
fun main() {
    println(/* inline block */ "Hello")
}`,
			expected: `package example

fun main() {
    println( "Hello")
}`,
		},
		{
			name: "mixed comment types",
			input: `package example

/* Header block comment
   spanning multiple lines */
// Line comment
fun main() /* function declaration comment */ {
    // Code comment
    println("Hello") /* trailing block */ // trailing line
}`,
			expected: `package example

fun main()  {
    
    println("Hello")  
}`,
		},
		{
			name: "comments before package declaration",
			input: `// Copyright notice
// License information

/* Package documentation
 * Provides main functionality
 */
package example

fun main() {
    println("Hello")
}`,
			expected: `
package example

fun main() {
    println("Hello")
}`,
		},
		{
			name: "comments at end of file",
			input: `package example

fun main() {
    println("Hello")
}
// End of file comment
/* Final block comment */
`,
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
    println(str1, str2) // This is a real comment
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
// 
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
			name: "multiple adjacent comment lines",
			input: `package example

// First comment
// Second comment
// Third comment

fun main() {
    println("Hello")
    
    // Comment group 1
    // Comment group 2
    // Comment group 3
    println("World")
}`,
			expected: `package example

fun main() {
    println("Hello")
    
    println("World")
}`,
		},
		{
			name: "comments with special characters",
			input: `package example

// Comment with UTF-8 characters: 你好, 世界! üñîçøðé
/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
fun main() {
    println("Hello")
}`,
			expected: `package example

fun main() {
    println("Hello")
}`,
		},
		{
			name: "Kotlin specific: KDoc comment",
			input: `package example

/**
 * This is a KDoc comment for the main function.
 * @param args command line arguments
 */
fun main(args: Array<String>) {
    println("Hello")
}`,
			expected: `package example

fun main(args: Array<String>) {
    println("Hello")
}`,
		},
		{
			name: "Kotlin specific: nested comments",
			input: `package example

/*
 * Outer comment
 * /* Nested comment */
 * Still in outer comment
 */
fun main() {
    println("Hello")
}`,
			expected: `package example

fun main() {
    println("Hello")
}`,
		},
		{
			name: "Kotlin specific: annotation comment directives",
			input: `package example

// @file:JvmName("MyFile")
// @file:Suppress("unused")

fun main() {
    // @Suppress("UNUSED_PARAMETER")
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
			expected: `// @file:JvmName("MyFile")
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
			expected: `// @file:JvmName("Example")
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