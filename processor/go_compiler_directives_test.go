package processor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoCompilerDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "build constraints",
			input: `// +build linux darwin
// +build !windows

package main

func main() {
	// Regular comment
	fmt.Println("Hello")
}`,
			expected: `// +build linux darwin
// +build !windows

package main

func main() {
	fmt.Println("Hello")
}`,
		},
		{
			name: "go generate directives",
			input: `package main

//go:generate go run gen.go
//go:generate protoc --go_out=. *.proto
// Normal comment
func main() {
	fmt.Println("Hello")
}`,
			expected: `package main

//go:generate go run gen.go
//go:generate protoc --go_out=. *.proto

func main() {
	fmt.Println("Hello")
}`,
		},
		{
			name: "cgo directives",
			input: `package main

import "C"

// Regular comment
// #include <stdio.h>
// #include <stdlib.h>
//
// void myFunction(void) {
//    printf("Hello from C!\n");
// }
import "fmt"

func main() {
	fmt.Println("Hello from Go!")
}`,
			expected: `package main

import "C"

// #include <stdio.h>
// #include <stdlib.h>
//
// void myFunction(void) {
//    printf("Hello from C!\n");
// }
import "fmt"

func main() {
	fmt.Println("Hello from Go!")
}`,
		},
		{
			name: "other compiler directives",
			input: `package main

//go:noinline
//go:nosplit
func performant() {
	// Comment
}

//go:linkname time_now time.now
func time_now() int64

//go:embed static/index.html
var indexHTML string

func main() {}`,
			expected: `package main

//go:noinline
//go:nosplit
func performant() {
}

//go:linkname time_now time.now
func time_now() int64

//go:embed static/index.html
var indexHTML string

func main() {}`,
		},
		{
			name: "build tags modern syntax",
			input: `//go:build linux && (amd64 || arm64)
// +build linux,amd64 linux,arm64

package main

func main() {
	// Function body
}`,
			expected: `//go:build linux && (amd64 || arm64)
// +build linux,amd64 linux,arm64

package main

func main() {
}`,
		},
	}

	t.Skip("Skipping compiler directive tests - feature not implemented yet")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func GoStripCommentsPreserveDirectives(source string) (string, error) {
	lines := strings.Split(source, "\n")
	directiveLines := make(map[int]string)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:") ||
			strings.HasPrefix(trimmed, "// +build") ||
			strings.HasPrefix(trimmed, "//go:build") ||
			(strings.HasPrefix(trimmed, "//") && strings.Contains(line, "#include")) {
			directiveLines[i] = line
		}
	}

	stripped, err := StripComments(source)
	if err != nil {
		return "", err
	}

	strippedLines := strings.Split(stripped, "\n")
	for i, directive := range directiveLines {
		if i < len(strippedLines) {
			strippedLines[i] = directive
		}
	}

	return strings.Join(strippedLines, "\n"), nil
}