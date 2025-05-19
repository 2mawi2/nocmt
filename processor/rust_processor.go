package processor

import (
	"nocmt/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

type RustSingleProcessor struct {
	*SingleLineCoreProcessor
}

func isRustSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" || node.Type() == "line_comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		if strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "///") && !strings.HasPrefix(trimmed, "//!") {
			lineStart := strings.LastIndex(sourceText[:node.StartByte()], "\n") + 1
			prefix := sourceText[lineStart:node.StartByte()]
			if strings.TrimSpace(prefix) == "" {
				return true
			}
		}
	}
	return false
}

func isRustDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#!") || strings.HasPrefix(trimmed, "#[")
}

func NewRustProcessor(preserveDirectivesFlag bool) *RustSingleProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"rust",
		rust.GetLanguage(),
		isRustSingleLineCommentNode,
		isRustDirective,
		nil,
	).WithPreserveDirectives(preserveDirectivesFlag).PreserveBlankRuns()

	return &RustSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func (p *RustSingleProcessor) GetLanguageName() string {
	return "rust"
}

func (p *RustSingleProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *RustSingleProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *RustSingleProcessor) StripComments(source string) (string, error) {
	return p.SingleLineCoreProcessor.StripComments(source)
}
