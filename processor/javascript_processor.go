package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type JavaScriptProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewJavaScriptProcessor(preserveDirectives bool) *JavaScriptProcessor {
	return &JavaScriptProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *JavaScriptProcessor) GetLanguageName() string {
	return "javascript"
}

func (p *JavaScriptProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *JavaScriptProcessor) StripComments(source string) (string, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	if p.preserveDirectives {
		return stripCommentsPreserveDirectives(source, p.isJSDirective, parser)
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	return removeComments(source, commentRanges), nil
}

func (p *JavaScriptProcessor) isJSDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "// @") ||
		strings.HasPrefix(trimmed, "/* @") ||
		strings.HasPrefix(trimmed, "//# sourceMappingURL=") ||
		strings.HasPrefix(trimmed, "//#") ||
		strings.HasPrefix(trimmed, "// =") ||
		strings.Contains(trimmed, "@preserve") ||
		strings.Contains(trimmed, "@license")
}
