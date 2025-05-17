package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSharpProcessor_FileBased(t *testing.T) {
	t.Run("WithDirectives", func(t *testing.T) {
		processor := NewCSharpProcessor(true)
		RunFileBasedTestCaseNormalized(t, processor, "../testdata/csharp/original.cs", "../testdata/csharp/expected.cs")
	})
	t.Run("WithoutDirectives_Simple", func(t *testing.T) {
		processor := NewCSharpProcessor(false)
		input := `/// <summary>XML Doc</summary>
#pragma warning disable CS1591 // A directive
public class Test {} // A comment`
		expected := `public class Test {}
` 
		actual, err := processor.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestCSharpProcessorGetLanguageName(t *testing.T) {
	processor := NewCSharpProcessor(false) 
	assert.Equal(t, "csharp", processor.GetLanguageName())
}

func TestCSharpProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewCSharpProcessor(true)
	processorWithoutDirectives := NewCSharpProcessor(false)

	assert.True(t, processorWithDirectives.PreserveDirectives())
	assert.False(t, processorWithoutDirectives.PreserveDirectives())
}

func TestIsCSharpDirective(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected bool
	}{
		{"If", "#if DEBUG", true},
		{"Else", "#else", true},
		{"Elif", "#elif TEST", true},
		{"Endif", "#endif", true},
		{"Define", "#define MY_CONSTANT", true},
		{"Undef", "#undef MY_CONSTANT", true},
		{"Region", "#region MyRegion", true},
		{"Endregion", "#endregion", true},
		{"PragmaWarning", "#pragma warning disable 1591", true},
		{"PragmaChecksum", "#pragma checksum \"file.cs\" \"{...}\" \"...\"", true},
		{"NullableEnable", "#nullable enable", true},
		{"LineHidden", "#line hidden", true},
		{"Error", "#error This is an error", true},
		{"Warning", "#warning This is a warning", true},
		{"StandardComment", "// This is a standard comment", false},
		{"XmlDocComment", "/// <summary>Test</summary>", false}, 
		{"CodeLine", "var x = 1; // #if DEBUG", false},
		{"EmptyLine", "", false},
		{"SpacedDirective", "  #if DEBUG  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isCSharpDirective(tt.line))
		})
	}
}

