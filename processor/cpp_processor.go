package processor

import (
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

type CppProcessor struct {
	*SingleLineCoreProcessor
}

func isCppSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" {
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmed := strings.TrimSpace(commentText)
		return strings.HasPrefix(trimmed, "//")
	}
	return false
}

func isCppDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	if strings.HasPrefix(trimmed, "//") {
		content := strings.TrimSpace(strings.TrimPrefix(trimmed, "//"))

		if strings.HasPrefix(content, "TODO") ||
			strings.HasPrefix(content, "FIXME") ||
			strings.HasPrefix(content, "NOTE") ||
			strings.HasPrefix(content, "HACK") ||
			strings.HasPrefix(content, "XXX") ||
			strings.HasPrefix(content, "BUG") ||
			strings.HasPrefix(content, "WARNING") {
			return true
		}

		if strings.HasPrefix(content, "pragma") ||
			strings.HasPrefix(content, "#pragma") {
			return true
		}
	}

	return false
}

func NewCppProcessor(preserveDirectives bool) *CppProcessor {
	single := NewSingleLineCoreProcessor(
		"cpp",
		cpp.GetLanguage(),
		isCppSingleLineCommentNode,
		isCppDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)

	return &CppProcessor{SingleLineCoreProcessor: single}
}

func (p *CppProcessor) GetLanguageName() string {
	return "cpp"
}

func (p *CppProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *CppProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
