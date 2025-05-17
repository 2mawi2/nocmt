package processor

import (
	"regexp"
	"strings"
)

type JavaProcessor struct {
	*CoreProcessor
}

func isJavaDirective(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "// @formatter:") ||
		strings.HasPrefix(trimmed, "// @SuppressWarnings") ||
		strings.HasPrefix(trimmed, "//CHECKSTYLE") ||
		strings.Contains(trimmed, "@SuppressWarnings") ||
		strings.Contains(trimmed, "CHECKSTYLE.OFF") ||
		strings.Contains(trimmed, "CHECKSTYLE.ON") ||
		strings.Contains(trimmed, "NOCHECKSTYLE") ||
		strings.Contains(trimmed, "NOSONAR") ||
		strings.Contains(trimmed, "NOFOLINT")
}

func NewJavaProcessor(preserveDirectives bool) *JavaProcessor {
	core := NewCoreProcessor(
		"java",
		nil,
		isJavaDirective,
		nil,
	).WithPreserveDirectives(preserveDirectives)
	return &JavaProcessor{CoreProcessor: core}
}

func (p *JavaProcessor) GetLanguageName() string {
	return "java"
}

func (p *JavaProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *JavaProcessor) StripComments(source string) (string, error) {
	reBlock := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	text := reBlock.ReplaceAllString(source, "")
	lines := strings.Split(text, "\n")
	var out strings.Builder
	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if strings.HasPrefix(trimmed, "//") {
			if p.preserveDirectives {
				if isJavaDirective(ln) {
					out.WriteString(ln)
					out.WriteString("\n")
				} else {
					out.WriteString("\n")
				}
			}
			continue
		}
		if idx := strings.Index(ln, "//"); idx >= 0 {
			ln = ln[:idx]
		}
		ln = strings.TrimRight(ln, " \t")
		out.WriteString(ln)
		out.WriteString("\n")
	}
	return out.String(), nil
}
