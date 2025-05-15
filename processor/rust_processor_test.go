package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRustStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		skip     bool
	}{
		{
			name: "strip line comments",
			input: `// This is a line comment
fn main() {
    // Another line comment
    println!("Hello");  // End of line comment
}
`,
			expected: `fn main() {
    println!("Hello");  
}
`,
		},
		{
			name: "strip block comments",
			input: `/* This is a
   multi-line block comment */
fn main() {
    /* Block comment */
    println!("Hello");
    /* Another
       block comment */
}
`,
			expected: `fn main() {
    println!("Hello");
}
`,
		},
		{
			name: "nested block comments",
			input: `/* Outer comment 
   /* Nested comment */
   continues here */
fn main() {
    println!("Hello");
}
`,
			expected: `fn main() {
    println!("Hello");
}
`,
		},
		{
			name: "mixed comment types",
			input: `// Header line comment
/* Block comment
   spanning multiple lines */
// Another line comment
fn main() {  // function declaration comment
    // Code comment
    println!("Hello");  // trailing line comment
    /* Block comment inside function */
}
`,
			expected: `fn main() {  
    println!("Hello");  
}
`,
		},
		{
			name: "comments at end of file",
			input: `fn main() {
    println!("Hello");
}
// End of file comment
/* Final block comment */
`,
			expected: `fn main() {
    println!("Hello");
}
`,
		},
		{
			name: "comments inside string literals",
			input: `fn main() {
    let str1 = "This is not a // comment";
    let str2 = "This is not a /* comment */ either";
    let str3 = r#"This raw string contains what looks like
    // a comment but it's not"#;
    println!("{} {} {}", str1, str2, str3);  // This is a real comment
}
`,
			expected: `fn main() {
    let str1 = "This is not a // comment";
    let str2 = "This is not a /* comment */ either";
    let str3 = r#"This raw string contains what looks like
    // a comment but it's not"#;
    println!("{} {} {}", str1, str2, str3);  
}
`,
		},
		{
			name: "empty comment lines",
			input: `//
// 
//    
fn main() {
    //
    println!("Hello");
}
`,
			expected: `fn main() {
    println!("Hello");
}
`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `// First comment
// Second comment
// Third comment

fn main() {
    println!("Hello");
    
    // Comment group 1
    // Comment group 2
    // Comment group 3
    println!("World");
}
`,
			expected: `fn main() {
    println!("Hello");
    
    println!("World");
}
`,
		},
		{
			name: "doc comments",
			input: `/// This is a doc comment for the function
/// It spans multiple lines
fn main() {
    //! This is an inner doc comment

    println!("Hello");
}

/// Doc comment for a struct
struct Point {
    /// Doc comment for x field
    x: i32,
    /// Doc comment for y field
    y: i32,
}
`,
			expected: `fn main() {
    println!("Hello");
}

struct Point {
    x: i32,
    y: i32,
}
`,
		},
		{
			name: "comments with special characters",
			input: `// Comment with UTF-8 characters: 你好, 世界! üñîçøðé
/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
fn main() {
    println!("Hello");
}
`,
			expected: `fn main() {
    println!("Hello");
}
`,
		},
		{
			name: "preserve attributes",
			input: `#![feature(test)]
#[derive(Debug)]
// Regular comment that should be removed
struct Point {
    #[deprecated]
    // This comment should be removed
    x: i32,
    /* This comment should also be removed */
    #[allow(dead_code)]
    y: i32,
}

#[cfg(test)]
mod tests {
    // Test comment
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
`,
			expected: `#![feature(test)]
#[derive(Debug)]
struct Point {
    #[deprecated]
    x: i32,
    #[allow(dead_code)]
    y: i32,
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
`,
		},
		{
			name: "preserve compiler directives if requested",
			input: `#![allow(unused_variables)]
// Regular comment
fn main() {
    #[allow(dead_code)]
    // This comment should be removed
    let x = 5;
}
`,
			expected: `#![allow(unused_variables)]
fn main() {
    #[allow(dead_code)]
    let x = 5;
}
`,
			skip: true, 
		},
	}

	for _, tc := range tests {
		if tc.skip {
			t.Logf("Skipping test case: %s", tc.name)
			continue
		}

		t.Run(tc.name, func(t *testing.T) {
			processor := NewRustProcessor(false)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRustProcessorInterface(t *testing.T) {
	processor := NewRustProcessor(false)
	assert.Equal(t, "rust", processor.GetLanguageName())
	assert.False(t, processor.PreserveDirectives())

	processorWithDirectives := NewRustProcessor(true)
	assert.True(t, processorWithDirectives.PreserveDirectives())
}

func TestRustPreserveDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserve attribute macros",
			input: `#![allow(unused_variables)]
// Regular comment
fn main() {
    #[allow(dead_code)]
    // This comment should be removed
    let x = 5;
}
`,
			expected: `#![allow(unused_variables)]
fn main() {
    #[allow(dead_code)]
    let x = 5;
}
`,
		},
		{
			name: "preserve cfg attributes",
			input: `// Regular comment
#[cfg(feature = "some_feature")]
// Another comment
fn conditional_function() {
    println!("This function is conditionally compiled");
}

// Comment before a module
#[cfg(test)]
mod tests {
    // Test comment
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
`,
			expected: `#[cfg(feature = "some_feature")]
fn conditional_function() {
    println!("This function is conditionally compiled");
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewRustProcessor(true)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRustStripCommentsErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unterminated block comment",
			input: "fn main() { /* This comment is not closed\n}",
		},
		{
			name:  "unterminated string",
			input: `fn main() { let s = "this string is not closed; }`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewRustProcessor(false)
			_, err := processor.StripComments(tc.input)
			assert.Error(t, err)
		})
	}
}