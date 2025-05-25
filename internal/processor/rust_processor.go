package processor

import (
	"nocmt/internal/config"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/rust"
)

type RustSingleProcessor struct {
	*SingleLineCoreProcessor
}

func isRustSingleLineCommentNode(node *sitter.Node, sourceText string) bool {
	if node.Type() == "comment" || node.Type() == "line_comment" {
		startByte := int(node.StartByte())
		endByte := int(node.EndByte())

		
		if endByte-startByte >= 2 && sourceText[startByte] == '/' && sourceText[startByte+1] == '/' {
			
			if endByte-startByte >= 3 {
				thirdChar := sourceText[startByte+2]
				if thirdChar == '/' || thirdChar == '!' {
					return false
				}
			}

			
			lineStart := findLineStartBeforePosition(sourceText, startByte)
			return isOnlyWhitespaceBeforePosition(sourceText, lineStart, startByte)
		}
	}
	return false
}

func findLineStartBeforePosition(text string, pos int) int {
	for i := pos - 1; i >= 0; i-- {
		if text[i] == '\n' {
			return i + 1
		}
	}
	return 0
}

func isOnlyWhitespaceBeforePosition(text string, start, end int) bool {
	for i := start; i < end; i++ {
		if text[i] != ' ' && text[i] != '\t' {
			return false
		}
	}
	return true
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
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
