package processor

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/bash"
)

type ShellProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewShellProcessor(preserveDirectives bool) *ShellProcessor {
	return &ShellProcessor{
		preserveDirectives: preserveDirectives,
	}
}

func (p *ShellProcessor) GetLanguageName() string {
	return "shell"
}

func (p *ShellProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *ShellProcessor) StripComments(source string) (string, error) {
	shebangRegex := regexp.MustCompile(`^(#!.*)$`)
	lines := strings.Split(source, "\n")
	shebang := ""
	if len(lines) > 0 && shebangRegex.MatchString(lines[0]) {
		shebang = lines[0]
	}

	parser := parsers.Get(bash.GetLanguage())
	defer parsers.Put(bash.GetLanguage(), parser)

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	var filteredRanges []CommentRange
	for _, r := range commentRanges {
		lineIdx := strings.Count(source[:int(r.StartByte)], "\n")
		if lineIdx == 0 {
			continue
		}
		if p.preserveDirectives && p.isShellDirective(lines[lineIdx]) {
			continue
		}
		filteredRanges = append(filteredRanges, r)
	}

	if len(filteredRanges) == 0 {
		return source, nil
	}

	result := removeComments(source, filteredRanges)

	resultLines := strings.Split(result, "\n")
	var finalLines []string
	for i, rl := range resultLines {
		if i == 0 && shebang != "" {
			finalLines = append(finalLines, shebang)
			continue
		}
		trimmed := strings.TrimRight(rl, " \t")
		if trimmed == "" {
			if i < len(lines) && strings.TrimSpace(lines[i]) == "" {
				finalLines = append(finalLines, rl)
			}
			continue
		}
		finalLines = append(finalLines, trimmed)
	}

	for len(finalLines) > 0 && finalLines[len(finalLines)-1] == "" {
		finalLines = finalLines[:len(finalLines)-1]
	}
	final := strings.Join(finalLines, "\n")
	if !strings.HasSuffix(final, "\n") {
		final += "\n"
	}
	return PreserveOriginalTrailingNewline(source, final), nil
}

func (p *ShellProcessor) isShellDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "# shellcheck")
}
