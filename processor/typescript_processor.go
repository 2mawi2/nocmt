package processor

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type TypeScriptProcessor struct {
	*CoreProcessor
}

func isTSDirective(line string) bool {
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

func NewTypeScriptProcessor(preserveDirectives bool) *TypeScriptProcessor {
	core := NewCoreProcessor(
		"typescript",
		typescript.GetLanguage(),
		isTSDirective,
		postProcessTypeScript,
	).WithPreserveDirectives(preserveDirectives)
	return &TypeScriptProcessor{CoreProcessor: core}
}

func (p *TypeScriptProcessor) StripComments(source string) (string, error) {
	lang := typescript.GetLanguage()
	parser := parsers.Get(lang)
	defer parsers.Put(lang, parser)

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	if p.commentConfig != nil {
		commentRanges = filterConfigIgnores(source, commentRanges, p.commentConfig)
	}

	for i, cr := range commentRanges {
		start := int(cr.StartByte)
		end := int(cr.EndByte)

		lineStart := strings.LastIndex(source[:start], "\n") + 1
		var lineEnd int
		if nl := strings.Index(source[end:], "\n"); nl >= 0 {
			lineEnd = end + nl + 1
		} else {
			lineEnd = len(source)
		}

		line := source[lineStart:lineEnd]
		without := strings.Replace(line, cr.Content, "", 1)

		if strings.TrimSpace(without) == "" {
			commentRanges[i].StartByte = uint32(lineStart)
			commentRanges[i].EndByte = uint32(lineEnd)
		}
	}

	if p.preserveDirectives && p.isDirective != nil {
		lines := splitIntoLines(source)
		linePos := calculateLinePositions(lines)

		directive := make(map[int]bool)
		for i, ln := range lines {
			if p.isDirective(ln) {
				directive[i] = true
			}
		}

		var keep []CommentRange
		for _, cr := range commentRanges {
			s, e := getCommentLineIndices(source, cr, linePos, lines)
			assoc := false
			for i := s; i <= e; i++ {
				if directive[i] {
					assoc = true
					break
				}
			}
			if !assoc {
				keep = append(keep, cr)
			}
		}
		commentRanges = keep
	}

	clean := removeComments(source, commentRanges)

	clean, err = postProcessTypeScript(clean, nil, false)
	if err != nil {
		return "", err
	}

	return normalizeText(clean), nil
}

func postProcessTypeScript(src string, _ []CommentRange, _ bool) (string, error) {
	s := regexp.MustCompile(`\n(?:[ \t]*\n){2,}`).ReplaceAllString(src, "\n\n")

	s = regexp.MustCompile(`\s+\)`).ReplaceAllString(s, ")")

	return s, nil
}
