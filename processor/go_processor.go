package processor

import (
	"strings"

	"github.com/smacker/go-tree-sitter/golang"
)

type GoProcessor struct {
	*CoreProcessor
}

func isGoDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, "//go:") || // e.g., //go:generate, //go:embed
		strings.HasPrefix(trimmed, "// +build") ||
		strings.HasPrefix(trimmed, "//go:build") {
		return true
	}
	if strings.HasPrefix(trimmed, "//") && (strings.Contains(line, "#cgo") || strings.Contains(line, "#include")) {
		return true
	}
	return false
}

func NewGoProcessor(preserveDirectives bool) *GoProcessor {
	core := NewCoreProcessor(
		"go",
		golang.GetLanguage(),
		isGoDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)
	return &GoProcessor{CoreProcessor: core}
}

func (p *GoProcessor) StripComments(source string) (string, error) {
	return p.CoreProcessor.StripComments(source)
}
