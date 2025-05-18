package processor

import (
	"nocmt/config"
	"testing"
)

func TestBaseProcessorCommentFiltering(t *testing.T) {
	tests := []struct {
		name     string
		comment  string
		patterns []string
		want     bool
	}{
		{
			name:     "no match without patterns",
			comment:  "// This is a comment",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "simple TODO match",
			comment:  "// TODO: implement this",
			patterns: []string{"TODO"},
			want:     true,
		},
		{
			name:     "prefix WHY match",
			comment:  "// WHY: because we need to",
			patterns: []string{"^\\s*//\\s*WHY"},
			want:     true,
		},
		{
			name:     "ticket number match",
			comment:  "// Fixes #1234",
			patterns: []string{"#\\d+"},
			want:     true,
		},
		{
			name:     "JIRA ticket match",
			comment:  "// TESTPROJECT-1250: Fixed login issue",
			patterns: []string{"TESTPROJECT-\\d+"},
			want:     true,
		},
		{
			name:     "no match with unrelated patterns",
			comment:  "// This is a regular comment",
			patterns: []string{"TODO", "FIXME", "#\\d+"},
			want:     false,
		},
		{
			name:     "match with one of multiple patterns",
			comment:  "// TODO: fix this later",
			patterns: []string{"FIXME", "TODO", "XXX"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New()
			err := cfg.SetCLIPatterns(tt.patterns)
			if err != nil {
				t.Fatalf("Failed to set patterns: %v", err)
			}

			base := BaseProcessor{
				commentConfig: cfg,
			}

			if got := base.ShouldIgnoreComment(tt.comment); got != tt.want {
				t.Errorf("BaseProcessor.ShouldIgnoreComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterCommentRanges(t *testing.T) {
	ranges := []CommentRange{
		{
			StartByte: 0,
			EndByte:   20,
			Content:   "// TODO: first task",
		},
		{
			StartByte: 25,
			EndByte:   45,
			Content:   "// Regular comment",
		},
		{
			StartByte: 50,
			EndByte:   80,
			Content:   "// This fixes #2345",
		},
	}

	cfg := config.New()
	err := cfg.SetCLIPatterns([]string{"TODO", "#\\d+"})
	if err != nil {
		t.Fatalf("Failed to set patterns: %v", err)
	}

	base := BaseProcessor{
		commentConfig: cfg,
	}

	filtered := base.filterCommentRanges(ranges)

	if len(filtered) != 1 {
		t.Errorf("Expected 1 comment range, got %d", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].Content != "// Regular comment" {
		t.Errorf("Expected to keep regular comment, got %s", filtered[0].Content)
	}
}

func TestEmptyLinePreservation(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   string
		skipReason string
	}{
		{
			name: "preserve empty lines without comments",
			input: `package main

import "fmt"

func main() {

	fmt.Println("Hello")

	// This is a comment
	fmt.Println("World")

}
`,
			expected: `package main

import "fmt"

func main() {

	fmt.Println("Hello")

	fmt.Println("World")

}
`,
		},
		{
			name: "preserve multiple consecutive empty lines",
			input: `package main

import "fmt"


// Comment between empty lines


func main() {
	fmt.Println("Hello")
}
`,
			expected: `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`,
		},
		{
			name: "remove empty lines with only comments on them",
			input: `package main

import "fmt"

// Comment on its own line
func main() {
	// Another comment line
	fmt.Println("Hello")
}
`,
			expected: `package main

import "fmt"

func main() {

	fmt.Println("Hello")
}
`,
		},
		{
			name: "preserve empty line at end of block",
			input: `package main

import "fmt"

func main() {
	if true {
		fmt.Println("Hello")
		// Comment followed by empty line

	}
}
`,
			expected: `package main

import "fmt"

func main() {
	if true {
		fmt.Println("Hello")

	}
}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.skipReason != "" {
				t.Skip(test.skipReason)
			}

			processor := NewGoProcessor(false)
			result, err := processor.StripComments(test.input)
			if err != nil {
				t.Fatalf("Error stripping comments: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, result)
			}
		})
	}
}
