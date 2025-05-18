package processor

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/typescript/tsx"
)

type TSXProcessor struct {
	*CoreProcessor
}

func isTSXDirective(line string) bool {
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

func NewTSXProcessor(preserveDirectives bool) *TSXProcessor {
	core := NewCoreProcessor(
		"tsx",
		tsx.GetLanguage(),
		isTSXDirective,
		postProcessTSX,
	).WithPreserveDirectives(preserveDirectives)
	return &TSXProcessor{CoreProcessor: core}
}

func (p *TSXProcessor) StripComments(source string) (string, error) {
	lang := tsx.GetLanguage()
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

	clean, err = postProcessTSX(clean, commentRanges, p.preserveDirectives)
	if err != nil {
		return "", err
	}

	return normalizeText(clean), nil
}

func postProcessTSX(content string, commentRanges []CommentRange, preserveDirectives bool) (string, error) {

	content, err := postProcessTypeScript(content, commentRanges, preserveDirectives)
	if err != nil {
		return "", err
	}

	emptyJSXCommentRegex := regexp.MustCompile(`{[ \t\r\n]*}`)
	content = emptyJSXCommentRegex.ReplaceAllString(content, "")

	inlineCommentInTemplateRegex := regexp.MustCompile(`('[^']*'|"[^"]*"|` + "`" + `[^` + "`" + `]*` + "`" + `)[ \t]*//.*`)
	content = inlineCommentInTemplateRegex.ReplaceAllString(content, "$1")

	malformedClosingRegex := regexp.MustCompile(`(\S+)[ \t]*\);(\s*\})`)
	content = malformedClosingRegex.ReplaceAllString(content, "$1\n  );$2")

	return content, nil
}
