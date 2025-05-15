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
			name: "go generate directive",
			input: `package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative protocol/test.proto

func main() {
	// This is a regular comment
}`,
			expected: `package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative protocol/test.proto

func main() {
}`,
		},
		{
			name: "cgo directive",
			input: `package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"

func main() {
	// Print something
	C.puts(C.CString("Hello, world"))
}`,
			expected: `package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"

func main() {
	C.puts(C.CString("Hello, world"))
}`,
		},
		{
			name: "go build tags",
			input: `//go:build linux && !windows
// +build linux,!windows

package main

func main() {
	// Comment
}`,
			expected: `//go:build linux && !windows
// +build linux,!windows

package main

func main() {
}`,
		},
		{
			name: "mixed directives and comments",
			input: `package main

//go:generate echo "Hello"
// This is a regular comment
//go:noinline
func example() {
	// Another comment
}`,
			expected: `package main

//go:generate echo "Hello"
//go:noinline
func example() {
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewGoProcessor(true)
			result, err := processor.StripComments(tt.input)
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
