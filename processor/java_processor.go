package processor

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
)

type JavaProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewJavaProcessor(preserveDirectives bool) *JavaProcessor {
	return &JavaProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *JavaProcessor) GetLanguageName() string {
	return "java"
}

func (p *JavaProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *JavaProcessor) StripComments(source string) (string, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(java.GetLanguage())

	if strings.Contains(source, "/*") && !strings.Contains(source, "*/") {
		return "", fmt.Errorf("unclosed block comment detected")
	}

	if p.preserveDirectives {
		return stripCommentsPreserveDirectives(source, p.isJavaDirective, parser)
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	return removeComments(source, commentRanges), nil
}

func (p *JavaProcessor) isJavaDirective(line string) bool {
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
