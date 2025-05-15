package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoStripCommentsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "multiple adjacent comment lines",
			input: `package main

// First comment
// Second comment
// Third comment

func main() {
	fmt.Println("Hello")
	
	// Comment group 1
	// Comment group 2
	// Comment group 3
	fmt.Println("World")
}`,
			expected: `package main


func main() {
	fmt.Println("Hello")
	
	fmt.Println("World")
}`,
		},
		{
			name: "comments with code-like syntax",
			input: `package main

// func fakeFuncInComment() {
//     fmt.Println("fake")
// }

/* if (fakeCondition) {
   doSomething();
} */

func main() {
	// var x = 10;
	fmt.Println("Hello")
	/* for i := 0; i < 10; i++ {
	   doSomething()
	} */
}`,
			expected: `package main


func main() {
	
	fmt.Println("Hello")
}`,
		},
		{
			name: "URLs and special formatting in comments",
			input: `package main

// https://example.com/path?query=value#fragment
/* http://test.org/
 * email@example.com
 */
func main() {
	// TODO: Implement this
	// FIXME: Fix this issue
	// NOTE: Important information
	fmt.Println("Hello")
}`,
			expected: `package main

func main() {
	
	fmt.Println("Hello")
}`,
		},
		{
			name: "nested block comments",
			input: `package main

/* Outer comment start
   /* Nested comment */
   Outer comment end */
func main() {
	fmt.Println("Hello")
}`,
			expected: `package main

   Outer comment end */
func main() {
	fmt.Println("Hello")
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
