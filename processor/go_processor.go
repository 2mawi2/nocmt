package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type GoProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewGoProcessor(preserveDirectives bool) *GoProcessor {
	return &GoProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *GoProcessor) GetLanguageName() string {
	return "go"
}

func (p *GoProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *GoProcessor) StripComments(source string) (string, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())

	if p.preserveDirectives {
		return p.stripCommentsPreserveDirectivesWithFiltering(source, p.isGoDirective, parser)
	}

	return p.stripCommentsWithFiltering(source, parser)
}

func (p *GoProcessor) isGoDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "//go:") ||
		strings.HasPrefix(trimmed, "// +build") ||
		strings.HasPrefix(trimmed, "//go:build") ||
		(strings.HasPrefix(trimmed, "//") && strings.Contains(line, "#include"))
}
