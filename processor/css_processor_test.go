package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSSStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "strip line comments",
			input: `body {
  color: red; /* This is a comment */
  font-size: 16px; /* Another comment */
}`,
			expected: `body {
  color: red; 
  font-size: 16px; 
}`,
		},
		{
			name: "strip block comments",
			input: `/* Header comment */
body {
  /* Property comment */
  color: red;
  /* Another 
     block comment */
  font-size: 16px;
}
/* Footer comment */`,
			expected: `body {
  color: red;
  font-size: 16px;
}
`,
		},
		{
			name: "nested selectors with comments",
			input: `.container { /* Container styles */
  width: 100%;
  /* Inner styles */
  .inner {
    /* More comments */
    padding: 10px;
  }
}`,
			expected: `.container { 
  width: 100%;
  .inner {
    padding: 10px;
  }
}`,
		},
		{
			name: "comments inside strings",
			input: `body {
  content: "This is not a /* comment */";
  font-family: "Times /* not a comment */ New Roman";
}`,
			expected: `body {
  content: "This is not a /* comment */";
  font-family: "Times /* not a comment */ New Roman";
}`,
		},
		{
			name: "multiple adjacent comment lines",
			input: `/* First comment */
/* Second comment */
/* Third comment */

body {
  /* Comment group 1 */
  /* Comment group 2 */
  /* Comment group 3 */
  color: blue;
}`,
			expected: `body {
  color: blue;
}`,
		},
		{
			name: "comments with special characters",
			input: `/* Comment with UTF-8 characters: 你好, 世界! üñîçøðé */
/* Block comment with symbols: 
   @#$%^&*()_+-=[]{}|;:'",.<>/? 
*/
body {
  color: green;
}`,
			expected: `body {
  color: green;
}`,
		},
		{
			name: "media queries with comments",
			input: `/* Media query comment */
@media screen and (max-width: 768px) { /* Responsive styles */
  body { /* Mobile styles */
    font-size: 14px;
  }
}`,
			expected: `@media screen and (max-width: 768px) { 
  body { 
    font-size: 14px;
  }
}`,
		},
		{
			name: "keyframes with comments",
			input: `/* Animation comment */
@keyframes fade { /* Keyframe comment */
  0% { /* Start */
    opacity: 0;
  }
  100% { /* End */
    opacity: 1;
  }
}`,
			expected: `@keyframes fade { 
  0% { 
    opacity: 0;
  }
  100% { 
    opacity: 1;
  }
}`,
		},
		{
			name: "complex CSS with multiple comment types",
			input: `/* Main styles */
.container {
  display: flex; /* Use flexbox */
  /* Set dimensions */
  width: 100%;
  height: 100vh;
  /* Colors */
  background-color: #f5f5f5;
}

/* Navigation styles */
nav {
  /* Sizing */
  width: 250px;
  /* Position */
  position: fixed;
  top: 0;
  left: 0;
}`,
			expected: `.container {
  display: flex; 
  width: 100%;
  height: 100vh;
  background-color: #f5f5f5;
}

nav {
  width: 250px;
  position: fixed;
  top: 0;
  left: 0;
}`,
		},
		{
			name: "CSS variables with comments",
			input: `:root {
  /* Primary colors */
  --primary: #007bff;
  /* Secondary colors */
  --secondary: #6c757d;
}`,
			expected: `:root {
  --primary: #007bff;
  --secondary: #6c757d;
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &CSSProcessor{preserveDirectives: false}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSSStripCommentsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name: "only comments",
			input: `/* Comment 1 */
/* Comment 2 */
/* Comment 3 */`,
			expected: "",
		},
		{
			name: "unterminated comment",
			input: `body {
  color: red;
  /* This comment is not closed
}`,
			expected: "",
		},
		{
			name: "escaped characters in strings",
			input: `body {
  content: "This contains a \\\" and a \\*/";
  color: red;
}`,
			expected: `body {
  content: "This contains a \\\" and a \\*/";
  color: red;
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &CSSProcessor{preserveDirectives: false}
			if tt.name == "unterminated comment" {
				_, err := processor.StripComments(tt.input)
				assert.Error(t, err)
			} else {
				result, err := processor.StripComments(tt.input)
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCSSPreserveDirectives(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expected           string
		preserveDirectives bool
	}{
		{
			name: "preserve at-rules 1",
			input: `/* Comment */
@charset "UTF-8"; /* Charset comment */
@import url('styles.css'); /* Import comment */
body {
  color: blue;
}`,
			expected: `@charset "UTF-8"; 
@import url('styles.css'); 
body {
  color: blue;
}`,
			preserveDirectives: true,
		},
		{
			name: "preserve at-rules 2",
			input: `/* Comment */
@media screen and (max-width: 768px) { /* Media query comment */
  body {
    font-size: 14px; /* Font size comment */
  }
}`,
			expected: `@media screen and (max-width: 768px) { 
  body { 
    font-size: 14px;
  }
}`,
			preserveDirectives: true,
		},
		{
			name: "don't preserve at-rules when preserveDirectives is false",
			input: `/* Comment */
@media screen and (max-width: 768px) { /* Media query comment */
  body {
    font-size: 14px; /* Font size comment */
  }
}`,
			expected: `@media screen and (max-width: 768px) { 
  body { 
    font-size: 14px;
  }
}`,
			preserveDirectives: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := &CSSProcessor{preserveDirectives: tt.preserveDirectives}
			result, err := processor.StripComments(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSSProcessorGetLanguageName(t *testing.T) {
	processor := &CSSProcessor{}
	assert.Equal(t, "css", processor.GetLanguageName())
}

func TestCSSProcessorPreserveDirectives(t *testing.T) {
	processor := &CSSProcessor{preserveDirectives: true}
	assert.True(t, processor.PreserveDirectives())

	processor = &CSSProcessor{preserveDirectives: false}
	assert.False(t, processor.PreserveDirectives())
}
