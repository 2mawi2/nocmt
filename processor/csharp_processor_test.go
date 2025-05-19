package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSharpProcessor_FileBased(t *testing.T) {
	t.Run("Default_PreserveDirectives", func(t *testing.T) {
		csharpProc := NewCSharpSingleProcessor(true)
		RunFileBasedTestCaseNormalized(t, csharpProc, "../testdata/csharp/original.cs", "../testdata/csharp/expected.cs")
	})
	t.Run("RemoveAll_NoDirectives", func(t *testing.T) {
		csharpProc := NewCSharpSingleProcessor(false)
		input := `/// <summary>XML Doc</summary>
#pragma warning disable CS1591 // A directive
public class Test {} // A comment`
		expected := `public class Test {}`
		actual, err := csharpProc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("No_Line_Artifacts_When_No_Comments", func(t *testing.T) {
		csharpProc := NewCSharpSingleProcessor(true)
		RunFileBasedTestCaseNormalized(t, csharpProc, "../testdata/csharp/original_noline_artifacts.cs", "../testdata/csharp/expected_no_line_artifacts.cs")
	})
	t.Run("BlockCommentsPreserved", func(t *testing.T) {
		proc := NewCSharpSingleProcessor(false)
		input := "/* Block comment */\npublic class C {}\n"
		actual, err := proc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, input, actual)
	})
	t.Run("PreserveExistingTrailingNewline", func(t *testing.T) {
		proc := NewCSharpSingleProcessor(false)
		input := "public class C {}\n"
		actual, err := proc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, input, actual)
	})
	t.Run("NoTrailingNewlineAdded", func(t *testing.T) {
		proc := NewCSharpSingleProcessor(false)
		input := "public class C {}"
		actual, err := proc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, input, actual)
	})
	t.Run("PreserveTrailingNewlineOnRemoval_NoNewline", func(t *testing.T) {
		proc := NewCSharpSingleProcessor(false)
		input := "namespace X {\n    class C { } // comment\n}" 
		expected := "namespace X {\n    class C { }\n}"         
		actual, err := proc.StripComments(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestCSharpProcessorGetLanguageName(t *testing.T) {
	processor := NewCSharpSingleProcessor(false)
	assert.Equal(t, "csharp", processor.GetLanguageName())
}

func TestCSharpProcessorPreserveDirectivesFlag(t *testing.T) {
	processorWithDirectives := NewCSharpSingleProcessor(true)
	processorWithoutDirectives := NewCSharpSingleProcessor(false)

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
			assert.Equal(t, tt.expected, checkCSharpDirective(tt.line))
		})
	}
}
