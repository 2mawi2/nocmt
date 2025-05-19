package processor

import (
	"nocmt/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type JavaScriptSingleProcessor struct {
	*SingleLineCoreProcessor
}

func isJavaScriptSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		return strings.HasPrefix(strings.TrimSpace(commentText), "//")
	}
	return false
}

func isJSDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "// @") ||
		strings.HasPrefix(trimmed, "/* @") ||
		strings.HasPrefix(trimmed, "//# sourceMappingURL=") ||
		strings.HasPrefix(trimmed, "//#") ||
		strings.HasPrefix(trimmed, "// =") ||
		strings.Contains(trimmed, "@preserve") ||
		strings.Contains(trimmed, "@license")
}

func NewJavaScriptProcessor(preserveDirectivesFlag bool) *JavaScriptSingleProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"javascript",
		javascript.GetLanguage(),
		isJavaScriptSingleLineCommentNode,
		isJSDirective,
		nil,
	).WithPreserveDirectives(preserveDirectivesFlag)

	return &JavaScriptSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func (p *JavaScriptSingleProcessor) GetLanguageName() string {
	return "javascript"
}

func (p *JavaScriptSingleProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *JavaScriptSingleProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *JavaScriptSingleProcessor) StripComments(source string) (string, error) {
	
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if cleaned == source {
		return source, nil
	}
	return cleaned, err
}
