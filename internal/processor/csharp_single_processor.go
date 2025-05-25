package processor

import (
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/csharp"
)

type CSharpSingleProcessor struct {
	*SingleLineCoreProcessor
}

var csharpDirectiveRegex = regexp.MustCompile(`^\s*#(if|else|elif|endif|define|undef|region|endregion|pragma|nullable|line|error|warning)\b`)

func checkCSharpDirective(line string) bool {
	return csharpDirectiveRegex.MatchString(line)
}

func NewCSharpSingleProcessor(preserveDirectivesFlag bool) *CSharpSingleProcessor {
	isSLCommentNode := func(node *sitter.Node, sourceText string) bool {
		if node.Type() != "comment" {
			return false
		}
		commentText := sourceText[node.StartByte():node.EndByte()]
		trimmedText := strings.TrimSpace(commentText)

		
		if strings.HasPrefix(trimmedText, "//") && !strings.HasPrefix(trimmedText, "///") {
			if preserveDirectivesFlag {
				start := int(node.StartByte())
				lineStart := strings.LastIndex(sourceText[:start], "\n")
				if lineStart == -1 {
					lineStart = 0
				} else {
					lineStart++
				}
				line := sourceText[lineStart:start]
				if checkCSharpDirective(line) {
					return false
				}
			}
			return true
		}

		
		return false
	}

	
	simplifiedPostProcess := func(source string, preserveDirectives bool) (string, error) {
		lines := strings.Split(source, "\n")
		var resultLines []string

		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)

			
			if strings.HasPrefix(trimmedLine, "///") && !preserveDirectives {
				continue
			}

			
			if strings.HasPrefix(trimmedLine, "#") && !preserveDirectives {
				continue
			}

			resultLines = append(resultLines, line)
		}

		return strings.Join(resultLines, "\n"), nil
	}

	singleLineCore := NewSingleLineCoreProcessor(
		"csharp",
		csharp.GetLanguage(),
		isSLCommentNode,
		checkCSharpDirective,
		simplifiedPostProcess,
	).WithPreserveDirectives(preserveDirectivesFlag).PreserveBlankRuns()

	return &CSharpSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}

func (p *CSharpSingleProcessor) StripComments(source string) (string, error) {
	cleaned, err := p.SingleLineCoreProcessor.StripComments(source)
	if err != nil {
		return "", err
	}
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}
