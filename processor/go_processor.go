package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type GoSingleProcessor struct {
	*SingleLineCoreProcessor
}

func checkGoDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//go:") ||
		strings.HasPrefix(trimmed, "// +build") ||
		strings.HasPrefix(trimmed, "//go:build") {
		return true
	}
	if strings.HasPrefix(trimmed, "//") && (strings.Contains(line, "#cgo") || strings.Contains(line, "#include")) {
		return true
	}
	return false
}

func isGoSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		return strings.HasPrefix(strings.TrimSpace(commentText), "//")
	}
	return false
}

func postProcessGoSingleLine(source string, preserveDirectives bool) (string, error) {
	sourceLines := strings.Split(source, "\n")
	var tempLines []string
	for i, line := range sourceLines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if i > 0 && (i < len(sourceLines)-1 && strings.TrimSpace(sourceLines[i-1]) != "") {
				tempLines = append(tempLines, "")
			}
			continue
		}
		tempLines = append(tempLines, line)
	}
	var resultLines []string
	for i, line := range tempLines {
		trimmed := strings.TrimSpace(line)
		resultLines = append(resultLines, line)
		if !preserveDirectives && strings.HasSuffix(trimmed, "{") && i+1 < len(tempLines) && strings.TrimSpace(tempLines[i+1]) != "" {
			resultLines = append(resultLines, "")
		}
	}
	return strings.Join(resultLines, "\n"), nil
}

func NewGoSingleProcessor(preserveDirectivesFlag bool) *GoSingleProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"go",
		golang.GetLanguage(),
		isGoSingleLineCommentNode,
		checkGoDirective,
		postProcessGoSingleLine,
	).WithPreserveDirectives(preserveDirectivesFlag).PreserveBlankRuns()

	return &GoSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func NewGoProcessor(preserveDirectivesFlag bool) *GoSingleProcessor {
	return NewGoSingleProcessor(preserveDirectivesFlag)
}

func (p *GoSingleProcessor) StripComments(source string) (string, error) {
	return p.SingleLineCoreProcessor.StripComments(source)
}
