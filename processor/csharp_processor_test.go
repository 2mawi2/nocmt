package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSharpStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "strip line comments",
			input: `// This is a line comment
using System;

// Another line comment
class Program
{
	// Method comment
	static void Main()
	{
		Console.WriteLine("Hello");  // End of line comment
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");  
	}
}`,
		},
		{
			name: "strip block comments",
			input: `/* This is a
   block comment */
using System;

class Program
{
	static void Main()
	{
		Console.WriteLine(/* inline block */ "Hello");
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine( "Hello");
	}
}`,
		},
		{
			name: "mixed comment types",
			input: `/* Header block comment
   spanning multiple lines */
// Line comment
using System;

class Program /* class declaration comment */
{
	// Code comment
	static void Main()
	{
		Console.WriteLine("Hello"); /* trailing block */ // trailing line
	}
}`,
			expected: `using System;

class Program 
{
	static void Main()
	{
		Console.WriteLine("Hello"); 
	}
}`,
		},
		{
			name: "comments before namespace/class declaration",
			input: `// Copyright notice
// License information

/* Class documentation
 * Provides main functionality
 */
using System;

namespace MyApp
{
    class Program
    {
        static void Main()
        {
            Console.WriteLine("Hello");
        }
    }
}`,
			expected: `using System;

namespace MyApp
{
    class Program
    {
        static void Main()
        {
            Console.WriteLine("Hello");
        }
    }
}`,
		},
		{
			name: "comments at end of file",
			input: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}
// End of file comment
/* Final block comment */
`,
			expected: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}
`,
		},
		{
			name: "comments inside string literals",
			input: `using System;

class Program
{
	static void Main()
	{
		string str1 = "This is not a // comment";
		string str2 = "This is not a /* block comment */ either";
		Console.WriteLine(str1, str2); // This is a real comment
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		string str1 = "This is not a // comment";
		string str2 = "This is not a /* block comment */ either";
		Console.WriteLine(str1, str2); 
	}
}`,
		},
		{
			name: "empty comment lines",
			input: `//
// 
//    
using System;

class Program
{
	//
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `// First comment
// Second comment
// Third comment

using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
		
		// Comment group 1
		// Comment group 2
		// Comment group 3
		Console.WriteLine("World");
	}
}`,
			expected: `
using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
		
		Console.WriteLine("World");
	}
}`,
		},
		{
			name: "comments with special characters",
			input: `// Comment with UTF-8 characters: 你好, 世界! üñîçøðé
/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		Console.WriteLine("Hello");
	}
}`,
		},
		{
			name: "comments within complex structures",
			input: `using System;

class Program
{
	static void Main()
	{
		if (true) { // conditional comment
			for (int i = 0; i < 10; i++) { // loop comment
				switch (i) { // switch comment
				case 1: // case comment
					Console.WriteLine(i);
					break;
				default: /* default comment */
					break;
				}
			}
		} else /* else comment */ {
			return;
		}
	}
}`,
			expected: `using System;

class Program
{
	static void Main()
	{
		if (true) { 
			for (int i = 0; i < 10; i++) { 
				switch (i) { 
				case 1: 
					Console.WriteLine(i);
					break;
				default: 
					break;
				}
			}
		} else  {
			return;
		}
	}
}`,
		},
		{
			name: "XML documentation comments",
			input: `using System;

/// <summary>
/// This is a test class.
/// </summary>
public class Program
{
	/// <summary>
	/// Main entry point.
	/// </summary>
	public static void Main()
	{
		// Regular comment
		Console.WriteLine("Hello");
	}
}`,
			expected: `using System;

public class Program
{
	public static void Main()
	{
		Console.WriteLine("Hello");
	}
}`,
		},
		{
			name: "preprocessor directives",
			input: `using System;

#if DEBUG
// Debug-only code
Console.WriteLine("Debug mode");
#else
// Release-only code
Console.WriteLine("Release mode");
#endif

