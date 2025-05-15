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
		return p.stripCommentsPreserveDirectives(source, parser)
	}

	return p.stripCommentsWithFiltering(source, parser)
}

func (p *GoProcessor) stripCommentsPreserveDirectives(source string, parser *sitter.Parser) (string, error) {
	lines := strings.Split(source, "\n")
	directiveLines := make(map[int]bool)

	for i, line := range lines {
		if p.isGoDirective(line) {
			directiveLines[i] = true
		}
	}

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	filteredRanges := make([]CommentRange, 0)

	for _, r := range commentRanges {
		startLine, endLine := FindCommentLineNumbers(source, r)

		shouldPreserve := false
		for line := startLine; line <= endLine; line++ {
			if directiveLines[line-1] {
				shouldPreserve = true
				break
			}
		}

		if !shouldPreserve {
			filteredRanges = append(filteredRanges, r)
		}
	}

	if p.commentConfig != nil {
		var ignoreFilteredRanges []CommentRange
		for _, r := range filteredRanges {
			if !p.ShouldIgnoreComment(r.Content) {
				ignoreFilteredRanges = append(ignoreFilteredRanges, r)
			}
		}
		filteredRanges = ignoreFilteredRanges
	}

	return removeComments(source, filteredRanges), nil
}

func (p *GoProcessor) isGoDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "//go:") ||
		strings.HasPrefix(trimmed, "// +build") ||
		strings.HasPrefix(trimmed, "//go:build") ||
		(strings.HasPrefix(trimmed, "//") && strings.Contains(line, "#include"))
}