package processor

import (
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/javascript"
)

type JavaScriptProcessor struct {
	*CoreProcessor
}

func isJSDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "// @") ||
		strings.HasPrefix(trimmed, "/* @") ||
		strings.HasPrefix(trimmed, "//# sourceMappingURL=") ||
		strings.HasPrefix(trimmed, "//#") ||
		strings.HasPrefix(trimmed, "// =") ||
		strings.Contains(trimmed, "@preserve") ||
		strings.Contains(trimmed, "@license")
}

func NewJavaScriptProcessor(preserveDirectives bool) *JavaScriptProcessor {
	core := NewCoreProcessor(
		"javascript",
		javascript.GetLanguage(),
		isJSDirective,
		postProcessJavaScript,
	).WithPreserveDirectives(preserveDirectives)

	return &JavaScriptProcessor{CoreProcessor: core}
}

func (p *JavaScriptProcessor) StripComments(source string) (string, error) {
	if !strings.ContainsRune(source, '\n') && strings.Contains(source, `\n`) {
		source = strings.ReplaceAll(source, `\n`, "\n")
	}
	lang := javascript.GetLanguage()
	parser := parsers.Get(lang)
	defer parsers.Put(lang, parser)

	commentRanges, err := parseCode(parser, source)
	if err != nil {
		return "", err
	}

	if p.commentConfig != nil {
		commentRanges = filterConfigIgnores(source, commentRanges, p.commentConfig)
	}

	lines := splitIntoLines(source)
	lineStarts := calculateLinePositions(lines)

	for i, cr := range commentRanges {
		start := int(cr.StartByte)
		end := int(cr.EndByte)

		lineIdx := -1
		for idx, pos := range lineStarts {
			if start >= pos && start < pos+len(lines[idx])+1 {
				lineIdx = idx
				break
			}
		}
		if lineIdx == -1 {
			continue
		}

		lineStart := lineStarts[lineIdx]
		lineEnd := lineStart + len(lines[lineIdx])
		if lineEnd < len(source) {
			lineEnd++
		}

		if start < lineStart || end > lineEnd || lineStart > lineEnd {
			continue
		}

		before := strings.TrimSpace(source[lineStart:start])
		after := strings.TrimSpace(source[end:lineEnd])
		if before == "" && after == "" {

			newEnd := end
			if end > lineStart && source[end-1] == '\n' {
				newEnd = end - 1
			}
			commentRanges[i].StartByte = uint32(start)
			commentRanges[i].EndByte = uint32(newEnd)
		}
	}

	if p.preserveDirectives {
		var filtered []CommentRange
		for _, cr := range commentRanges {
			if p.isDirective != nil && p.isDirective(strings.TrimSpace(cr.Content)) {
				continue
			}
			filtered = append(filtered, cr)
		}
		commentRanges = filtered
	}

	for i, cr := range commentRanges {
		content := cr.Content
		if strings.HasPrefix(content, "/*") &&
			strings.Contains(content, "\n") {
			nl := strings.Index(content, "\n")
			if nl >= 0 {
				commentRanges[i].EndByte = cr.StartByte + uint32(nl+1)
			}
		}
	}

	clean := removeComments(source, commentRanges)

	clean, err = postProcessJavaScript(clean, nil, false)
	if err != nil {
		return "", err
	}

	return normalizeText(clean), nil
}

func postProcessJavaScript(src string, _ []CommentRange, _ bool) (string, error) {

	s := regexp.MustCompile(`\n(?:[ \t]*\n){2,}`).ReplaceAllString(src, "\n\n")
	s = regexp.MustCompile(`\s+\)`).ReplaceAllString(s, ")")
	return s, nil
}
