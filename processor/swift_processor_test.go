package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwiftStripComments(t *testing.T) {
	t.Run("FileBased", func(t *testing.T) {
		processor := NewSwiftProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/swift/original.swift", "../testdata/swift/expected.swift")
	})

	tests := []struct {
		name     string
		input    string
		expected string
		skip     bool
	}{
		{
			name: "strip line comments",
			input: `// This is a line comment
func main() {
    // Another line comment
    print("Hello")  // End of line comment
}
`,
			expected: `func main() {
    print("Hello")  
}
`,
		},
		{
			name: "strip block comments",
			input: `/* This is a
   multi-line block comment */
func main() {
    /* Block comment */
    print("Hello")
    /* Another
       block comment */
}
`,
			expected: `func main() {
    print("Hello")
}
`,
		},
		{
			name: "nested block comments",
			input: `/* Outer comment 
   /* Nested comment */
   continues here */
func main() {
    print("Hello")
}
`,
			expected: `func main() {
    print("Hello")
}
`,
		},
		{
			name: "mixed comment types",
			input: `// Header line comment
/* Block comment
   spanning multiple lines */
// Another line comment
func main() {  // function declaration comment
    // Code comment
    print("Hello")  // trailing line comment
    /* Block comment inside function */
}
`,
			expected: `func main() {  
    print("Hello")  
}
`,
		},
		{
			name: "comments at end of file",
			input: `func main() {
    print("Hello")
}
// End of file comment
/* Final block comment */
`,
			expected: `func main() {
    print("Hello")
}
`,
		},
		{
			name: "comments inside string literals",
			input: `func main() {
    let str1 = "This is not a // comment"
    let str2 = "This is not a /* comment */ either"
    let str3 = """
    This multi-line string contains what looks like
    // a comment but it's not
    """
    print("\(str1) \(str2) \(str3)")  // This is a real comment
}
`,
			expected: `func main() {
    let str1 = "This is not a // comment"
    let str2 = "This is not a /* comment */ either"
    let str3 = """
    This multi-line string contains what looks like
    // a comment but it's not
    """
    print("\(str1) \(str2) \(str3)")  
}
`,
		},
		{
			name: "empty comment lines",
			input: `//
// 
//    
func main() {
    //
    print("Hello")
}
`,
			expected: `func main() {
    print("Hello")
}
`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `// First comment
// Second comment
// Third comment

func main() {
    print("Hello")
    
    // Comment group 1
    // Comment group 2
    // Comment group 3
    print("World")
}
`,
			expected: `func main() {
    print("Hello")
    
    print("World")
}
`,
		},
		{
			name: "documentation comments",
			input: `/// This is a documentation comment for the function
/// It spans multiple lines
func main() {
    /// This is a documentation comment inside the function
    print("Hello")
}

/// Documentation comment for a struct
struct Point {
    /// Documentation comment for x property
    var x: Int
    /// Documentation comment for y property
    var y: Int
}
`,
			expected: `func main() {
    print("Hello")
}

struct Point {
    var x: Int
    var y: Int
}
`,
		},
		{
			name: "markdown documentation comments",
			input: `/**
 This is a multi-line documentation comment
 with markdown formatting
 
 # Example Usage
 ` + "```" + `
 let p = Point(x: 10, y: 20)
 ` + "```" + `
 */
struct Point {
    var x: Int
    var y: Int
}
`,
			expected: `struct Point {
    var x: Int
    var y: Int
}
`,
		},
		{
			name: "comments with special characters",
			input: `// Comment with UTF-8 characters: 你好, 世界! üñîçøðé
/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
func main() {
    print("Hello")
}
`,
			expected: `func main() {
    print("Hello")
}
`,
		},
		{
			name: "preserve attributes and compiler directives",
			input: `@available(iOS 13.0, *)
// Regular comment that should be removed
struct ContentView {
    @State
    // This comment should be removed
    private var counter: Int = 0
    /* This comment should also be removed */
    @IBOutlet
    var label: UILabel!
}

#if DEBUG
// Debug-only code comment
func debugPrint() {
    print("Debug mode")
}
#endif
`,
			expected: `@available(iOS 13.0, *)
struct ContentView {
    @State
    private var counter: Int = 0
    @IBOutlet
    var label: UILabel!
}

#if DEBUG
func debugPrint() {
    print("Debug mode")
}
#endif
`,
		},
		{
			name: "preserve MARK, TODO, FIXME annotations",
			input: `class MyViewController: UIViewController {
    // MARK: - Properties
    var data: [String] = []
    
    // MARK: - Lifecycle Methods
    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
    }
    
    // MARK: - Private Methods
    private func setupUI() {
        // TODO: Implement the UI setup
        // FIXME: This causes a memory leak
    }
}
`,
			expected: `class MyViewController: UIViewController {
    // MARK: - Properties
    var data: [String] = []
    
    // MARK: - Lifecycle Methods
    override func viewDidLoad() {
        super.viewDidLoad()
        setupUI()
    }
    
    // MARK: - Private Methods
    private func setupUI() {
        // TODO: Implement the UI setup
        // FIXME: This causes a memory leak
    }
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
			processor := NewSwiftProcessor(false)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSwiftProcessorInterface(t *testing.T) {
	processor := NewSwiftProcessor(false)
	assert.Equal(t, "swift", processor.GetLanguageName())
	assert.False(t, processor.PreserveDirectives())

	processorWithDirectives := NewSwiftProcessor(true)
	assert.True(t, processorWithDirectives.PreserveDirectives())
}

func TestSwiftPreserveDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserve Swift attributes and property wrappers",
			input: `@available(iOS 13.0, *)
// Regular comment
struct ContentView {
    @State
    // This comment should be removed
    private var counter: Int = 0
}
`,
			expected: `@available(iOS 13.0, *)
struct ContentView {
    @State
    private var counter: Int = 0
}
`,
		},
		{
			name: "preserve Swift compiler directives",
			input: `// Regular comment
#if DEBUG
// Debug comment
func debugPrint() {
    print("Debug mode")
}
#endif

// Comment before directive
#if os(iOS)
func iOSOnly() {
    print("iOS only")
}
#elseif os(macOS)
func macOSOnly() {
    print("macOS only")
}
#else
func otherPlatform() {
    print("Other platform")
}
#endif
`,
			expected: `#if DEBUG
func debugPrint() {
    print("Debug mode")
}
#endif

#if os(iOS)
func iOSOnly() {
    print("iOS only")
}
#elseif os(macOS)
func macOSOnly() {
    print("macOS only")
}
#else
func otherPlatform() {
    print("Other platform")
}
#endif
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewSwiftProcessor(true)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSwiftStripCommentsErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "unterminated block comment",
			input: "func main() { /* This comment is not closed\n}",
		},
		{
			name:  "unterminated string",
			input: `func main() { let s = "this string is not closed; }`,
		},
		{
			name: "unterminated multi-line string",
			input: `func main() { let s = """
                     this multi-line string is not closed
                   }`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewSwiftProcessor(false)
			_, err := processor.StripComments(tc.input)
			assert.Error(t, err)
		})
	}
}
