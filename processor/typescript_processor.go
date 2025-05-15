package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

type TypeScriptProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewTypeScriptProcessor(preserveDirectives bool) *TypeScriptProcessor {
	return &TypeScriptProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *TypeScriptProcessor) GetLanguageName() string {
	return "typescript"
}

func (p *TypeScriptProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *TypeScriptProcessor) StripComments(source string) (string, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	if p.preserveDirectives {
		return stripCommentsPreserveDirectives(source, p.isTSDirective, parser)
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	return removeComments(source, commentRanges), nil
}

func (p *TypeScriptProcessor) isTSDirective(line string) bool {
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
