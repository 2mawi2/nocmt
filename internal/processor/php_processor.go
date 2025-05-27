package processor

import (
	"nocmt/internal/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/php"
)

type PHPSingleProcessor struct {
	*SingleLineCoreProcessor
}

func isPHPSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	nodeType := node.Type()
	if nodeType == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		return strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#")
	}
	return false
}

func isPHPDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#!/") ||
		strings.HasPrefix(trimmed, "<?php") ||
		strings.HasPrefix(trimmed, "?>") ||
		strings.HasPrefix(trimmed, "// @") ||
		strings.HasPrefix(trimmed, "/* @") ||
		strings.HasPrefix(trimmed, "# @") ||
		strings.Contains(trimmed, "@preserve") ||
		strings.Contains(trimmed, "@license") ||
		strings.Contains(trimmed, "@codingStandardsIgnore") ||
		strings.Contains(trimmed, "@phan-") ||
		strings.Contains(trimmed, "@phpstan-") ||
		strings.Contains(trimmed, "@psalm-")
}

func NewPHPProcessor(preserveDirectivesFlag bool) *PHPSingleProcessor {
	singleLineCore := NewSingleLineCoreProcessor(
		"php",
		php.GetLanguage(),
		isPHPSingleLineCommentNode,
		isPHPDirective,
		nil,
	).WithPreserveDirectives(preserveDirectivesFlag)

	return &PHPSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func (p *PHPSingleProcessor) GetLanguageName() string {
	return "php"
}

func (p *PHPSingleProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *PHPSingleProcessor) SetCommentConfig(cfg *config.Config) {
	p.commentConfig = cfg
}

func (p *PHPSingleProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
