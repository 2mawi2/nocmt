package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	java "github.com/tree-sitter/tree-sitter-java"
)

type JavaProcessor struct {
	*SingleLineCoreProcessor
}

func isJavaDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	
	// Don't treat comments as directives
	if strings.HasPrefix(trimmed, "//") {
		return false
	}
	
	// Preserve Java annotations (e.g., @Override, @SuppressWarnings, etc.)
	return strings.HasPrefix(trimmed, "@") ||
		strings.HasPrefix(trimmed, "#!") ||
		strings.Contains(trimmed, "@")
}

func isJavaSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	nodeType := node.Type()
	
	if nodeType == "line_comment" || nodeType == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		
		// Only single-line comments (//), not multi-line (/* */) or doc comments (/** */)
		return strings.HasPrefix(trimmed, "//")
	}
	return false
}

func NewJavaProcessor(preserveDirectives bool) *JavaProcessor {
	single := NewSingleLineCoreProcessor(
		"java",
		java.GetLanguage(),
		isJavaSingleLineCommentNode,
		isJavaDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)
	return &JavaProcessor{SingleLineCoreProcessor: single}
}

func (p *JavaProcessor) GetLanguageName() string {
	return "java"
}

func (p *JavaProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *JavaProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}