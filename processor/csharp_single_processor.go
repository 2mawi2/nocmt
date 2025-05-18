package processor

import (
	"context"
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

func handleBlockComments(source string) (string, error) {
	parser := sitter.NewParser()
	parser.SetLanguage(csharp.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, []byte(source))
	if err != nil {
		return source, err
	}
	defer tree.Close()

	var blockCommentRanges []CommentRange

	Walk(tree.RootNode(), func(node *sitter.Node) bool {
		if node.Type() == "comment" {
			commentText := source[node.StartByte():node.EndByte()]
			trimmedText := strings.TrimSpace(commentText)
			if strings.HasPrefix(trimmedText, "/*") && strings.HasSuffix(trimmedText, "*/") {
				blockCommentRanges = append(blockCommentRanges, CommentRange{
					StartByte: node.StartByte(),
					EndByte:   node.EndByte(),
				})
				return false
			}
		}
		return true
	})

	for i := range blockCommentRanges {
		for j := i + 1; j < len(blockCommentRanges); j++ {
			if blockCommentRanges[i].StartByte < blockCommentRanges[j].StartByte {
				blockCommentRanges[i], blockCommentRanges[j] = blockCommentRanges[j], blockCommentRanges[i]
			}
		}
	}

	resultBytes := []byte(source)
	for _, r := range blockCommentRanges {
		start := int(r.StartByte)
		end := int(r.EndByte)

		if start >= len(resultBytes) || end > len(resultBytes) || start > end {
			continue
		}

		resultBytes = append(resultBytes[:start], resultBytes[end:]...)
	}

	return string(resultBytes), nil
}

func postProcessCSharpSingleLine(source string, preserveDirectives bool) (string, error) {

	originalLines := strings.Split(source, "\n")
	processedSource := source
	if !preserveDirectives {
		var err error
		processedSource, err = handleBlockComments(source)
		if err != nil {
			return source, err
		}
	}

	lines := strings.Split(processedSource, "\n")
	var resultLines []string
	var prevLineBlank = false

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if strings.HasPrefix(trimmedLine, "///") && !preserveDirectives {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "#") && !preserveDirectives {
			continue
		}

		if preserveDirectives && checkCSharpDirective(line) && !strings.Contains(line, "//") {
			if i < len(originalLines) {
				if idx := strings.Index(originalLines[i], "//"); idx != -1 {
					inline := originalLines[i][idx:]
					line = strings.TrimRight(line, " \t") + "  " + inline
				}
			}
		}

		isBlankLine := trimmedLine == ""

		if isBlankLine {
			if !prevLineBlank && i > 0 {
				resultLines = append(resultLines, line)
			}
			prevLineBlank = true
		} else {
			resultLines = append(resultLines, line)
			prevLineBlank = false
		}
	}

	return strings.Join(resultLines, "\n"), nil
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
	singleLineCore := NewSingleLineCoreProcessor(
		"csharp",
		csharp.GetLanguage(),
		isSLCommentNode,
		checkCSharpDirective,
		postProcessCSharpSingleLine,
	).WithPreserveDirectives(preserveDirectivesFlag)

	return &CSharpSingleProcessor{
		SingleLineCoreProcessor: singleLineCore,
	}
}
