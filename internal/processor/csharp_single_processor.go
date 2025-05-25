package processor

import (
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

type CSharpSingleProcessor struct {
	*SingleLineCoreProcessor
}

var csharpDirectiveRegex = regexp.MustCompile(`^\s*#(if|else|elif|endif|define|undef|region|endregion|pragma|nullable|line|error|warning)\b`)

func checkCSharpDirective(line string) bool {
	return csharpDirectiveRegex.MatchString(line)
}

func isCommentNode(node *sitter.Node) bool {
	return node.Type() == "comment"
}

func isLineComment(commentText string) bool {
	trimmed := strings.TrimSpace(commentText)
	return strings.HasPrefix(trimmed, "//")
}

func isDocumentationComment(commentText string) bool {
	trimmed := strings.TrimSpace(commentText)
	return strings.HasPrefix(trimmed, "///")
}

func isRegularLineComment(commentText string) bool {
	return isLineComment(commentText) && !isDocumentationComment(commentText)
}

func shouldPreserveCommentForDirective(node *sitter.Node, sourceText string) bool {
	commentStartPosition := int(node.StartByte())
	lineStartPosition := findLineStartPosition(sourceText, commentStartPosition)
	lineBeforeComment := sourceText[lineStartPosition:commentStartPosition]

	return checkCSharpDirective(lineBeforeComment)
}

func isCSharpLineCommentNode(node *sitter.Node, sourceText string, preserveDirectives bool) bool {
	if !isCommentNode(node) {
		return false
	}

	commentText := sourceText[node.StartByte():node.EndByte()]

	if !isRegularLineComment(commentText) {
		return false
	}

	if preserveDirectives && shouldPreserveCommentForDirective(node, sourceText) {
		return false
	}

	return true
}

func isXmlDocumentationComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "///")
}

func isPreprocessorDirective(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

func shouldSkipLine(line string, preserveDirectives bool) bool {
	if !preserveDirectives && isXmlDocumentationComment(line) {
		return true
	}

	if !preserveDirectives && isPreprocessorDirective(line) {
		return true
	}

	return false
}

func filterCSharpLines(sourceLines []string, preserveDirectives bool) []string {
	var processedLines []string

	for _, line := range sourceLines {
		if shouldSkipLine(line, preserveDirectives) {
			continue
		}
		processedLines = append(processedLines, line)
	}

	return processedLines
}

func postProcessCSharpSource(source string, preserveDirectives bool) (string, error) {
	sourceLines := strings.Split(source, "\n")
	processedLines := filterCSharpLines(sourceLines, preserveDirectives)
	return strings.Join(processedLines, "\n"), nil
}

func createCSharpCommentNodeChecker(preserveDirectivesFlag bool) func(*sitter.Node, string) bool {
	return func(node *sitter.Node, sourceText string) bool {
		return isCSharpLineCommentNode(node, sourceText, preserveDirectivesFlag)
	}
}

func NewCSharpSingleProcessor(preserveDirectivesFlag bool) *CSharpSingleProcessor {
	commentNodeChecker := createCSharpCommentNodeChecker(preserveDirectivesFlag)

	singleLineCore := NewSingleLineCoreProcessor(
		"csharp",
		csharp.GetLanguage(),
		commentNodeChecker,
		checkCSharpDirective,
		postProcessCSharpSource,
	).WithPreserveDirectives(preserveDirectivesFlag).PreserveBlankRuns()

	return &CSharpSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func findLineStartPosition(text string, position int) int {
	lineStartPosition := strings.LastIndex(text[:position], "\n")
	if lineStartPosition == -1 {
		return 0
	}
	return lineStartPosition + 1
}

func (p *CSharpSingleProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
