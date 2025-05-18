package processor

import (
	"nocmt/config"
	"testing"
)

type mockProcessor struct {
	languageName       string
	preserveDirectives bool
	commentConfig      *config.Config
	stripCommentsFunc  func(string) (string, error)
}

func (m *mockProcessor) GetLanguageName() string {
	return m.languageName
}

func (m *mockProcessor) PreserveDirectives() bool {
	return m.preserveDirectives
}

func (m *mockProcessor) StripComments(source string) (string, error) {
	return m.stripCommentsFunc(source)
}

func (m *mockProcessor) SetCommentConfig(cfg *config.Config) {
	m.commentConfig = cfg
}

func newMockProcessor(languageName string, preserveDirectives bool) *mockProcessor {
	return &mockProcessor{
		languageName:       languageName,
		preserveDirectives: preserveDirectives,
		stripCommentsFunc: func(source string) (string, error) {
			return source, nil
		},
	}
}

func TestFindCommentLineNumbers(t *testing.T) {
	content := `package main

import (
	"fmt"
)

// Comment on line 7
func main() {
	// Comment on line 9
	fmt.Println("Hello")
	
	/*
	 * Multi-line comment
	 * spanning lines 12-14
	 */
	fmt.Println("World")
}
`

	tests := []struct {
		name          string
		comment       CommentRange
		expectedStart int
		expectedEnd   int
	}{
		{
			name: "Single-line comment",
			comment: CommentRange{
				StartByte: 42,
				EndByte:   59,
				Content:   "// Comment on line 7",
			},
			expectedStart: 7,
			expectedEnd:   8,
		},
		{
			name: "Indented single-line comment",
			comment: CommentRange{
				StartByte: 71,
				EndByte:   89,
				Content:   "// Comment on line 9",
			},
			expectedStart: 9,
			expectedEnd:   9,
		},
		{
			name: "Multi-line comment",
			comment: CommentRange{
				StartByte: 114,
				EndByte:   168,
				Content: `/*
	 * Multi-line comment
	 * spanning lines 12-14
	 */`,
			},
			expectedStart: 12,
			expectedEnd:   15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := FindCommentLineNumbers(content, tt.comment)
			if start != tt.expectedStart {
				t.Errorf("FindCommentLineNumbers() start = %v, want %v", start, tt.expectedStart)
			}
			if end != tt.expectedEnd {
				t.Errorf("FindCommentLineNumbers() end = %v, want %v", end, tt.expectedEnd)
			}
		})
	}
}

func TestCommentOverlapsModifiedLines(t *testing.T) {
	tests := []struct {
		name          string
		commentStart  int
		commentEnd    int
		modifiedLines map[int]bool
		wantOverlaps  bool
	}{
		{
			name:         "Comment overlaps single modified line",
			commentStart: 5,
			commentEnd:   5,
			modifiedLines: map[int]bool{
				5: true,
			},
			wantOverlaps: true,
		},
		{
			name:         "Comment overlaps one of multiple modified lines",
			commentStart: 10,
			commentEnd:   12,
			modifiedLines: map[int]bool{
				5:  true,
				12: true,
				20: true,
			},
			wantOverlaps: true,
		},
		{
			name:         "Comment doesn't overlap any modified lines",
			commentStart: 7,
			commentEnd:   9,
			modifiedLines: map[int]bool{
				5:  true,
				15: true,
			},
			wantOverlaps: false,
		},
		{
			name:          "Empty modified lines",
			commentStart:  7,
			commentEnd:    9,
			modifiedLines: map[int]bool{},
			wantOverlaps:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			overlaps := CommentOverlapsModifiedLines(tt.commentStart, tt.commentEnd, tt.modifiedLines)
			if overlaps != tt.wantOverlaps {
				t.Errorf("CommentOverlapsModifiedLines() = %v, want %v", overlaps, tt.wantOverlaps)
			}
		})
	}
}

func TestFilterCommentsForRemoval(t *testing.T) {
	content := `package main

import (
	"fmt"
)

// Regular comment
func main() {
	// TODO: Preserve this comment
	fmt.Println("Hello")
	
	/*
	 * Another comment
	 */
	fmt.Println("World")
	
	//go:generate Some directive
}
	//go:generate Some directive
`
	commentConfig := config.New()
	_ = commentConfig.SetCLIPatterns([]string{"TODO"})

	modifiedLines := map[int]bool{
		7:  true,
		9:  true, // TODO comment
		13: true,
		17: true,
	}

	proc := newMockProcessor("go", true)
	proc.SetCommentConfig(commentConfig)

	comments := []CommentRange{
		{
			StartByte: 42,
			EndByte:   58,
			Content:   "// Regular comment",
		},
		{
			// TODO comment - should be preserved by pattern
			StartByte: 78, // Position of "// TODO: Preserve this comment"
			EndByte:   104,
			Content:   "// TODO: Preserve this comment",
		},
		{
			StartByte: 127,
			EndByte:   153,
			Content:   "/*\n\t * Another comment\n\t */",
		},
		{
			StartByte: 177,
			EndByte:   198,
			Content:   "//go:generate Some directive",
		},
	}

	commentsToRemove := FilterCommentsForRemoval(comments, content, modifiedLines, proc, true, commentConfig)

	// The TODO comment should be preserved due to pattern matching
	if len(commentsToRemove) != 2 {
		t.Errorf("FilterCommentsForRemoval() returned %d comments, want 2", len(commentsToRemove))
	}

	for _, comment := range commentsToRemove {
		if comment.Content == "// TODO: Preserve this comment" {
			t.Errorf("FilterCommentsForRemoval() incorrectly included comment with TODO pattern")
		}
		if comment.Content == "//go:generate Some directive" {
			t.Errorf("FilterCommentsForRemoval() incorrectly included directive comment")
		}
	}
}

func TestLanguageParsers(t *testing.T) {
	tests := []struct {
		name      string
		processor *mockProcessor
		wantNil   bool
	}{
		{"Go", newMockProcessor("go", false), false},
		{"JavaScript", newMockProcessor("javascript", false), false},
		{"Python", newMockProcessor("python", false), false},
		{"Rust", newMockProcessor("rust", false), false},
		{"Bash", newMockProcessor("bash", false), false},
		{"CSS", newMockProcessor("css", false), false},
		{"C#", newMockProcessor("csharp", false), false},
		{"Unknown", newMockProcessor("unknown", false), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := GetParserForProcessor(tt.processor)
			if (parser == nil) != tt.wantNil {
				if tt.wantNil {
					t.Errorf("%s parser returned non-nil, expected nil for unsupported language", tt.name)
				} else {
					t.Errorf("%s parser returned nil, expected a valid tree-sitter parser", tt.name)
				}
			}
		})
	}
}
