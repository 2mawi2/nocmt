package processor

import (
	"strings"

	"github.com/smacker/go-tree-sitter/csharp"
)

type CSharpProcessor struct {
	*CoreProcessor
}

func isCSharpDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	directivePrefixes := []string{
		"#if", "#else", "#elif", "#endif",
		"#define", "#undef", "#region", "#endregion",
		"#pragma", "#nullable", "#line", "#error", "#warning",
	}
	for _, prefix := range directivePrefixes {
		if strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

func postProcessCSharp(source string, _ []CommentRange, preserveDirectives bool) (string, error) {
	lines := strings.Split(source, "\n")
	var resultLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "///") {
			continue
		}

		if !preserveDirectives && isCSharpDirective(line) { 
			continue
		}

		resultLines = append(resultLines, line)
	}
	return strings.Join(resultLines, "\n"), nil
}

func NewCSharpProcessor(preserveDirectivesFlag bool) *CSharpProcessor {
	core := NewCoreProcessor(
		"csharp",
		csharp.GetLanguage(),
		isCSharpDirective,
		postProcessCSharp, 
	).WithPreserveDirectives(preserveDirectivesFlag)
	return &CSharpProcessor{CoreProcessor: core}
}

