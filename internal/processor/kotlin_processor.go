package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/kotlin"
)

type KotlinProcessor struct {
	*SingleLineCoreProcessor
}

func isKotlinDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	
	if strings.HasPrefix(trimmed, "//") {
		return false
	}
	
	
	return strings.HasPrefix(trimmed, "@") ||
		strings.HasPrefix(trimmed, "#!") ||
		strings.Contains(trimmed, "@")
}

func isKotlinSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	nodeType := node.Type()
	
	if nodeType == "line_comment" || nodeType == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		
		
		return strings.HasPrefix(trimmed, "//")
	}
	return false
}

func NewKotlinProcessor(preserveDirectives bool) *KotlinProcessor {
	single := NewSingleLineCoreProcessor(
		"kotlin",
		kotlin.GetLanguage(),
		isKotlinSingleLineCommentNode,
		isKotlinDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)
	return &KotlinProcessor{SingleLineCoreProcessor: single}
}

func (p *KotlinProcessor) GetLanguageName() string {
	return "kotlin"
}

func (p *KotlinProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *KotlinProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
