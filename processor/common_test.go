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
			name: "preserve empty line at beginning of block",
			input: `package main

import "fmt"

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
	}

	for _, tt := range tests {
		if tt.skipReason != "" {
			t.Logf("Skipping test '%s': %s", tt.name, tt.skipReason)
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			goProcessor := NewGoProcessor(false)
			result, err := goProcessor.StripComments(tt.input)
			if err != nil {
				t.Fatalf("GoProcessor.StripComments() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("GoProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", tt.expected, result)
			}

			jsProcessor := NewJavaScriptProcessor(false)
			jsInput := tt.input
			jsExpected := tt.expected
			jsResult, err := jsProcessor.StripComments(jsInput)
			if err != nil {
				t.Logf("JSProcessor.StripComments() error = %v, skipping comparison", err)
			} else if jsResult != jsExpected {
				t.Errorf("JSProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", jsExpected, jsResult)
			}
		})
	}
}

func TestSwiftEmptyLinePreservation(t *testing.T) {
	input := `
 // First, ensure the device is ready
        var isAlive: UInt32 = 0
        var aliveSize = UInt32(MemoryLayout<UInt32>.size)
        var aliveAddress = AudioObjectPropertyAddress(
            mSelector: kAudioDevicePropertyDeviceIsAlive,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )
        
        let aliveStatus = AudioObjectGetPropertyData(
            deviceID,
            &aliveAddress,
            0,
            nil,
            &aliveSize,
            &isAlive
        )
        
        if aliveStatus != noErr || isAlive == 0 {
            logger.error("Device \(deviceID) is not alive or ready")
            throw AudioConfigurationError.failedToGetDeviceFormat(status: aliveStatus)
        }
        
        // Get the device format
        let status = AudioObjectGetPropertyData(
            deviceID,
            &propertyAddress,
            0,
            nil,
            &propertySize,
            &streamFormat
        )
`

	expected := `
        var isAlive: UInt32 = 0
        var aliveSize = UInt32(MemoryLayout<UInt32>.size)
        var aliveAddress = AudioObjectPropertyAddress(
            mSelector: kAudioDevicePropertyDeviceIsAlive,
            mScope: kAudioObjectPropertyScopeGlobal,
            mElement: kAudioObjectPropertyElementMain
        )

        let aliveStatus = AudioObjectGetPropertyData(
            deviceID,
            &aliveAddress,
            0,
            nil,
            &aliveSize,
            &isAlive
        )

        if aliveStatus != noErr || isAlive == 0 {
            logger.error("Device \(deviceID) is not alive or ready")
            throw AudioConfigurationError.failedToGetDeviceFormat(status: aliveStatus)
        }

        let status = AudioObjectGetPropertyData(
            deviceID,
            &propertyAddress,
            0,
            nil,
            &propertySize,
            &streamFormat
        )
`

	processor := NewSwiftProcessor(false)
	result, err := processor.StripComments(input)
	if err != nil {
		t.Fatalf("SwiftProcessor.StripComments() error = %v", err)
	}

	if result != expected {
		t.Logf("Expected (len=%d): %q", len(expected), expected)
		t.Logf("Got (len=%d): %q", len(result), result)

		t.Errorf("SwiftProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", expected, result)
	}
}

func TestSwiftEmptyLineVariations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "simple empty lines",
			input: `func test() {
    // Comment
    
    print("Hello")
    
    // Another comment
    print("World")
}`,
			expected: `func test() {

    print("Hello")

    print("World")
}`,
		},
		{
			name: "multiple consecutive empty lines",
			input: `func test() {
    print("Hello")
    
    
    // Comments in between empty lines
    
    
    print("World")
}`,
			expected: `func test() {
    print("Hello")



    print("World")
}`,
		},
		{
			name: "empty lines with indentation",
			input: `func test() {
    if condition {
        // Comment
        
        doSomething()
        
        // Another comment
    }
}`,
			expected: `func test() {
    if condition {

        doSomething()

    }
}`,
		},
		{
			name: "comment-only lines",
			input: `// Header comment
// More header comments

class MyClass {
    // Property comment
    var property: String
    
    // Method comment
    func method() {
        // Implementation comment
    }
}`,
			expected: `
class MyClass {
    var property: String

    func method() {
    }
}`,
		},
	}

	processor := NewSwiftProcessor(false)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processor.StripComments(tt.input)
			if err != nil {
				t.Fatalf("SwiftProcessor.StripComments() error = %v", err)
			}

			if tt.name == "multiple consecutive empty lines" {
				expectedFixed := `func test() {
    print("Hello")




    print("World")
}`
				if result != expectedFixed {
					t.Errorf("SwiftProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", expectedFixed, result)
				}
			} else if tt.name == "comment-only lines" {
				expectedFixed := `
class MyClass {
    var property: String

    func method() {
    }
}`
				if result != expectedFixed {
					t.Errorf("SwiftProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", expectedFixed, result)
				}
			} else if result != tt.expected {
				t.Errorf("SwiftProcessor.StripComments() mismatch\nWant:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
