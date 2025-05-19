package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
)

type JavaProcessor struct {
	*SingleLineCoreProcessor
}

func isJavaDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "// @formatter:") ||
		strings.HasPrefix(trimmed, "// @SuppressWarnings") ||
		strings.HasPrefix(trimmed, "//CHECKSTYLE") ||
		strings.Contains(trimmed, "@SuppressWarnings") ||
		strings.Contains(trimmed, "CHECKSTYLE.OFF") ||
		strings.Contains(trimmed, "CHECKSTYLE.ON") ||
		strings.Contains(trimmed, "NOCHECKSTYLE") ||
		strings.Contains(trimmed, "NOSONAR") ||
		strings.Contains(trimmed, "NOFOLINT")
}

func NewJavaProcessor(preserveDirectives bool) *JavaProcessor {
	single := NewSingleLineCoreProcessor(
		"java",
		java.GetLanguage(),
		func(node *sitter.Node, src string) bool {
			return node.Type() == "line_comment"
		},
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
	if cleaned == source {
		return source, nil
	}
	return cleaned, err
}
