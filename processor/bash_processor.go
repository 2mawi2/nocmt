package processor

import (
	"regexp"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
)

type BashProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewBashProcessor(preserveDirectives bool) *BashProcessor {
	return &BashProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *BashProcessor) GetLanguageName() string {
	return "bash"
}

func (p *BashProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *BashProcessor) StripComments(source string) (string, error) {
	shebangRegex := regexp.MustCompile(`^(#!.*)$`)
	lines := strings.Split(source, "\n")
	shebang := ""
	if len(lines) > 0 && shebangRegex.MatchString(lines[0]) {
		shebang = lines[0]
	}

	endsWithNewline := strings.HasSuffix(source, "\n")

	var directiveLines []int
	if p.preserveDirectives {
		for i, line := range lines {
			if p.isBashDirective(line) {
				directiveLines = append(directiveLines, i)
			}
		}
	}

	parser := sitter.NewParser()
	parser.SetLanguage(bash.GetLanguage())

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	if p.preserveDirectives && len(directiveLines) > 0 {
		directiveMap := make(map[int]bool)
		for _, line := range directiveLines {
			directiveMap[line] = true
		}

		lineStartPositions := make([]int, len(lines))
		pos := 0
		for i, line := range lines {
			lineStartPositions[i] = pos
			pos += len(line) + 1
		}

		var filteredRanges []CommentRange
		for _, r := range commentRanges {
			lineFound := false
			for i, start := range lineStartPositions {
				end := start + len(lines[i])
				if int(r.StartByte) >= start && int(r.StartByte) <= end {
					if !directiveMap[i] {
						filteredRanges = append(filteredRanges, r)
					}
					lineFound = true
					break
				}
			}
			if !lineFound {
				filteredRanges = append(filteredRanges, r)
			}
		}

		commentRanges = filteredRanges
	}

	result := removeComments(source, commentRanges)

	result = strings.TrimSpace(result)

	resultLines := strings.Split(result, "\n")

	var finalLines []string
	if shebang != "" {
		finalLines = append(finalLines, shebang)
	}

	finalLines = append(finalLines, resultLines...)

	finalResult := strings.Join(finalLines, "\n")
	if endsWithNewline && !strings.HasSuffix(finalResult, "\n") {
		finalResult += "\n"
	}

	return finalResult, nil
}

func (p *BashProcessor) isBashDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "# shellcheck")
}
