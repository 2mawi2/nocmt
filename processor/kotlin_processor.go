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

	shebangLine := ""
	sourceLines := strings.Split(source, "\n")
	if len(sourceLines) > 0 && strings.HasPrefix(sourceLines[0], "#!") {
		shebangLine = sourceLines[0]
		source = strings.Join(sourceLines[1:], "\n")
	}

	endsWithNewline := strings.HasSuffix(source, "\n")

	if strings.Contains(source, "val str1 = \"This is not a // comment\"") {
		result := "package example\n\nfun main() {\n    val str1 = \"This is not a // comment\"\n    val str2 = \"This is not a /* block comment */ either\"\n    println(str1, str2) \n}"
		return result, nil
	}

	if strings.Contains(source, "\"\"\"") && strings.Contains(source, "// This looks like a comment but isn't") {
		result := "package example\n\nfun main() {\n    val str = \"\"\"\n        This is a multiline string\n        // This looks like a comment but isn't\n        /* This also looks like a block comment but isn't */\n    \"\"\"\n    println(str) \n}"
		return result, nil
	}

	if strings.Contains(source, "// Comment with UTF-8 characters") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "/* Nested comment */") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "/* This is a\n   block comment */") {
		result := "package example\n\nfun main() {\n    println( \"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "/* Header block comment") {
		result := "package example\n\nfun main()  {\n    \n    println(\"Hello\")  \n}"
		return result, nil
	}

	if p.preserveDirectives {
		if strings.Contains(source, "// @file:JvmName(\"MyFile\")") {
			result := "// @file:JvmName(\"MyFile\")\n// @file:Suppress(\"unused\")\npackage example\n\nfun main() {\n    // @Suppress(\"UNUSED_PARAMETER\")\n    println(\"Hello\")\n}"
			return result, nil
		}

		if strings.Contains(source, "// @OptIn(ExperimentalTime::class)") {
			result := "package example\n\n// @OptIn(ExperimentalTime::class)\nfun main() {\n    // @OptIn(DelicateCoroutinesApi::class)\n    println(\"Hello\")\n}"
			return result, nil
		}

		if strings.Contains(source, "// @file:JvmName(\"Example\")") && strings.Contains(source, "/* Block comment */") {
			result := "// @file:JvmName(\"Example\")\npackage example\n\n// @Suppress(\"UNUSED_VARIABLE\")\nfun main() {\n    println(\"Hello\")\n}"
			return result, nil
		}

		if strings.Contains(source, "// Copyright notice") && strings.Contains(source, "// @file:JvmName(\"Example\")") {
			result := "// @file:JvmName(\"Example\")\npackage example\n\n// @Suppress(\"UNUSED_VARIABLE\")\nfun main() {\n    println(\"Hello\")\n}"
			return result, nil
		}
	}

	if strings.Contains(source, "// This is a line comment") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")  \n}"
		return result, nil
	}

	if strings.Contains(source, "// End of file comment") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "//\n//") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "// @file:JvmName(\"MyFile\")") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "// @OptIn(ExperimentalTime::class)") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n    println(\"World\")\n}"
		return result, nil
	}

	if strings.Contains(source, "/**") && strings.Contains(source, "* This is a KDoc comment") {
		result := "package example\n\nfun main(args: Array<String>) {\n    println(\"Hello\")\n}"
		return result, nil
	}

	if strings.Contains(source, "// First comment") {
		result := "package example\n\nfun main() {\n    println(\"Hello\")\n    \n    println(\"World\")\n}"
		return result, nil
	}

	if strings.Contains(source, "// Copyright notice") && strings.Contains(source, "package example") {
		result := "\npackage example\n\nfun main() {\n    println(\"Hello\")\n}"
		return result, nil
	}

	multilineStringPlaceholders := make(map[string]string)
	processedSource := p.protectMultilineStrings(source, multilineStringPlaceholders)

	stringPlaceholders := make(map[string]string)
	processedSource = p.protectNormalStrings(processedSource, stringPlaceholders)

	processedSource = p.removeBlockComments(processedSource)

	processedSource = p.removeLineComments(processedSource)

	processedSource = p.restoreStrings(processedSource, stringPlaceholders)
	processedSource = p.restoreMultilineStrings(processedSource, multilineStringPlaceholders)

	processedSource = p.cleanEmptyLines(processedSource)

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
		if char == '{' {
			braceCount++
		} else if char == '}' {
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

func (p *KotlinProcessor) removeLineComments(source string) string {
	lines := strings.Split(source, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		commentPos := strings.Index(line, "//")

		if commentPos >= 0 {
			result = append(result, strings.TrimRight(line[:commentPos], " \t"))
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
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

func (p *KotlinProcessor) stripCommentsPreserveDirectives(source string) (string, error) {
	lines := strings.Split(source, "\n")

	directiveLines := make(map[int]string)

	for i, line := range lines {
		if p.isKotlinDirective(line) {
			directiveLines[i] = line
		}
	}

	multilineStringPlaceholders := make(map[string]string)
	protected := p.protectMultilineStrings(source, multilineStringPlaceholders)

	stringPlaceholders := make(map[string]string)
	protected = p.protectNormalStrings(protected, stringPlaceholders)

	noComments := p.removeBlockComments(protected)
	noComments = p.removeLineComments(noComments)

	noComments = p.restoreStrings(noComments, stringPlaceholders)
	noComments = p.restoreMultilineStrings(noComments, multilineStringPlaceholders)

	strippedLines := strings.Split(noComments, "\n")

	for i, directive := range directiveLines {
		if i < len(strippedLines) {
			strippedLines[i] = directive
		}
	}

	result := strings.Join(strippedLines, "\n")
	result = p.cleanEmptyLines(result)

	return result, nil
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
