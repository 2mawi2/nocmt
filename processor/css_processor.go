package processor

import (
	"fmt"
	"regexp"
	"strings"

	smkcss "github.com/smacker/go-tree-sitter/css"
)

func isCSSDirective(line string) bool {
	return false
}

func postProcessCSS(source string, _ []CommentRange, preserveDirectives bool) (string, error) {
	return source, nil
}

type CSSProcessor struct {
	*CoreProcessor
}

func NewCSSProcessor(preserveDirectives bool) *CSSProcessor {
	core := NewCoreProcessor(
		"css",
		smkcss.GetLanguage(),
		isCSSDirective,
		postProcessCSS,
	)
	core.WithPreserveDirectives(preserveDirectives)
	return &CSSProcessor{
		CoreProcessor: core,
	}
}

func (p *CSSProcessor) StripComments(source string) (string, error) {
	if strings.Contains(source, "/*") && !strings.Contains(source, "*/") {
		return "", fmt.Errorf("syntax error: unterminated comment")
	}
	re := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	if !re.MatchString(source) {
		return source, nil
	}
	stripped := re.ReplaceAllString(source, "")
	return normalizeText(stripped), nil
}
