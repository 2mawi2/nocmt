package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPythonStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		skip     bool
	}{
		{
			name: "strip line comments",
			input: `#!/usr/bin/env python3
# This is a line comment
def main():
    # Another line comment
    print("Hello")  # End of line comment
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")  
`,
		},
		{
			name: "strip multi-line string comments (triple quotes)",
			input: `#!/usr/bin/env python3
"""This is a
   multi-line string comment"""
def main():
    print("""This is not a comment but a string""")
    '''
    Multi-line string with single quotes used as a doc string
    '''
    print("Hello")
`,
			expected: `#!/usr/bin/env python3
def main():
    print("""This is not a comment but a string""")
    print("Hello")
`,
		},
		{
			name: "preserve shebang line",
			input: `#!/usr/bin/env python3
# License information
# Author information

def main():
    print("Hello")
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")
`,
		},
		{
			name: "mixed comment types",
			input: `#!/usr/bin/env python3
# Header line comment
"""Module documentation
spanning multiple lines"""
# Another line comment
def main():  # function declaration comment
    # Code comment
    print("Hello")  # trailing line comment
    """
    This is a multi-line string
    used as a comment
    """
`,
			expected: `#!/usr/bin/env python3
def main():  
    print("Hello")  
`,
		},
		{
			name: "comments at end of file",
			input: `#!/usr/bin/env python3
def main():
    print("Hello")
# End of file comment
"""Final doc string comment"""
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")
`,
		},
		{
			name: "comments inside string literals",
			input: `#!/usr/bin/env python3
def main():
    str1 = "This is not a # comment"
    str2 = 'This is not a # comment either'
    str3 = """This string contains what looks like
    # a comment but it's not"""
    print(str1, str2, str3)  # This is a real comment
`,
			expected: `#!/usr/bin/env python3
def main():
    str1 = "This is not a # comment"
    str2 = 'This is not a # comment either'
    str3 = """This string contains what looks like
    # a comment but it's not"""
    print(str1, str2, str3)  
`,
		},
		{
			name: "empty comment lines",
			input: `#!/usr/bin/env python3
#
# 
#    
def main():
    #
    print("Hello")
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")
`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `#!/usr/bin/env python3
# First comment
# Second comment
# Third comment

def main():
    print("Hello")
    
    # Comment group 1
    # Comment group 2
    # Comment group 3
    print("World")
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")
    
    print("World")
`,
		},
		{
			name: "docstrings for functions and classes",
			input: `#!/usr/bin/env python3
class MyClass:
    """
    This is a class docstring.
    It describes the class purpose.
    """
    
    def __init__(self):
        """Initialize the class instance."""
        self.value = 42
        
    def my_method(self):
        """
        This is a method docstring.
        It should be removed.
        """
        # This is a comment
        return self.value
`,
			expected: `#!/usr/bin/env python3
class MyClass:
    
    def __init__(self):
        self.value = 42
        
    def my_method(self):
        return self.value
`,
		},
		{
			name: "Python f-strings with hash",
			input: `#!/usr/bin/env python3
def main():
    count = 42
    print(f"Count is {count}")  # This is a comment
    print(f"Hash symbol in f-string: #{count}")  # Another comment
`,
			expected: `#!/usr/bin/env python3
def main():
    count = 42
    print(f"Count is {count}")  
    print(f"Hash symbol in f-string: #{count}")  
`,
		},
		{
			name: "Python type hints with comments",
			input: `#!/usr/bin/env python3
def add(a: int, b: int) -> int:  # Function with type hints
    """Add two numbers and return the result."""
    # Calculate sum
    return a + b  # Return result

# Variable with type hint
x: int = 5  # Initialize x with 5
`,
			expected: `#!/usr/bin/env python3
def add(a: int, b: int) -> int:  
    return a + b  

x: int = 5  
`,
		},
		{
			name: "comments with special characters",
			input: `#!/usr/bin/env python3
# Comment with UTF-8 characters: 你好, 世界! üñîçøðé
"""Triple quoted string with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
"""
def main():
    print("Hello")
`,
			expected: `#!/usr/bin/env python3
def main():
    print("Hello")
`,
		},
		{
			name: "preserve type comments if requested",
			input: `#!/usr/bin/env python3
# Normal comment
x = []  # type: list[int]
def func(arg):
    # type: (str) -> int
    return len(arg)
`,
			expected: `#!/usr/bin/env python3
x = []  # type: list[int]
def func(arg):
    # type: (str) -> int
    return len(arg)
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
			processor := NewPythonProcessor(false)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPythonProcessorInterface(t *testing.T) {
	processor := NewPythonProcessor(false)
	assert.Equal(t, "python", processor.GetLanguageName())
	assert.False(t, processor.PreserveDirectives())

	processorWithDirectives := NewPythonProcessor(true)
	assert.True(t, processorWithDirectives.PreserveDirectives())
}

func TestPythonPreserveDirectives(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserve type hints",
			input: `#!/usr/bin/env python3
# Regular comment
x = []  # type: list[int]
def func(arg):
    # type: (str) -> int
    return len(arg)

y = 5  # normal comment
`,
			expected: `#!/usr/bin/env python3
x = []  # type: list[int]
def func(arg):
    # type: (str) -> int
    return len(arg)

y = 5  
`,
		},
		{
			name: "preserve specific directives",
			input: `#!/usr/bin/env python3
# mypy: ignore-errors
# pylint: disable=unused-import
# fmt: off
import os
import sys
# fmt: on

# Regular comment
def main():
    # Another comment
    print("Hello")
`,
			expected: `#!/usr/bin/env python3
# mypy: ignore-errors
# pylint: disable=unused-import
# fmt: off
import os
import sys
# fmt: on

def main():
    print("Hello")
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewPythonProcessor(true)
			result, err := processor.StripComments(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestPythonStripCommentsErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid syntax",
			input: "def invalid syntax(:)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := NewPythonProcessor(false)
			_, err := processor.StripComments(tc.input)
			assert.Error(t, err)
		})
	}
}