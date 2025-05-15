package processor

import (
	"fmt"
	"regexp"
	"strings"
)

type KotlinProcessor struct {
	BaseProcessor
	preserveDirectives bool
}

func NewKotlinProcessor(preserveDirectives bool) *KotlinProcessor {
	return &KotlinProcessor{
		preserveDirectives: preserveDirectives,
	}
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

	// Handle shebang line if present
	shebangLine := ""
	sourceLines := strings.Split(source, "\n")
	if len(sourceLines) > 0 && strings.HasPrefix(sourceLines[0], "#!") {
		shebangLine = sourceLines[0]
		source = strings.Join(sourceLines[1:], "\n")
	}

	endsWithNewline := strings.HasSuffix(source, "\n")

	// If preserving directives is needed, identify directive lines
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

	// Process strings
	multilineStringPlaceholders := make(map[string]string)
	processedSource := p.protectMultilineStrings(source, multilineStringPlaceholders)

	stringPlaceholders := make(map[string]string)
	processedSource = p.protectNormalStrings(processedSource, stringPlaceholders)

	// Remove comments except for directives
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

	// Restore strings
	processedSource = p.restoreStrings(processedSource, stringPlaceholders)
	processedSource = p.restoreMultilineStrings(processedSource, multilineStringPlaceholders)

	// Restore directives if needed
	if p.preserveDirectives && len(directiveLines) > 0 {
		lines := strings.Split(processedSource, "\n")

		for i, directive := range directiveLines {
			if i < len(lines) {
				lines[i] = directive
			}
		}

		processedSource = strings.Join(lines, "\n")
	}

	// Clean up whitespace and empty lines
	processedSource = p.cleanEmptyLines(processedSource)

	// Restore shebang line if it was present
	if shebangLine != "" {
		processedSource = shebangLine + "\n" + processedSource
	}

	// Ensure the newline at the end matches the original
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

func (p *KotlinProcessor) cleanEmptyLines(source string) string {
	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		if i > 0 && strings.TrimSpace(lines[i]) == "" && strings.TrimSpace(lines[i-1]) == "" {
			continue
		}

		result = append(result, lines[i])
	}

	return strings.Join(result, "\n")
}

// isKotlinDirective checks if a line contains a Kotlin annotation/directive
func (p *KotlinProcessor) isKotlinDirective(line string) bool {
	trimmed := strings.TrimSpace(line)

	return strings.Contains(trimmed, "@file:") ||
		strings.HasPrefix(trimmed, "// @") ||
		strings.Contains(trimmed, "@Suppress") ||
		strings.Contains(trimmed, "@OptIn")
}
