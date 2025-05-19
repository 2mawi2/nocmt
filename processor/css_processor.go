package processor

import (
	"fmt"
	"nocmt/config"
	"regexp"
	"strings"
)

func isCSSDirective(line string) bool {
	return false
}

func postProcessCSS(source string, _ []CommentRange, preserveDirectives bool) (string, error) {
	return source, nil
}

type CSSProcessor struct {
	preserveDirectives bool
}

func NewCSSProcessor(preserveDirectives bool) *CSSProcessor {
	return &CSSProcessor{preserveDirectives: preserveDirectives}
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
	cleaned := normalizeText(stripped)
	return PreserveOriginalTrailingNewline(source, cleaned), nil
}

func (p *CSSProcessor) GetLanguageName() string {
	return "css"
}

func (p *CSSProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *CSSProcessor) SetCommentConfig(cfg *config.Config) {}
