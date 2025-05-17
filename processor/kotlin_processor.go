package processor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/smacker/go-tree-sitter/kotlin"
)

func postProcessKotlin(src string, _ []CommentRange, _ bool) (string, error) {
	onlyWS := regexp.MustCompile(`(?m)^[ \t]+\n`)
	s := onlyWS.ReplaceAllString(src, "\n")

	multiBlank := regexp.MustCompile(`\n(?:[ \t]*\n)+`)
	s = multiBlank.ReplaceAllString(s, "\n\n")

	return s, nil
}

type KotlinProcessor struct{ *CoreProcessor }

func NewKotlinProcessor(preserve bool) *KotlinProcessor {
	core := NewCoreProcessor(
		"kotlin",
		kotlin.GetLanguage(),
		isKotlinDirective,
		postProcessKotlin,
	).WithPreserveDirectives(preserve)

	return &KotlinProcessor{CoreProcessor: core}
}

func (p *KotlinProcessor) GetLanguageName() string {
	return "kotlin"
}

func (p *KotlinProcessor) PreserveDirectives() bool {
	return p.preserveDirectives
}

func (p *KotlinProcessor) StripComments(source string) (string, error) {
	err := p.validateKotlinSyntax(source)
	if err != nil {
		return "", err
	}

	shebangLine := ""
	sourceLines := strings.Split(source, "\n")
	if len(sourceLines) > 0 && strings.HasPrefix(sourceLines[0], "#!") {
		shebangLine = sourceLines[0]
		source = strings.Join(sourceLines[1:], "\n")
	}

	endsWithNewline := strings.HasSuffix(source, "\n")

	var directiveLines map[int]string
	if p.preserveDirectives {
		directiveLines = make(map[int]string)
		lines := strings.Split(source, "\n")

		for i, line := range lines {
			if p.isKotlinDirective(line) {
				directiveLines[i] = line
			}
		}
	}

	multilineStringPlaceholders := make(map[string]string)
	processedSource := p.protectMultilineStrings(source, multilineStringPlaceholders)

	stringPlaceholders := make(map[string]string)
	processedSource = p.protectNormalStrings(processedSource, stringPlaceholders)

	lines := strings.Split(processedSource, "\n")
	for i := range lines {
		if !p.preserveDirectives || directiveLines[i] == "" {
			commentPos := strings.Index(lines[i], "//")
			if commentPos >= 0 {
				lines[i] = strings.TrimRight(lines[i][:commentPos], " \t")
			}
		}
	}
	processedSource = strings.Join(lines, "\n")

	processedSource = p.removeBlockComments(processedSource)

	processedSource = p.restoreStrings(processedSource, stringPlaceholders)
	processedSource = p.restoreMultilineStrings(processedSource, multilineStringPlaceholders)

	var postErr error
	processedSource, postErr = postProcessKotlin(
		processedSource, nil, p.preserveDirectives)
	if postErr != nil {
		return "", postErr
	}

	if shebangLine != "" {
		processedSource = shebangLine + "\n" + processedSource
	}

	if endsWithNewline && !strings.HasSuffix(processedSource, "\n") {
		processedSource += "\n"
	} else if !endsWithNewline && strings.HasSuffix(processedSource, "\n") {
		processedSource = processedSource[:len(processedSource)-1]
	}

	return processedSource, nil
}

func (p *KotlinProcessor) validateKotlinSyntax(source string) error {
	if strings.Contains(source, "/*") && !strings.Contains(source, "*/") {
		return fmt.Errorf("invalid Kotlin syntax: unclosed block comment")
	}

	braceCount := 0
	for _, char := range source {
		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount < 0 {
				return fmt.Errorf("invalid Kotlin syntax: unmatched closing brace")
			}
		}
	}

	if braceCount != 0 {
		return fmt.Errorf("invalid Kotlin syntax: unmatched braces")
	}

	return nil
}

func (p *KotlinProcessor) protectMultilineStrings(source string, placeholders map[string]string) string {
	multilineRegex := regexp.MustCompile(`"""[\s\S]*?"""`)

	return multilineRegex.ReplaceAllStringFunc(source, func(match string) string {
		placeholder := fmt.Sprintf("__MULTILINE_STRING_PLACEHOLDER_%d__", len(placeholders))
		placeholders[placeholder] = match
		return placeholder
	})
}

func (p *KotlinProcessor) protectNormalStrings(source string, placeholders map[string]string) string {
	stringRegex := regexp.MustCompile(`"[^"\\]*(?:\\.[^"\\]*)*"`)

	return stringRegex.ReplaceAllStringFunc(source, func(match string) string {
		placeholder := fmt.Sprintf("__STRING_PLACEHOLDER_%d__", len(placeholders))
		placeholders[placeholder] = match
		return placeholder
	})
}

func (p *KotlinProcessor) removeBlockComments(source string) string {
	result := source
	for strings.Contains(result, "/*") && strings.Contains(result, "*/") {
		openPos := strings.Index(result, "/*")
		closePos := strings.Index(result[openPos:], "*/")

		if openPos == -1 || closePos == -1 {
			break
		}

		closePos = openPos + closePos + 2

		before := result[:openPos]
		after := result[closePos:]

		commentText := result[openPos:closePos]
		newlines := strings.Count(commentText, "\n")

		replacement := strings.Repeat("\n", newlines)
		result = before + replacement + after
	}

	return result
}

func (p *KotlinProcessor) restoreStrings(source string, placeholders map[string]string) string {
	result := source
	for placeholder, original := range placeholders {
		result = strings.ReplaceAll(result, placeholder, original)
	}
	return result
}

func (p *KotlinProcessor) restoreMultilineStrings(source string, placeholders map[string]string) string {
	result := source
	for placeholder, original := range placeholders {
		result = strings.ReplaceAll(result, placeholder, original)
	}
	return result
}

func (p *KotlinProcessor) isKotlinDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	return strings.Contains(trimmed, "@file:") ||
		strings.HasPrefix(trimmed, "// @") ||
		strings.Contains(trimmed, "@Suppress") ||
		strings.Contains(trimmed, "@OptIn")
}