class Program
{
	static void Main()
	{
		#region Setup
		// Setup code
		var x = 10;
		#endregion

		Console.WriteLine(x);
	}
}`,
			expected: `using System;

#if DEBUG
Console.WriteLine("Debug mode");
#else
Console.WriteLine("Release mode");
#endif

class Program
{
	static void Main()
	{
		#region Setup
		var x = 10;
		#endregion

		Console.WriteLine(x);
	}
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &CSharpProcessor{preserveDirectives: false}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSharpStripCommentsWithDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserve pragma directives",
			input: `using System;

#pragma warning disable CS1591
// This comment should be removed
class Program
{
    static void Main()
    {
        // This should also be removed
        Console.WriteLine("Hello");
        #pragma warning restore CS1591
    }
}`,
			expected: `using System;

#pragma warning disable CS1591
class Program
{
    static void Main()
    {
        Console.WriteLine("Hello");
        #pragma warning restore CS1591
    }
}`,
		},
		{
			name: "preserve nullable directives",
			input: `using System;

// Top comment
#nullable enable
// Comment after directive
class Program
{
    static void Main()
    {
        // Method comment
        string? nullableString = null;
        Console.WriteLine(nullableString);
        #nullable disable
    }
}`,
			expected: `using System;

#nullable enable
class Program
{
    static void Main()
    {
        string? nullableString = null;
        Console.WriteLine(nullableString);
        #nullable disable
    }
}`,
		},
		{
			name: "preserve disable warnings",
			input: `using System;

// This is a regular comment
#pragma warning disable IDE0051 // Remove unused private members
// Another comment
class Program
{
    // This should be removed
    private int UnusedField; // This comment should be removed
}`,
			expected: `using System;

#pragma warning disable IDE0051 // Remove unused private members
class Program
{
    private int UnusedField; 
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &CSharpProcessor{preserveDirectives: true}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSharpProcessorGetLanguageName(t *testing.T) {
	processor := &CSharpProcessor{}
	assert.Equal(t, "csharp", processor.GetLanguageName())
}

func TestCSharpProcessorPreserveDirectives(t *testing.T) {
	processor := &CSharpProcessor{preserveDirectives: true}
	assert.True(t, processor.PreserveDirectives())

	processor = &CSharpProcessor{preserveDirectives: false}
	assert.False(t, processor.PreserveDirectives())
}

func TestCSharpDirectiveDetection(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{
			name:     "pragma warning disable",
			line:     "#pragma warning disable CS1591",
			expected: true,
		},
		{
			name:     "pragma warning restore",
			line:     "#pragma warning restore CS1591",
			expected: true,
		},
		{
			name:     "nullable enable",
			line:     "#nullable enable",
			expected: true,
		},
		{
			name:     "nullable disable",
			line:     "#nullable disable",
			expected: true,
		},
		{
			name:     "region directive",
			line:     "#region Setup",
			expected: true,
		},
		{
			name:     "endregion directive",
			line:     "#endregion",
			expected: true,
		},
		{
			name:     "if directive",
			line:     "#if DEBUG",
			expected: true,
		},
		{
			name:     "else directive",
			line:     "#else",
			expected: true,
		},
		{
			name:     "endif directive",
			line:     "#endif",
			expected: true,
		},
		{
			name:     "define directive",
			line:     "#define DEBUG",
			expected: true,
		},
		{
			name:     "undef directive",
			line:     "#undef DEBUG",
			expected: true,
		},
		{
			name:     "line directive",
			line:     "#line 100",
			expected: true,
		},
		{
			name:     "error directive",
			line:     "#error This is an error",
			expected: true,
		},
		{
			name:     "warning directive",
			line:     "#warning This is a warning",
			expected: true,
		},
		{
			name:     "regular comment",
			line:     "// This is a regular comment",
			expected: false,
		},
		{
			name:     "code line",
			line:     "var x = 10;",
			expected: false,
		},
	}

	processor := &CSharpProcessor{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.isCSharpDirective(tt.line)
			assert.Equal(t, tt.expected, result)
		})
	}
}
