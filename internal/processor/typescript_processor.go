package processor

import (
	"nocmt/internal/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type TypeScriptProcessor struct {
	*SingleLineCoreProcessor
}

func isTypeScriptSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		return strings.HasPrefix(strings.TrimSpace(commentText), "//")
	}
	return false
}

func isTSDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "// @") ||
		strings.HasPrefix(trimmed, "/* @") ||
		strings.HasPrefix(trimmed, "//# sourceMappingURL=") ||
		strings.HasPrefix(trimmed, "//#") ||
		strings.HasPrefix(trimmed, "// =") ||
		strings.Contains(trimmed, "@preserve") ||
		strings.Contains(trimmed, "@license") {
		return true
	}
	return strings.HasPrefix(trimmed, "// @ts-") ||
		strings.HasPrefix(trimmed, "/* @ts-") ||
		strings.Contains(trimmed, "@ts-ignore") ||
		strings.Contains(trimmed, "@ts-nocheck") ||
		strings.Contains(trimmed, "@ts-check") ||
		strings.Contains(trimmed, "@ts-expect-error") ||
		strings.Contains(trimmed, "@jsx ") ||
		strings.HasPrefix(trimmed, "/// <reference")
}

func NewTypeScriptProcessor(preserveDirectives bool) *TypeScriptProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"typescript",
		typescript.GetLanguage(),
		isTypeScriptSingleLineCommentNode,
		isTSDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)

	return &TypeScriptProcessor{SingleLineCoreProcessor: singleLineCore}
}

func (p *TypeScriptProcessor) GetLanguageName() string {
	return "typescript"
}

func (p *TypeScriptProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *TypeScriptProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *TypeScriptProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
